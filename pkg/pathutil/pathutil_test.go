package pathutil

import (
	"os"
	"path/filepath"
	"testing"
)

// TestNormalizePath checks that NormalizePath returns the correct directory
func TestNormalizePath(t *testing.T) {
	os.Setenv("XDG_CONFIG_HOME", "/tmp/.config")

	home, _ := os.UserHomeDir()
	configDir, _ := os.UserConfigDir()

	cases := []struct {
		in, out string
	}{
		{"~/.config/cointop/config.toml", filepath.Join(home, ".config/cointop/config.toml")},
		{":HOME:/.cointop/config.toml", filepath.Join(home, "/.cointop/config.toml")},
		{":PREFERRED_CONFIG_HOME:/cointop/config.toml", filepath.Join(configDir, "/cointop/config.toml")},
	}
	for _, c := range cases {
		got := NormalizePath(c.in)
		if got != c.out {
			t.Errorf("NormalizePath(%q) == %q, want %q", c.in, got, c.out)
		}
	}
}
