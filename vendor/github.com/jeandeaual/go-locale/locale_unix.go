//go:build !windows && !darwin && !js && !android
// +build !windows,!darwin,!js,!android

package locale

import (
	"os"
	"strings"
)

func splitLocales(locales string) []string {
	// If the user set different locales, they might be set in $LANGUAGE,
	// separated by a colon
	return strings.Split(locales, ":")
}

func getLangFromEnv() string {
	locale := ""

	// Check the following environment variables for the language information
	// See https://www.gnu.org/software/gettext/manual/html_node/Locale-Environment-Variables.html
	for _, env := range [...]string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		locale = os.Getenv(env)
		if len(locale) > 0 {
			break
		}
	}

	if locale == "C" || locale == "POSIX" {
		return locale
	}

	// Check LANGUAGE if localization is enabled (not set to "C")
	// See https://www.gnu.org/software/gettext/manual/html_node/The-LANGUAGE-variable.html#The-LANGUAGE-variable
	languages := os.Getenv("LANGUAGE")
	if len(languages) > 0 {
		return languages
	}

	return locale
}

func getUnixLocales() []string {
	locale := getLangFromEnv()
	if locale == "C" || locale == "POSIX" || len(locale) == 0 {
		return nil
	}

	return splitLocales(locale)
}

// GetLocale retrieves the IETF BCP 47 language tag set on the system.
func GetLocale() (string, error) {
	unixLocales := getUnixLocales()
	if unixLocales == nil {
		return "", nil
	}

	language, region := splitLocale(unixLocales[0])
	locale := language
	if len(region) > 0 {
		locale = strings.Join([]string{language, region}, "-")
	}

	return locale, nil
}

// GetLocales retrieves the IETF BCP 47 language tags set on the system.
func GetLocales() ([]string, error) {
	unixLocales := getUnixLocales()
	if unixLocales == nil {
		return nil, nil
	}

	locales := make([]string, 0, len(unixLocales))

	for _, unixLocale := range unixLocales {
		language, region := splitLocale(unixLocale)
		locale := language
		if len(region) > 0 {
			locale = strings.Join([]string{language, region}, "-")
		}
		locales = append(locales, locale)
	}

	return locales, nil
}

// GetLanguage retrieves the IETF BCP 47 language tag set on the system and
// returns the language part of the tag.
func GetLanguage() (string, error) {
	language := ""

	unixLocales := getUnixLocales()
	if unixLocales == nil {
		return "", nil
	}

	language, _ = splitLocale(unixLocales[0])

	return language, nil
}

// GetRegion retrieves the IETF BCP 47 language tag set on the system and
// returns the region part of the tag.
func GetRegion() (string, error) {
	region := ""

	unixLocales := getUnixLocales()
	if unixLocales == nil {
		return "", nil
	}

	_, region = splitLocale(unixLocales[0])

	return region, nil
}
