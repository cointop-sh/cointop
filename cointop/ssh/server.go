//+build !windows

package ssh

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
	"time"
	"unsafe"

	"github.com/creack/pty"
	"github.com/gliderlabs/ssh"
	"github.com/miguelmota/cointop/cointop/common/pathutil"
	gossh "golang.org/x/crypto/ssh"
)

// DefaultHostKeyFile ...
var DefaultHostKeyFile = "~/.ssh/id_rsa"

// Config ...
type Config struct {
	Port             uint
	Address          string
	IdleTimeout      time.Duration
	ExecutableBinary string
	HostKeyFile      string
}

// Server ...
type Server struct {
	port             uint
	address          string
	idleTimeout      time.Duration
	executableBinary string
	sshServer        *ssh.Server
	hostKeyFile      string
}

// NewServer ...
func NewServer(config *Config) *Server {
	hostKeyFile := DefaultHostKeyFile
	if config.HostKeyFile != "" {
		hostKeyFile = config.HostKeyFile
	}

	hostKeyFile = pathutil.NormalizePath(hostKeyFile)

	return &Server{
		port:             config.Port,
		address:          config.Address,
		idleTimeout:      config.IdleTimeout,
		executableBinary: config.ExecutableBinary,
		hostKeyFile:      hostKeyFile,
	}
}

// ListenAndServe ...
func (s *Server) ListenAndServe() error {
	s.sshServer = &ssh.Server{
		Addr:        fmt.Sprintf("%s:%v", s.address, s.port),
		IdleTimeout: s.idleTimeout,
		Handler: func(sshSession ssh.Session) {
			cmdUserArgs := sshSession.Command()
			ptyReq, winCh, isPty := sshSession.Pty()
			if !isPty {
				io.WriteString(sshSession, "Error: Non-interactive terminals are not supported")
				sshSession.Exit(1)
				return
			}

			tempDir, err := createTempDir()
			if err != nil {
				fmt.Println(err)
				return
			}

			configPath := fmt.Sprintf("%s/config", tempDir)
			colorsDir := pathutil.NormalizePath("~/.config/cointop/colors")

			cmdCtx, cancelCmd := context.WithCancel(sshSession.Context())
			defer cancelCmd()

			flags := []string{
				"--reset",
				"--silent",
				"--cache-dir",
				tempDir,
				"--config",
				configPath,
				"--colors-dir",
				colorsDir,
			}

			for i, arg := range cmdUserArgs {
				if i == 0 {
					continue
				}

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
			os.Remove(configPath)
		},
		PtyCallback: func(ctx ssh.Context, pty ssh.Pty) bool {
			// TODO: check public key hash
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

// Shutdown ...
func (s *Server) Shutdown() {
	s.sshServer.Close()
}

// setWinsize ...
func setWinsize(f *os.File, w, h int) {
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
}

// createTempDir ...
// TODO: load saved configuration based on ssh public key hash
func createTempDir() (string, error) {
	return ioutil.TempDir("", "")
}
