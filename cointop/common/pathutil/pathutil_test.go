package pathutil

import (
	"os"
	"path/filepath"
	"testing"
)

// TestNormalizePath checks that NormalizePath returns the correct directory
func TestNormalizePath(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"~/.config/cointop/config.toml", filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "/cointop/config.toml")},
	}
	for _, c := range cases {
		got := NormalizePath(c.in)
		if got != c.want {
			t.Errorf("NormalizePath(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
