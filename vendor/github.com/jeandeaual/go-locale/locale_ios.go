//go:build ios
// +build ios

package locale

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework UIKit

const char *getLocale();
const char *getLocales();
*/
import "C"
import (
	"strings"
)

// GetLocale retrieves the IETF BCP 47 language tag set on the system.
func GetLocale() (string, error) {
	return C.GoString(C.getLocale()), nil
}

// GetLocales retrieves the IETF BCP 47 language tags set on the system.
func GetLocales() ([]string, error) {
	return strings.Split(C.GoString(C.getLocales()), ","), nil
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
