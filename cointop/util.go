package cointop

import (
	"bytes"
	"encoding/gob"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/miguelmota/cointop/cointop/common/open"
)

// OpenLink opens the url in a browser
func (ct *Cointop) OpenLink() error {
	ct.debuglog("openLink()")
	open.URL(ct.RowLink())
	return nil
}

// GetBytes returns the interface in bytes form
func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UserPreferredHomeDir returns the preferred home directory for the user
func UserPreferredHomeDir() string {
	var home string

	if runtime.GOOS == "windows" {
		home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	} else if runtime.GOOS == "linux" {
		home = os.Getenv("XDG_CONFIG_HOME")
	}

	if home == "" {
		home, _ = os.UserHomeDir()
	}

	return home
}

// NormalizePath normalizes and extends the path string
func NormalizePath(path string) string {
	// expand tilde
	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(UserPreferredHomeDir(), path[2:])
	}

	path = strings.Replace(path, "/", string(filepath.Separator), -1)

	return path
}

// Slugify returns a slugified string
func Slugify(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	return s
}
