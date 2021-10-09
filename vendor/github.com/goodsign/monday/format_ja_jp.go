package monday

import (
	"strings"
)

// ============================================================
// Format rules for "ja_JP" locale: Japanese
// ============================================================

var longDayNamesJaJP = map[string]string{
	"Sunday":    "日曜日",
	"Monday":    "月曜日",
	"Tuesday":   "火曜日",
	"Wednesday": "水曜日",
	"Thursday":  "木曜日",
	"Friday":    "金曜日",
	"Saturday":  "土曜日",
}

var shortDayNamesJaJP = map[string]string{
	"Sun": "日",
	"Mon": "月",
	"Tue": "火",
	"Wed": "水",
	"Thu": "木",
	"Fri": "金",
	"Sat": "土",
}

var longMonthNamesJaJP = map[string]string{
	"January":   "1月",
	"February":  "2月",
	"March":     "3月",
	"April":     "4月",
	"May":       "5月",
	"June":      "6月",
	"July":      "7月",
	"August":    "8月",
	"September": "9月",
	"October":   "10月",
	"November":  "11月",
	"December":  "12月",
}

var shortMonthNamesJaJP = map[string]string{
	"Jan": "1月",
	"Feb": "2月",
	"Mar": "3月",
	"Apr": "4月",
	"May": "5月",
	"Jun": "6月",
	"Jul": "7月",
	"Aug": "8月",
	"Sep": "9月",
	"Oct": "10月",
	"Nov": "11月",
	"Dec": "12月",
}

var periodsJaJP = map[string]string{
	"am": "午前",
	"pm": "午後",
	"AM": "午前",
	"PM": "午後",
}

func parseFuncJaCommon(locale Locale) internalParseFunc {
	return func(layout, value string) string {
		// This special case is needed because ja_JP... contains month names
		// that consist of a number, a delimiter, and '月'. Example: "October" = "10 月"
		//
		// This means that probably default time package layout IDs like 'January' or 'Jan'
		// shouldn't be used in ja_JP. But this is a time-compatible package, so someone
		// might actually use those and we need to replace those before doing standard procedures.
		for k, v := range knownMonthsLongReverse[locale] {
			value = strings.Replace(value, k, v, -1)
		}

		value = commonFormatFunc(value, layout,
			knownDaysShortReverse[locale], knownDaysLongReverse[locale],
			knownMonthsShortReverse[locale], knownMonthsLongReverse[locale], knownPeriods[locale])

		// knownPeriodsReverse has hash collisions
		for k, v := range knownPeriodsReverse[locale] {
			targetValue := strings.ToLower(v)
			if strings.Index(layout, "PM") != -1 {
				targetValue = strings.ToUpper(v)
			}
			value = strings.Replace(value, k, targetValue, -1)
		}

		return value
	}
}
