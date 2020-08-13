package pathutil

import (
	"os"
	"path/filepath"
	"strings"
)

// UserPreferredHomeDir returns the preferred home directory for the user
func UserPreferredHomeDir() (string, bool) {
	var isConfigDir bool

	home, _ := os.UserConfigDir()
	isConfigDir = true

	if home == "" {
		home, _ = os.UserHomeDir()
		isConfigDir = false
	}

	return home, isConfigDir
}

// NormalizePath normalizes and extends the path string
func NormalizePath(path string) string {
	// expand tilde
	if strings.HasPrefix(path, "~/") {
		home, isConfigDir := UserPreferredHomeDir()
		if !isConfigDir {
			path = filepath.Join(home, path[2:])
		}
		path = filepath.Join(home, path[10:])
	}

	path = strings.Replace(path, "/", string(filepath.Separator), -1)

	return path
}
