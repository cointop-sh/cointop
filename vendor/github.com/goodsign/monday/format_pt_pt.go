package monday

import "strings"

// ============================================================
// Format rules for "pt_PT" locale: Portuguese (Portugal)
// ============================================================

var longDayNamesPtPT = map[string]string{
	"Sunday":    "Domingo",
	"Monday":    "Segunda-feira",
	"Tuesday":   "Terça-feira",
	"Wednesday": "Quarta-feira",
	"Thursday":  "Quinta-feira",
	"Friday":    "Sexta-feira",
	"Saturday":  "Sábado",
}

var shortDayNamesPtPT = map[string]string{
	"Sun": "dom",
	"Mon": "seg",
	"Tue": "ter",
	"Wed": "qua",
	"Thu": "qui",
	"Fri": "sex",
	"Sat": "sáb",
}

var longMonthNamesPtPT = map[string]string{
	"January":   "Janeiro",
	"February":  "Fevereiro",
	"March":     "Março",
	"April":     "Abril",
	"May":       "Maio",
	"June":      "Junho",
	"July":      "Julho",
	"August":    "Agosto",
	"September": "Setembro",
	"October":   "Outubro",
	"November":  "Novembro",
	"December":  "Dezembro",
}

var shortMonthNamesPtPT = map[string]string{
	"Jan": "Jan",
	"Feb": "Fev",
	"Mar": "Mar",
	"Apr": "Abr",
	"May": "Mai",
	"Jun": "Jun",
	"Jul": "Jul",
	"Aug": "Ago",
	"Sep": "Set",
	"Oct": "Out",
	"Nov": "Nov",
	"Dec": "Dez",
}

func parseFuncPtCommon(locale Locale) internalParseFunc {
	return func(layout, value string) string {
		// This special case is needed because Pt_PT/Pt_BR/... contains day-of-week names
		// that consist of two words and a delimiter (like 'terça-feira'). These
		// should be replaced before using the standard procedure correctly.
		for k, v := range knownDaysLongReverse[locale] {
			value = strings.Replace(value, k, v, -1)
		}

		return commonFormatFunc(value, layout,
			knownDaysShortReverse[locale], knownDaysLongReverse[locale],
			knownMonthsShortReverse[locale], knownMonthsLongReverse[locale], knownPeriods[locale])
	}
}
