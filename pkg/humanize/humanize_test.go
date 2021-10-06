package humanize

import (
	"fmt"
	"testing"
	"time"
)

// TestMonetary tests monetary formatting
func TestMonetary(t *testing.T) {
	if Monetaryf(834142.3256, 2) != "834,142.3256" {
		t.FailNow()
	}
}

func TestScale(t *testing.T) {
	scaleTests := map[float64]string{
		5.54 * 1e12:  "5.5T",
		4.44 * 1e9:   "4.4B",
		3.34 * 1e6:   "3.3M",
		2.24 * 1e3:   "2.2K",
		1.1:          "1.1",
		0.06:         "0.1",
		0.04:         "0.0",
		-5.54 * 1e12: "-5.5T",
	}

	for value, expected := range scaleTests {
		volScale, volSuffix := Scale(value)
		result := fmt.Sprintf("%.1f%s", volScale, volSuffix)
		if result != expected {
			t.Fatalf("Expected %f to scale to '%s' but got '%s'\n", value, expected, result)
		}
	}
}

func TestScaleNumeric(t *testing.T) {
	scaleTests := map[float64]string{
		5.54 * 1e12:    "5.5T",
		4.44 * 1e9:     "4.4B",
		3.34 * 1e6:     "3.3M",
		2.24 * 1e3:     "2.2K",
		1.1:            "1.1",
		0.0611:         "0.06",
		-5.5432 * 1e12: "-5.5T",
	}

	for value, expected := range scaleTests {
		result := ScaleNumericf(value, 2)
		if result != expected {
			t.Fatalf("Expected %f to scale to '%s' but got '%s'\n", value, expected, result)
		}
	}
}

func TestFormatTime(t *testing.T) {
	s := FormatTime(time.Now(), "Jan 2006")
	t.Logf("First: %s", s)
	if Monetaryf(834142.3256, 2) != "834,142.3256" {
		t.FailNow()
	}
}
