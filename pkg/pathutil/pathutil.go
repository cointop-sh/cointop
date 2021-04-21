package pathutil

import (
	"os"
	"path/filepath"
	"strings"
)

// UserPreferredConfigDir returns the preferred config directory for the user
func UserPreferredConfigDir() string {
	defaultConfigDir := "~/.config"

	config, err := os.UserConfigDir()
	if err != nil {
		return defaultConfigDir
	}

	if config == "" {
		return defaultConfigDir
	}

	return config
}

// UserPreferredCacheDir returns the preferred cache directory for the user
func UserPreferredCacheDir() string {
	defaultCacheDir := "/tmp"

	cache, err := os.UserCacheDir()
	if err != nil {
		return defaultCacheDir
	}

	if cache == "" {
		return defaultCacheDir
	}

	return cache
}

// UserPreferredHomeDir returns the preferred home directory for the user
func UserPreferredHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return home
}

// NormalizePath normalizes and extends the path string
func NormalizePath(path string) string {
	userHome := UserPreferredHomeDir()
	userConfigHome := UserPreferredConfigDir()
	userCacheHome := UserPreferredCacheDir()

	// expand tilde
	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(userHome, path[2:])
	}

	path = strings.Replace(path, ":HOME:", userHome, -1)
	path = strings.Replace(path, ":PREFERRED_CONFIG_HOME:", userConfigHome, -1)
	path = strings.Replace(path, ":PREFERRED_CACHE_HOME:", userCacheHome, -1)
	path = strings.Replace(path, "/", string(filepath.Separator), -1)

	return filepath.Clean(path)
}
