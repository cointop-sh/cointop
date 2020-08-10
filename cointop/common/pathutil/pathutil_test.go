package pathutil

import (
	"os"
	"path/filepath"
	"testing"
)

// TestNormalizePath checks that NormalizePath returns the correct directory
func TestNormalizePath(t *testing.T) {
	home, _ := os.UserHomeDir()
	configDir, _ := os.UserConfigDir()
	cases := []struct {
		in, want string
	}{
		{"~/.config/cointop/config.toml", filepath.Join(configDir, "/cointop/config.toml")},
		{"~/.config/cointop/config.toml", filepath.Join(home, ".config/cointop/config.toml")},
		{"~/.config/cointop/config.toml", filepath.Join(configDir, "/cointop/config.toml")},
		{"~/.config/cointop/config.toml", filepath.Join(home, ".config/cointop/config.toml")},
	}
	for i, c := range cases {
		got := NormalizePath(c.in)
		if i > 1 {
			home = ""
			configDir = ""
		}
		if got != c.want {
			t.Errorf("NormalizePath(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
