package humanize

import (
	"testing"
)

// TestMonetary tests monetary formatting
func TestMonetary(t *testing.T) {
	if Monetaryf(834142.3256, 2) != "834,142.3256" {
		t.FailNow()
	}
}
