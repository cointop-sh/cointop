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

// UserHomeDir returns home directory for the user
func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	} else if runtime.GOOS == "linux" {
		home := os.Getenv("XDG_CONFIG_HOME")
		if home != "" {
			return home
		}
	}
	return os.Getenv("HOME")
}

// NormalizePath normalizes and extends the path string
func NormalizePath(path string) string {
	// expand tilde
	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(UserHomeDir(), path[2:])
	}

	return path
}

// Slugify returns a slugified string
func (ct *Cointop) Slugify(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	return s
}
