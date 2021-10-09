# go-locale

[![PkgGoDev](https://pkg.go.dev/badge/github.com/jeandeaual/go-locale)](https://pkg.go.dev/github.com/jeandeaual/go-locale)
[![Go Report Card](https://goreportcard.com/badge/github.com/jeandeaual/go-locale)](https://goreportcard.com/report/github.com/jeandeaual/go-locale)
[![Coverage Status](https://coveralls.io/repos/github/jeandeaual/go-locale/badge.svg?branch=master)](https://coveralls.io/github/jeandeaual/go-locale?branch=master)
[![test](https://github.com/jeandeaual/go-locale/workflows/test/badge.svg)](https://github.com/jeandeaual/go-locale/actions?query=workflow%3Atest)

Go library used to retrieve the current locale(s) of the operating system.

## OS Support

* Windows\
    Using [`GetUserDefaultLocaleName`](https://docs.microsoft.com/en-us/windows/win32/api/winnls/nf-winnls-getuserdefaultlocalename) and [`GetSystemDefaultLocaleName`](https://docs.microsoft.com/en-us/windows/win32/api/winnls/nf-winnls-getsystemdefaultlocalename).
* macOS\
    Using `defaults read -g AppleLocale` and `defaults read -g AppleLanguages` (since environment variables like `LANG` are not usually set on macOS).
* Unix-like systems (Linux, BSD, etc.)\
    Using the `LANGUAGE`, `LC_ALL`, `LC_MESSAGES` and `LANG` environment variables.
* WASM (JavaScript)\
    Using [`navigator.language`](https://developer.mozilla.org/en-US/docs/Web/API/NavigatorLanguage/language) and [`navigator.languages`](https://developer.mozilla.org/en-US/docs/Web/API/NavigatorLanguage/languages).
* iOS\
    Using [`[NSLocale preferredLanguages]`](https://developer.apple.com/documentation/foundation/nslocale/1415614-preferredlanguages).
* Android\
    Using [`getResources().getConfiguration().getLocales`](https://developer.android.com/reference/android/content/res/Configuration#getLocales()) for Android N or later, or [`getResources().getConfiguration().locale`](https://developer.android.com/reference/android/content/res/Configuration#locale) otherwise.

    Note: for Android, you'll first need to call `SetRunOnJVM`, depending on which mobile framework you're using:
    * For [Fyne](https://fyne.io/):

        ```go
        import (
        	"github.com/fyne-io/mobile/app"
        	"github.com/jeandeaual/go-locale"
        )

        func init() {
        	locale.SetRunOnJVM(app.RunOnJVM)
        }
        ```

    * For [gomobile](https://github.com/golang/go/wiki/Mobile):

        ```go
        import (
        	"golang.org/x/mobile/app"
        	"github.com/jeandeaual/go-locale"
        )

        func init() {
        	locale.SetRunOnJVM(app.RunOnJVM)
        }
        ```

## Usage

## GetLocales

`GetLocales` returns the user's preferred locales, by order of preference, as a slice of [IETF BCP 47 language tag](https://tools.ietf.org/rfc/bcp/bcp47.txt) (e.g. `[]string{"en-US", "fr-FR", "ja-JP"}`).

This works if the user set multiple languages on macOS and other Unix systems.
Otherwise, it returns a slice with a single locale.

```go
userLocales, err := locale.GetLocales()
if err == nil {
	fmt.Println("Locales:", userLocales)
}
```

This can be used with [golang.org/x/text](https://godoc.org/golang.org/x/text) or [go-i18n](https://github.com/nicksnyder/go-i18n) to set the localizer's language preferences:

```go
import (
	"github.com/jeandeaual/go-locale"
	"golang.org/x/text/message"
)

func main() {
	userLocales, _ := locale.GetLocales()
	p := message.NewPrinter(message.MatchLanguage(userLocales...))
	...
}
```

```go
import (
	"github.com/jeandeaual/go-locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func main() {
	userLocales, _ := locale.GetLocales()
	bundle := i18n.NewBundle(language.English)
	localizer := i18n.NewLocalizer(bundle, userLocales...)
	...
}
```

For a complete example, see [here](examples/getlocale-gui/main.go).

## GetLocale

`GetLocale` returns the current locale as defined in [IETF BCP 47](https://tools.ietf.org/rfc/bcp/bcp47.txt) (e.g. `"en-US"`).

```go
userLocale, err := locale.GetLocale()
if err == nil {
	fmt.Println("Locale:", userLocale)
}
```

## GetLanguage

`GetLanguage` returns the current language as an [ISO 639](http://en.wikipedia.org/wiki/ISO_639) language code (e.g. `"en"`).

```go
userLanguage, err := locale.GetLanguage()
if err == nil {
	fmt.Println("Language:", userLocale)
}
```

## GetRegion

`GetRegion` returns the current language as an [ISO 3166](http://en.wikipedia.org/wiki/ISO_3166-1) country code (e.g. `"US"`).

```go
userRegion, err := locale.GetRegion()
if err == nil {
	fmt.Println("Region:", userRegion)
}
```

## Aknowledgements

Inspired by [jibber_jabber](https://github.com/cloudfoundry-attic/jibber_jabber).
