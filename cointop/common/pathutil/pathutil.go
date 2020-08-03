package pathutil

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

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
