package humanize

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/goodsign/monday"
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

// Attempt to determine the locale from the current environment. If usage is provided, treat it as the name of a LC_xxx environment variable.
// LANGUAGE
// LC_ALL       Will override the setting of all other LC_* variables.
// LC_MONETARY  Sets the locale for the LC_MONETARY category.
// LC_NUMERIC   Sets the locale for the LC_NUMERIC category.
// LC_TIME      Sets the locale for the LC_TIME category.
// LANG         Used as a substitute for any unset LC_* variable.  If LANG is unset, it will act as if set to "C"
// Local is language[_territory][.codeset] [@modifier]
func DetectLocale(usage string) string {
	if lc, ok := os.LookupEnv("LANGUAGE"); ok {
		return lc
	}
	if lc, ok := os.LookupEnv("LC_ALL"); ok {
		return lc
	}
	if usage != "" {
		if lc, ok := os.LookupEnv(strings.ToUpper(usage)); ok {
			return lc
		}
	}
	if lc, ok := os.LookupEnv("LANG"); ok {
		return lc
	}
	return "C"
}

// func DetectLanguage(usage string) language.Tag {
// 	lc := DetectLocale(usage)
// 	if lc == "C" {
// 		lc = "en"
// 	}
// 	return language.Make(lc)
// }

// Locale is language[_territory][.codeset] [@modifier]
func FormatTime(time time.Time, layout string) string {
	// Attempt to use the environment to determine monday.Locale
	xxx := strings.Split(DetectLocale("LC_TIME"), ".")[0]
	bits := strings.Split(xxx, "_")

	// Look for a supported Locale with default to en_US
	var locale monday.Locale = monday.LocaleEnUS // default
	if len(bits) == 2 {
		lookFor := monday.Locale(strings.ToLower(bits[0]) + "_" + strings.ToUpper(bits[1]))
		for _, v := range monday.ListLocales() {
			if v == lookFor {
				locale = v
			}
		}
	}

	return monday.Format(time, layout, locale)
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
