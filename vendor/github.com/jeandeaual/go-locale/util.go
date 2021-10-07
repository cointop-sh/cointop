// +build !android

package locale

import (
	"strings"
)

// SetRunOnJVM is a noop, this function is only valid on Android
func SetRunOnJVM(fn func(fn func(vm, env, ctx uintptr) error) error) {}

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
