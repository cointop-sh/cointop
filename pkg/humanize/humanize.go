package humanize

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Numericf produces a string from of the given number with give fixed precision
// in base 10 with thousands separators after every three orders of magnitude
// using a thousands and decimal spearator according to LC_NUMERIC; defaulting "en".
//
// e.g. Numericf(834142.32, 2) -> "834,142.32"
func Numericf(value float64, precision int) string {
	return f(value, precision, "LC_NUMERIC", true)
}

// Monetaryf produces a string from of the given number give minimum precision
// in base 10 with thousands separators after every three orders of magnitude
// using thousands and decimal spearator according to LC_MONETARY; defaulting "en".
//
// e.g. Monetaryf(834142.3256, 2) -> "834,142.3256"
func Monetaryf(value float64, precision int) string {
	return f(value, precision, "LC_MONETARY", false)
}

// f formats given value v, with d decimal places using thousands and decimal
// separator according to language found in given locale environment variable e.
// If r is true the decimal places are fixed to the given d otherwise d is the
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
