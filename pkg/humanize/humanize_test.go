package humanize

import (
	"fmt"
	"testing"
)

// TestMonetary tests monetary formatting
func TestMonetary(t *testing.T) {
	if Monetaryf(834142.3256, 2) != "834,142.3256" {
		t.FailNow()
	}

	scaleTests := map[float64]string {
		5.5 * 1e12: "5.5T",
		4.4 * 1e9: "4.4B",
		3.3 * 1e6: "3.3M",
		2.2 * 1e3: "2200.0",
		1.1: "1.1",
		0.06: "0.1",
		-5.5 * 1e12: "-5.5T",
		-4.4 * 1e9: "-4.4B",
		-3.3 * 1e6: "-3.3M",
	}

	for value, expected := range scaleTests {
		volScale, volSuffix := Scale(value)
		result := fmt.Sprintf("%.1f%s", volScale, volSuffix)
		if result != expected {
			t.Fatalf("Expected %f to scale to '%s' but got '%s'\n", value, expected, result)
		}
	}
}
