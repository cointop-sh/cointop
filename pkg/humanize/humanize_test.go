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
	testData := map[string]map[string]string{
		"en_GB": {
			"Monday 2 January 2006": "Wednesday 12 March 2014",
			"Jan 2006":              "Mar 2014",
			"02 Jan 2006":           "12 Mar 2014",
			"02/01/2006":            "12/03/2014",
		},
		"en_US": {
			"Monday 2 January 2006": "Wednesday 12 March 2014",
			"Jan 2006":              "Mar 2014",
			"02 Jan 2006":           "12 Mar 2014",
			"02/01/2006":            "12/03/2014", // ??
		},
		"fr_FR": {
			"Monday 2 January 2006": "mercredi 12 mars 2014",
			"Jan 2006":              "mars 2014",
			"02 Jan 2006":           "12 mars 2014",
			"02/01/2006":            "12/03/2014",
		},
		"de_DE": {
			"Monday 2 January 2006": "Mittwoch 12 März 2014",
			"Jan 2006":              "Mär 2014",
			"02 Jan 2006":           "12 Mär 2014",
			"02/01/2006":            "12/03/2014",
		},
	}

	testTime := time.Date(2014, 3, 12, 0, 0, 0, 0, time.Local)
	for locale, tests := range testData {
		for layout, result := range tests {
			s := formatTimeExplicit(testTime, layout, locale)
			if s != result {
				t.Fatalf("Expected layout '%s' in locale %s to render '%s' but got '%s'", layout, locale, result, s)
			}

		}
	}
}
