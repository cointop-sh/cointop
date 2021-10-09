/*
Package monday is a minimalistic translator for month and day of week names in time.Date objects

Introduction

Monday is not an alternative to standard time package. It is a temporary solution to use while
the internationalization features are not ready.

That's why monday doesn't create any additional parsing algorithms, layout identifiers. It is just
a wrapper for time.Format and time.ParseInLocation and uses all the same layout IDs, constants, etc.

Usage

Format usage:

    t := time.Date(2013, 4, 12, 0, 0, 0, 0, time.UTC)
    layout := "2 January 2006 15:04:05 MST"

    translationEnUS := monday.Format(t, layout, monday.LocaleEnUS)  // Instead of t.Format(layout)
    translationRuRU := monday.Format(t, layout, monday.LocaleRuRU)  // Instead of t.Format(layout)
    ...

Parse usage:
    layout := "2 January 2006 15:04:05 MST"

    // Instead of time.ParseInLocation(layout, "12 April 2013 00:00:00 MST", time.UTC)
    parsed := monday.ParseInLocation(layout, "12 April 2013 00:00:00 MST", time.UTC, monday.LocaleEnUS))
    parsed2 = monday.ParseInLocation(layout, "12 апреля 2013 00:00:00 MST", time.UTC, monday.LocaleRuRU))
    ...

Thread safety

Monday initializes all its data once in the init func and then uses only
func calls and local vars. Thus, it's thread-safe and doesn't need any mutexes to be
used with.

*/
package monday
