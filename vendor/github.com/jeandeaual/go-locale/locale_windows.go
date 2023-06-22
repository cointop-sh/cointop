//go:build windows
// +build windows

package locale

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// LocaleNameMaxLength is the maximum length of a locale name on Windows.
// See https://docs.microsoft.com/en-us/windows/win32/intl/locale-name-constants.
const LocaleNameMaxLength uint32 = 85

func getWindowsLocaleFromProc(syscall string) (string, error) {
	dll, err := windows.LoadDLL("kernel32")
	if err != nil {
		return "", fmt.Errorf("could not find the kernel32 DLL: %v", err)
	}

	proc, err := dll.FindProc(syscall)
	if err != nil {
		return "", fmt.Errorf("could not find the %s proc in kernel32: %v", syscall, err)
	}

	buffer := make([]uint16, LocaleNameMaxLength)

	// See https://docs.microsoft.com/en-us/windows/win32/api/winnls/nf-winnls-getuserdefaultlocalename
	// and https://docs.microsoft.com/en-us/windows/win32/api/winnls/nf-winnls-getsystemdefaultlocalename
	// GetUserDefaultLocaleName and GetSystemDefaultLocaleName both take a buffer and a buffer size,
	// and return the length of the locale name (0 if not found).
	ret, _, err := proc.Call(uintptr(unsafe.Pointer(&buffer[0])), uintptr(LocaleNameMaxLength))
	if ret == 0 {
		return "", fmt.Errorf("locale not found when calling %s: %v", syscall, err)
	}

	return windows.UTF16ToString(buffer), nil
}

func getWindowsLocale() (string, error) {
	var (
		locale string
		err    error
	)

	for _, proc := range [...]string{"GetUserDefaultLocaleName", "GetSystemDefaultLocaleName"} {
		locale, err = getWindowsLocaleFromProc(proc)
		if err == nil {
			return locale, err
		}
	}

	return locale, err
}

// GetLocale retrieves the IETF BCP 47 language tag set on the system.
func GetLocale() (string, error) {
	locale, err := getWindowsLocale()
	if err != nil {
		return "", fmt.Errorf("cannot determine locale: %v", err)
	}

	return locale, err
}

// GetLocales retrieves the IETF BCP 47 language tags set on the system.
func GetLocales() ([]string, error) {
	locale, err := GetLocale()
	if err != nil {
		return nil, err
	}

	return []string{locale}, nil
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
