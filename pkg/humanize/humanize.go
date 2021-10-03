package humanize

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

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
