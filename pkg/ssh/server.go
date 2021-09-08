//+build !windows

package ssh

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/creack/pty"
	"github.com/gliderlabs/ssh"
	"github.com/miguelmota/cointop/pkg/pathutil"
	gossh "golang.org/x/crypto/ssh"
)

// DefaultHostKeyFile is default SSH key path
var DefaultHostKeyFile = "~/.ssh/id_rsa"

// UserConfigTypePublicKey is constant
var UserConfigTypePublicKey = "public-key"

// UserConfigTypeUsername is constant
var UserConfigTypeUsername = "username"

// UserConfigTypeNone is constant
var UserConfigTypeNone = "none"

var UserConfigTypes = []string{UserConfigTypePublicKey, UserConfigTypeUsername, UserConfigTypeNone}

// Config is config struct
type Config struct {
	Port             uint
	Address          string
	IdleTimeout      time.Duration
	MaxTimeout       time.Duration
	ExecutableBinary string
	HostKeyFile      string
	MaxSessions      uint
	UserConfigType   string
	ColorsDir        string
}

// Server is server struct
type Server struct {
	port             uint
	address          string
	idleTimeout      time.Duration
	maxTimeout       time.Duration
	executableBinary string
	sshServer        *ssh.Server
	hostKeyFile      string
	maxSessions      uint
	sessionCount     uint
	userConfigType   string
	colorsDir        string
}

// NewServer returns a new server instance
func NewServer(config *Config) *Server {
	hostKeyFile := DefaultHostKeyFile
	if config.HostKeyFile != "" {
		hostKeyFile = config.HostKeyFile
	}
	userConfigType := config.UserConfigType
	if userConfigType == "" {
		userConfigType = UserConfigTypePublicKey
	}
	validateUserConfigType(userConfigType)
	hostKeyFile = pathutil.NormalizePath(hostKeyFile)
	colorsDir := pathutil.NormalizePath("~/.config/cointop/colors")
	if config.ColorsDir != "" {
		colorsDir = pathutil.NormalizePath(config.ColorsDir)
	}
	return &Server{
		port:             config.Port,
		address:          config.Address,
		idleTimeout:      config.IdleTimeout,
		maxTimeout:       config.MaxTimeout,
		executableBinary: config.ExecutableBinary,
		hostKeyFile:      hostKeyFile,
		maxSessions:      config.MaxSessions,
		userConfigType:   userConfigType,
		colorsDir:        colorsDir,
	}
}

// ListenAndServe starts the server
func (s *Server) ListenAndServe() error {
	s.sshServer = &ssh.Server{
		Addr:        fmt.Sprintf("%s:%v", s.address, s.port),
		IdleTimeout: s.idleTimeout,
		MaxTimeout:  s.maxTimeout,
		Handler: func(sshSession ssh.Session) {
			if s.maxSessions > 0 {
				s.sessionCount++
				defer func() {
					s.sessionCount--
				}()
				if s.sessionCount > s.maxSessions {
					io.WriteString(sshSession, "Error: Maximum sessions reached. Must wait until session slot is available.")
					sshSession.Exit(1)
					return
				}
			}

			cmdUserArgs := sshSession.Command()
			ptyReq, winCh, isPty := sshSession.Pty()
			if !isPty {
				io.WriteString(sshSession, "Error: Non-interactive terminals are not supported")
				sshSession.Exit(1)
				return
			}

			configDir := ""
			configDirKey := ""

			switch s.userConfigType {
			case UserConfigTypePublicKey:
				pubKey := sshSession.PublicKey()
				if pubKey != nil {
					pubBytes := pubKey.Marshal()
					if len(pubBytes) > 0 {
						hash := sha256.Sum256(pubBytes)
						configDirKey = fmt.Sprintf("%x", hash)
					}
				}
			case UserConfigTypeUsername:
				user := sshSession.User()
				if user != "" {
					configDirKey = user
				}
			}

			if configDirKey != "" {
				configDir = fmt.Sprintf("/tmp/cointop_config/%s", configDirKey)
				err := os.MkdirAll(configDir, 0700)
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			if configDir == "" {
				tempDir, err := createTempDir()
				if err != nil {
					fmt.Println(err)
					return
				}
				configDir = tempDir
				defer os.RemoveAll(configDir)
			}

			configPath := fmt.Sprintf("%s/config", configDir)

			cmdCtx, cancelCmd := context.WithCancel(sshSession.Context())
			defer cancelCmd()

			flags := []string{
				"--reset",
				"--silent",
				"--cache-dir",
				configDir,
				"--config",
				configPath,
				"--colors-dir",
				s.colorsDir,
			}

			if len(cmdUserArgs) > 0 {
				if cmdUserArgs[0] == "cointop" {
					cmdUserArgs = cmdUserArgs[1:]
				}
			}

			if len(cmdUserArgs) > 0 {
				if cmdUserArgs[0] == "holdings" {
					flags = []string{
						"--config",
						configPath,
					}
				}
			}

			for _, arg := range cmdUserArgs {
				flags = append(flags, arg)
			}

			cmd := exec.CommandContext(cmdCtx, s.executableBinary, flags...)
			cmd.Env = append(sshSession.Environ(), fmt.Sprintf("TERM=%s", ptyReq.Term))

			f, err := pty.Start(cmd)
			if err != nil {
				io.WriteString(sshSession, err.Error())
			}

			defer f.Close()

			go func() {
				for win := range winCh {
					setWinsize(f, win.Width, win.Height)
				}
			}()

			go func() {
				io.Copy(f, sshSession)
			}()

			io.Copy(sshSession, f)
			f.Close()
			cmd.Wait()
		},
		PtyCallback: func(ctx ssh.Context, pty ssh.Pty) bool {
			return true
		},
		PublicKeyHandler: func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true
		},
		PasswordHandler: func(ctx ssh.Context, password string) bool {
			return true
		},
		KeyboardInteractiveHandler: func(ctx ssh.Context, challenger gossh.KeyboardInteractiveChallenge) bool {
			return true
		},
	}

	if _, err := os.Stat(s.hostKeyFile); os.IsNotExist(err) {
		return errors.New("SSH key is required to start server")
	}

	err := s.sshServer.SetOption(ssh.HostKeyFile(s.hostKeyFile))
	if err != nil {
		return err
	}

	return s.sshServer.ListenAndServe()
}

// Shutdown shuts down the server
func (s *Server) Shutdown() {
	s.sshServer.Close()
}

// setWinsize sets the PTY window size
func setWinsize(f *os.File, w, h int) {
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
}

// createTempDir creates a temporary directory
func createTempDir() (string, error) {
	return ioutil.TempDir("", "")
}

func validateUserConfigType(userConfigType string) {
	for _, validType := range UserConfigTypes {
		if validType == userConfigType {
			return
		}
	}
	panic(fmt.Errorf("invalid user config type. Acceptable values are: %s", strings.Join(UserConfigTypes, ",")))
}
