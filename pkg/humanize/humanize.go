package humanize

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/goodsign/monday"
	"github.com/jeandeaual/go-locale"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var cachedSystemLocale = ""

// Numericf produces a string from of the given number with give fixed precision
// in base 10 with thousands separators after every three orders of magnitude
// using thousands and decimal separator according to LC_NUMERIC; defaulting "en".
//
// e.g. Numericf(834142.32, 2) -> "834,142.32"
func Numericf(value float64, precision int) string {
	return f(value, precision, "LC_NUMERIC", true)
}

// Monetaryf produces a string from of the given number give minimum precision
// in base 10 with thousands separators after every three orders of magnitude
// using thousands and decimal separator according to LC_MONETARY; defaulting "en".
//
// e.g. Monetaryf(834142.3256, 2) -> "834,142.3256"
func Monetaryf(value float64, precision int) string {
	return f(value, precision, "LC_MONETARY", false)
}

// borrowed from go-locale/util.go
func splitLocale(locale string) (string, string) {
	// Remove the encoding, if present
	formattedLocale := strings.Split(locale, ".")[0]
	// Normalize by replacing the hyphens with underscores
	formattedLocale = strings.Replace(formattedLocale, "-", "_", -1)

	// Split at the underscore
	split := strings.Split(formattedLocale, "_")
	language := split[0]
	territory := ""
	if len(split) > 1 {
		territory = split[1]
	}

	return language, territory
}

// GetLocale returns the current locale as defined in IETF BCP 47 (e.g. "en-US").
// The envvar provided is checked first (eg LC_TIME), before the platform-specific defaults.
func getLocale(envvar string) string {
	userLocale := "en-US" // default language-REGION
	// First try looking up envar directly
	envlang, ok := os.LookupEnv(envvar)
	if ok {
		language, region := splitLocale(envlang)
		userLocale = language
		if len(region) > 0 {
			userLocale = strings.Join([]string{language, region}, "-")
		}
	} else {
		// Then use (cached) system-specific locale
		if cachedSystemLocale == "" {
			if loc, err := locale.GetLocale(); err == nil {
				userLocale = loc
				cachedSystemLocale = loc
			}
		} else {
			userLocale = cachedSystemLocale
		}
	}
	return userLocale
}

// formatTimeExplicit formats the given time using the prescribed layout with the provided userLocale
func formatTimeExplicit(time time.Time, layout string, userLocale string) string {
	mondayLocale := monday.Locale(strings.Replace(userLocale, "-", "_", 1))
	return monday.Format(time, layout, mondayLocale)
}

// FormatTime is a dropin replacement time.Format(layout) that uses system locale + LC_TIME
func FormatTime(time time.Time, layout string) string {
	return formatTimeExplicit(time, layout, getLocale("LC_TIME"))
}

// f formats given value, with precision decimal places using thousands and decimal
// separator according to language found in given locale environment variable e.
// If fixed is true the decimal places are fixed to the given precision otherwise d is the
// minimum of decimal places until the first 0.
func f(value float64, precision int, envvar string, fixed bool) string {
	parts := strings.Split(strconv.FormatFloat(value, 'f', -1, 64), ".")
	if !fixed && len(parts) > 1 {
		for ; precision < len(parts[1]); precision += 1 {
			if parts[1][precision] == '0' {
				break
			}
		}
	}

	envlang, ok := os.LookupEnv(envvar)
	if !ok {
		envlang = "en"
	}
	lang := language.Make(envlang)

	format := fmt.Sprintf("%%.%df", precision)
	return message.NewPrinter(lang).Sprintf(format, value)
}

// Scale returns a scaled-down version of value and a suffix to add (M,B,etc.)
func Scale(value float64) (float64, string) {
	type scalingUnit struct {
		value  float64
		suffix string
	}

	//  quadrillion, quintrillion, sextillion, septillion, octillion, nonillion, and decillion
	var scales = [...]scalingUnit{
		{value: 1e12, suffix: "T"},
		{value: 1e9, suffix: "B"},
		{value: 1e6, suffix: "M"},
		{value: 1e3, suffix: "K"},
	}

	for _, scale := range scales {
		if math.Abs(value) > scale.value {
			return value / scale.value, scale.suffix
		}
	}
	return value, ""
}

// ScaleNumericf scales a large number down using a suffix, then formats it with the
// prescribed number of significant digits.
func ScaleNumericf(value float64, digits int) string {
	value, suffix := Scale(value)

	// Round the scaled value to a certain number of significant figures
	var s string
	if math.Abs(value) < 1 {
		s = Numericf(value, digits)
	} else {
		numDigits := len(fmt.Sprintf("%.0f", math.Abs(value)))
		if numDigits >= digits {
			s = Numericf(value, 0)
		} else {
			s = Numericf(value, digits-numDigits)
		}
	}

	return s + suffix
}
