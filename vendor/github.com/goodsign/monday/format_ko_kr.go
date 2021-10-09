package monday

import "strings"

// ============================================================
// Format rules for "ko_KR" locale: Korean (Korea)
// ============================================================

var longDayNamesKoKR = map[string]string{
	"Sunday":    "일요일",
	"Monday":    "월요일",
	"Tuesday":   "화요일",
	"Wednesday": "수요일",
	"Thursday":  "목요일",
	"Friday":    "금요일",
	"Saturday":  "토요일",
}

var shortDayNamesKoKR = map[string]string{
	"Sun": "일",
	"Mon": "월",
	"Tue": "화",
	"Wed": "수",
	"Thu": "목",
	"Fri": "금",
	"Sat": "토",
}

var longMonthNamesKoKR = map[string]string{
	"January":   "1월",
	"February":  "2월",
	"March":     "3월",
	"April":     "4월",
	"May":       "5월",
	"June":      "6월",
	"July":      "7월",
	"August":    "8월",
	"September": "9월",
	"October":   "10월",
	"November":  "11월",
	"December":  "12월",
}

var shortMonthNamesKoKR = map[string]string{
	"Jan": "1월",
	"Feb": "2월",
	"Mar": "3월",
	"Apr": "4월",
	"May": "5월",
	"Jun": "6월",
	"Jul": "7월",
	"Aug": "8월",
	"Sep": "9월",
	"Oct": "10월",
	"Nov": "11월",
	"Dec": "12월",
}

var periodsKoKR = map[string]string{
	"am": "오전",
	"pm": "오후",
	"AM": "오전",
	"PM": "오후",
}

func parseFuncKoCommon(locale Locale) internalParseFunc {
	return func(layout, value string) string {
		// This special case is needed because ko_KR... contains month names
		// that consist of a number, a delimiter, and '월'. Example: "September" = "9 월"
		//
		// This means that probably default time package layout IDs like 'January' or 'Jan'
		// shouldn't be used in ko_KR. But this is a time-compatible package, so someone
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
