//go:build darwin && !ios
// +build darwin,!ios

package locale

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
)

func execCommand(cmd string, args ...string) (status int, out string, err error) {
	var bytesOut []byte
	status = -1
	command := exec.Command(cmd, args...)

	// Execute the command and get the standard and error outputs
	bytesOut, err = command.CombinedOutput()
	out = string(bytesOut)
	if err != nil {
		return
	}

	// Check the status code
	if w, ok := command.ProcessState.Sys().(syscall.WaitStatus); ok {
		status = w.ExitStatus()
	}

	return
}

// GetLocale retrieves the IETF BCP 47 language tag set on the system.
func GetLocale() (string, error) {
	_, output, err := execCommand("defaults", "read", "-g", "AppleLocale")
	if err != nil {
		return "", fmt.Errorf("cannot determine locale: %v (output: %s)", err, output)
	}

	// defaults read -g AppleLocale can return a string containing additional
	// information after the locale, e.g. "en_US@currency=USD"
	if idx := strings.Index(output, "@"); idx != -1 {
		output = output[:idx]
	}

	return strings.TrimRight(strings.Replace(output, "_", "-", 1), "\n"), nil
}

// appleLanguagesRegex is used to parse the output of "defaults read -g AppleLanguages"
// e.g.:
// (
//     en,
//     "fr-FR",
//     "ja-JP"
// )
var appleLanguagesRegex = regexp.MustCompile(`([a-z]{2}(?:-[A-Z]{2})?)`)

// GetLocales retrieves the IETF BCP 47 language tags set on the system.
func GetLocales() ([]string, error) {
	_, output, err := execCommand("defaults", "read", "-g", "AppleLanguages")
	if err != nil {
		return nil, fmt.Errorf("cannot determine locale: %v (output: %s)", err, output)
	}

	matches := appleLanguagesRegex.FindAllStringSubmatch(output, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("invalid output from \"defaults read -g AppleLanguages\": %s", output)
	}

	locales := make([]string, 0, len(matches))

	for _, match := range matches {
		locales = append(locales, match[1])
	}

	return locales, nil
}

// GetLanguage retrieves the IETF BCP 47 language tag set on the system and
// returns the language part of the tag.
func GetLanguage() (string, error) {
	language := ""

	locale, err := GetLocale()
	if err == nil {
		language, _ = splitLocale(locale)
	}

	return language, err
}

// GetRegion retrieves the IETF BCP 47 language tag set on the system and
// returns the region part of the tag.
func GetRegion() (string, error) {
	region := ""

	locale, err := GetLocale()
	if err == nil {
		_, region = splitLocale(locale)
	}

	return region, err
}
