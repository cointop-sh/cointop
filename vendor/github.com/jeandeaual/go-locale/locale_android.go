//go:build android
// +build android

package locale

/*
#cgo LDFLAGS: -landroid -llog

#include <stdlib.h>

const char *getLocales(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx);
const char *getLocale(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx);
const char *getLanguage(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx);
const char *getRegion(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx);
*/
import "C"
import (
	"errors"
	"strings"
	"unsafe"
)

var (
	errRunOnJVMNotSet error = errors.New("you first need to call SetRunOnJVM")
	runOnJVM          func(fn func(vm, env, ctx uintptr) error) error
)

// SetRunOnJVM sets the RunOnJVM function that will be called by this library.
// This can either be "golang.org/x/mobile/app".RunOnJVM or "github.com/fyne-io/mobile/app".RunOnJVM,
// depending on the mobile framework you're using (both can't be imported at the same time).
//
// RunOnJVM runs fn on a new goroutine locked to an OS thread with a JNIEnv.
//
// RunOnJVM blocks until the call to fn is complete. Any Java
// exception or failure to attach to the JVM is returned as an error.
//
// The function fn takes vm, the current JavaVM*,
// env, the current JNIEnv*, and
// ctx, a jobject representing the global android.context.Context.
func SetRunOnJVM(fn func(fn func(vm, env, ctx uintptr) error) error) {
	runOnJVM = fn
}

// GetLocale retrieves the IETF BCP 47 language tag set on the system.
func GetLocale() (string, error) {
	if runOnJVM == nil {
		return "", errRunOnJVMNotSet
	}

	locale := ""

	err := runOnJVM(func(vm, env, ctx uintptr) error {
		chars := C.getLocale(C.uintptr_t(vm), C.uintptr_t(env), C.uintptr_t(ctx))
		locale = C.GoString(chars)
		C.free(unsafe.Pointer(chars))
		return nil
	})

	return locale, err
}

// GetLocales retrieves the IETF BCP 47 language tags set on the system.
func GetLocales() ([]string, error) {
	if runOnJVM == nil {
		return nil, errRunOnJVMNotSet
	}

	locales := ""

	err := runOnJVM(func(vm, env, ctx uintptr) error {
		chars := C.getLocales(C.uintptr_t(vm), C.uintptr_t(env), C.uintptr_t(ctx))
		locales = C.GoString(chars)
		C.free(unsafe.Pointer(chars))
		return nil
	})

	return strings.Split(locales, ","), err
}

// GetLanguage retrieves the IETF BCP 47 language tag set on the system and
// returns the language part of the tag.
func GetLanguage() (string, error) {
	if runOnJVM == nil {
		return "", errRunOnJVMNotSet
	}

	language := ""

	err := runOnJVM(func(vm, env, ctx uintptr) error {
		chars := C.getLanguage(C.uintptr_t(vm), C.uintptr_t(env), C.uintptr_t(ctx))
		language = C.GoString(chars)
		C.free(unsafe.Pointer(chars))
		return nil
	})

	return language, err
}

// GetRegion retrieves the IETF BCP 47 language tag set on the system and
// returns the region part of the tag.
func GetRegion() (string, error) {
	if runOnJVM == nil {
		return "", errRunOnJVMNotSet
	}

	region := ""

	err := runOnJVM(func(vm, env, ctx uintptr) error {
		chars := C.getRegion(C.uintptr_t(vm), C.uintptr_t(env), C.uintptr_t(ctx))
		region = C.GoString(chars)
		C.free(unsafe.Pointer(chars))
		return nil
	})

	return region, err
}
