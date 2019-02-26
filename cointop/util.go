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

func (ct *Cointop) openLink() error {
	open.URL(ct.rowLink())
	return nil
}

func getBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func userHomeDir() string {
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

func normalizePath(path string) string {
	// expand tilde
	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(userHomeDir(), path[2:])
	}

	return path
}

func (ct *Cointop) slugify(s string) string {
	s = strings.ToLower(s)
	return s
}
