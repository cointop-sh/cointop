Description
====

Monday is a minimalistic translator for month and day of week names in time.Date objects. Supports 20+ different locales.
Written in pure [Go](http://golang.org).

![Go](https://github.com/goodsign/monday/workflows/Go/badge.svg)

Installing
====

```
go get github.com/goodsign/monday
```

Usage
====

Format
---------------------

Given that you already use [time.Format](http://golang.org/pkg/time/#Time.Format) somewhere in your code,
to translate your output you should import monday and replace

```go
  yourTime.Format(yourLayout)
```

with

```go
  // Change LocaleEnUS to the locale you want to use for translation
  monday.Format(yourTime, yourLayout, monday.LocaleEnUS)
```

Parse
---------------------

Given that you already use [time.ParseInLocation](http://golang.org/pkg/time/#ParseInLocation) somewhere in your code,
to parse input string in a different language you should import monday and replace

```go
  time.ParseInLocation(yourLayout, yourString, yourLocation)
```

with

```go
  // Change LocaleEnUS to the locale you want to use for translation
  monday.ParseInLocation(yourLayout, yourString, yourLocation, monday.LocaleEnUS)
```

Predefined formats
---------------------

Monday declares some predefined formats: Full, Long, Medium, Short, DateTime formats for each locale. E.g. to get
short format for any locale you can use map:

```go
monday.ShortFormatsByLocale[locale]
```

Usage notes
-----------

**Monday** is not an alternative to standard **time** package. It is a temporary solution to use while
the internationalization features are not ready.

That's why **monday** doesn't create any additional parsing algorithms, layout identifiers. It is just
a wrapper for time.Format and time.ParseInLocation and uses all the same layout IDs, constants, etc.

So, the changes you need to temporarily switch to **monday** (while the internationalization features are being developed)
are minimal: you preserve your layout, your time object, your parsed date string formats and the only change is
the func call itself.

Locales
----

Supported locales are listed in **locale.go** file.

```
const (
    LocaleEnUS = "en_US" // English (United States)
    LocaleEnGB = "en_GB" // English (United Kingdom)
    LocaleDaDK = "da_DK" // Danish (Denmark)
    LocaleNlBE = "nl_BE" // Dutch (Belgium)
    LocaleNlNL = "nl_NL" // Dutch (Netherlands)
    LocaleFiFI = "fi_FI" // Finnish (Finland)
    LocaleFrFR = "fr_FR" // French (France)
    LocaleFrCA = "fr_CA" // French (Canada)
    LocaleDeDE = "de_DE" // German (Germany)
    LocaleHuHU = "hu_HU" // Hungarian (Hungary)
    LocaleItIT = "it_IT" // Italian (Italy)
    LocaleNnNO = "nn_NO" // Norwegian Nynorsk (Norway)
    LocaleNbNO = "nb_NO" // Norwegian Bokmål (Norway)
    LocalePlPL = "pl_PL" // Polish (Poland)
    LocalePtPT = "pt_PT" // Portuguese (Portugal)
    LocalePtBR = "pt_BR" // Portuguese (Brazil)
    LocaleRoRO = "ro_RO" // Romanian (Romania)
    LocaleRuRU = "ru_RU" // Russian (Russia)
    LocaleEsES = "es_ES" // Spanish (Spain)
    LocaleCaES = "ca_ES" // Catalan (Spain)
    LocaleSvSE = "sv_SE" // Swedish (Sweden)
    LocaleTrTR = "tr_TR" // Turkish (Turkey)
    LocaleUkUA = "uk_UA" // Ukrainian (Ukraine)
    LocaleBgBG = "bg_BG" // Bulgarian (Bulgaria)
    LocaleZhCN = "zh_CN" // Chinese (Mainland)
    LocaleZhTW = "zh_TW" // Chinese (Taiwan)
    LocaleZhHK = "zh_HK" // Chinese (Hong Kong)
    LocaleKoKR = "ko_KR" // Korean (Korea)
    LocaleJaJP = "ja_JP" // Japanese (Japan)
    LocaleElGR = "el_GR" // Greek (Greece)
    LocaleIdID = "id_ID" // Indonesian (Indonesia)
    LocaleFrGP = "fr_GP" // French (Guadeloupe)
    LocaleFrLU = "fr_LU" // French (Luxembourg)
    LocaleFrMQ = "fr_MQ" // French (Martinique)
    LocaleFrGF = "fr_GF" // French (French Guiana)
    LocaleFrGF = "fr_RE" // French (Reunion)
    LocaleCsCZ = "cs_CZ" // Czech (Czech Republic)
    LocaleSlSI = "sl_SI" // Slovenian (Slovenia)
)
```

LocaleDetector
====

```go
    var timeLocaleDetector *monday.LocaleDetector = monday.NewLocaleDetector()
    dateTime, err := timeLocaleDetector.Parse(layout,datestr)
```
parses datetime with **unknown** locale (for now - layout must be defined, as for time.Parse())

useful for text parsing tools/crawlers (f.e.: rss-feeds crawler)

TODO:
 * make LocaleDetector insensitive to whitespaces count

Thread-safety
====

**Monday** initializes all its data once in the **init** func and then uses only
func calls and local vars. Thus, it's thread-safe and doesn't need any mutexes to be
used with.

Monday Licence
==========

The **Monday** library is released under the [BSD Licence](http://opensource.org/licenses/bsd-license.php)

[LICENCE file](https://github.com/goodsign/monday/blob/master/LICENCE)

Thanks
==========

* [Martin Angers](https://github.com/PuerkitoBio)
* Andrey Mirtchovski
* [mikespook](https://github.com/mikespook)
* [Luis Azevedo](https://github.com/braceta)
* [imikod](https://github.com/imikod)
* [Renato Serra](https://github.com/RenatoSerra22)
* [Zachary Stewart](https://github.com/ztstewart)
