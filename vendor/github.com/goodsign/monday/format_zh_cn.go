package monday

import "strings"

// ============================================================
// Format rules for "zh_CN" locale: Chinese (Mainland)
// ============================================================

var longDayNamesZhCN = map[string]string{
	"Sunday":    "星期日",
	"Monday":    "星期一",
	"Tuesday":   "星期二",
	"Wednesday": "星期三",
	"Thursday":  "星期四",
	"Friday":    "星期五",
	"Saturday":  "星期六",
}

var shortDayNamesZhCN = map[string]string{
	"Sun": "日",
	"Mon": "一",
	"Tue": "二",
	"Wed": "三",
	"Thu": "四",
	"Fri": "五",
	"Sat": "六",
}

var longMonthNamesZhCN = map[string]string{
	"January":   "1 月",
	"February":  "2 月",
	"March":     "3 月",
	"April":     "4 月",
	"May":       "5 月",
	"June":      "6 月",
	"July":      "7 月",
	"August":    "8 月",
	"September": "9 月",
	"October":   "10 月",
	"November":  "11 月",
	"December":  "12 月",
}

var shortMonthNamesZhCN = map[string]string{
	"Jan": "1",
	"Feb": "2",
	"Mar": "3",
	"Apr": "4",
	"May": "5",
	"Jun": "6",
	"Jul": "7",
	"Aug": "8",
	"Sep": "9",
	"Oct": "10",
	"Nov": "11",
	"Dec": "12",
}

func parseFuncZhCommon(locale Locale) internalParseFunc {
	return func(layout, value string) string {
		// This special case is needed because Zh_CN/Zh/HK/... contains month names
		// that consist of a number, a delimiter, and '月'. Example: "October" = "10 月"
		//
		// This means that probably default time package layout IDs like 'January' or 'Jan'
		// shouldn't be used in Zh_*. But this is a time-compatible package, so someone
		// might actually use those and we need to replace those before doing standard procedures.
		for k, v := range knownMonthsLongReverse[locale] {
			value = strings.Replace(value, k, v, -1)
		}

		return commonFormatFunc(value, layout,
			knownDaysShortReverse[locale], knownDaysLongReverse[locale],
			knownMonthsShortReverse[locale], knownMonthsLongReverse[locale], knownPeriods[locale])
	}
}
