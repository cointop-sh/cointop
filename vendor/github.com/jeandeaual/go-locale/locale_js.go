//go:build js && wasm
// +build js,wasm

package locale

import (
	"errors"
	"syscall/js"
)

func getNavigatorObject() js.Value {
	return js.Global().Get("navigator")
}

// GetLocale retrieves the IETF BCP 47 language tag set on the system.
func GetLocale() (string, error) {
	navigator := getNavigatorObject()
	if navigator.IsUndefined() {
		return "", errors.New("couldn't get window.navigator")
	}

	language := navigator.Get("language")
	if language.IsUndefined() {
		return "", errors.New("couldn't get window.navigator.language")
	}

	return language.String(), nil
}

// GetLocales retrieves the IETF BCP 47 language tags set on the system.
func GetLocales() ([]string, error) {
	navigator := getNavigatorObject()
	if navigator.IsUndefined() {
		return nil, errors.New("couldn't get window.navigator")
	}

	languages := navigator.Get("languages")
	if languages.IsUndefined() {
		return nil, errors.New("couldn't get window.navigator.languages")
	}

	locales := make([]string, 0, languages.Length())

	// Convert the Javascript object to a string slice
	for i := 0; i < languages.Length(); i++ {
		locales = append(locales, languages.Index(i).String())
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
