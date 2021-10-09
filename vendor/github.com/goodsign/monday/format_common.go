package monday

import "strings"

func findInString(where string, what string, foundIndex *int, trimRight *int) (found bool) {
	ind := strings.Index(strings.ToLower(where), strings.ToLower(what))
	if ind != -1 {
		*foundIndex = ind
		*trimRight = len(where) - ind - len(what)
		return true
	}

	return false
}

// commonFormatFunc is used for languages which don't have changed forms of month names dependent
// on their position (after day or standalone)
func commonFormatFunc(value, format string,
	knownDaysShort, knownDaysLong, knownMonthsShort, knownMonthsLong, knownPeriods map[string]string) string {
	l := stringToLayoutItems(value)
	f := stringToLayoutItems(format)
	if len(l) != len(f) {
		return value // layouts does not matches
	}

	sb := &strings.Builder{}
	sb.Grow(32) // Reasonable default size that should fit most strings.

	for i, v := range l {

		var knw map[string]string

		// number of symbols before replaced term
		foundIndex := 0
		trimRight := 0
		lowerCase := false
		switch {
		case findInString(f[i].item, "Monday", &foundIndex, &trimRight):
			knw = knownDaysLong
		case findInString(f[i].item, "Mon", &foundIndex, &trimRight):
			knw = knownDaysShort
		case findInString(f[i].item, "January", &foundIndex, &trimRight):
			knw = knownMonthsLong
		case findInString(f[i].item, "Jan", &foundIndex, &trimRight):
			knw = knownMonthsShort
		case findInString(f[i].item, "PM", &foundIndex, &trimRight):
			knw = knownPeriods
		case findInString(f[i].item, "pm", &foundIndex, &trimRight):
			lowerCase = true
			knw = knownPeriods
		}

		knw = mapToLowerCase(knw)

		if knw != nil {
			trimmedItem := strings.ToLower(v.item[foundIndex : len(v.item)-trimRight])

			tr, ok := knw[trimmedItem]
			if lowerCase == true {
				tr = strings.ToLower(tr)
			}

			if ok {
				sb.WriteString(v.item[:foundIndex])
				sb.WriteString(tr)
				sb.WriteString(v.item[len(v.item)-trimRight:])
			} else {
				sb.WriteString(v.item)
			}
		} else {
			sb.WriteString(v.item)
		}
	}
	return sb.String()
}

func hasDigitBefore(l []dateStringLayoutItem, position int) bool {
	if position >= 2 {
		return l[position-2].isDigit && len(l[position-2].item) <= 2
	}
	return false
}

// commonGenitiveFormatFunc is used for languages with genitive forms of names, like Russian.
func commonGenitiveFormatFunc(value, format string,
	knownDaysShort, knownDaysLong, knownMonthsShort, knownMonthsLong,
	knownMonthsGenShort, knownMonthsGenLong, knownPeriods map[string]string) string {

	l := stringToLayoutItems(value)
	f := stringToLayoutItems(format)

	if len(l) != len(f) {
		return value // layouts does not matches
	}

	sb := &strings.Builder{}
	sb.Grow(32) // Reasonable default size that should fit most strings.

	for i, v := range l {
		lowerCase := false
		var knw map[string]string
		switch f[i].item {
		case "Mon":
			knw = knownDaysShort
		case "Monday":
			knw = knownDaysLong
		case "Jan":
			if hasDigitBefore(l, i) {
				knw = knownMonthsGenShort
			} else {
				knw = knownMonthsShort
			}
		case "January":
			if hasDigitBefore(l, i) {
				knw = knownMonthsGenLong
			} else {
				knw = knownMonthsLong
			}
		case "PM":
			knw = knownPeriods
		case "pm":
			lowerCase = true
			knw = knownPeriods
		}

		knw = mapToLowerCase(knw)

		if knw != nil {
			tr, ok := knw[strings.ToLower(v.item)]
			if !ok {
				sb.WriteString(v.item)
				continue
			}
			if lowerCase == true {
				tr = strings.ToLower(tr)
			}
			sb.WriteString(tr)
		} else {
			sb.WriteString(v.item)
		}
	}
	return sb.String()
}

func createCommonFormatFunc(locale Locale) internalFormatFunc {
	return func(value, layout string) (res string) {
		return commonFormatFunc(value, layout,
			knownDaysShort[locale], knownDaysLong[locale], knownMonthsShort[locale], knownMonthsLong[locale], knownPeriods[locale])
	}
}

func createCommonFormatFuncWithGenitive(locale Locale) internalFormatFunc {
	return func(value, layout string) (res string) {
		return commonGenitiveFormatFunc(value, layout,
			knownDaysShort[locale], knownDaysLong[locale], knownMonthsShort[locale], knownMonthsLong[locale],
			knownMonthsGenitiveShort[locale], knownMonthsGenitiveLong[locale], knownPeriods[locale])
	}
}

func createCommonParseFunc(locale Locale) internalParseFunc {
	return func(layout, value string) string {
		return commonFormatFunc(value, layout,
			knownDaysShortReverse[locale], knownDaysLongReverse[locale],
			knownMonthsShortReverse[locale], knownMonthsLongReverse[locale], knownPeriodsReverse[locale])
	}
}

func createCommonParsetFuncWithGenitive(locale Locale) internalParseFunc {
	return func(layout, value string) (res string) {
		return commonGenitiveFormatFunc(value, layout,
			knownDaysShortReverse[locale], knownDaysLongReverse[locale],
			knownMonthsShortReverse[locale], knownMonthsLongReverse[locale],
			knownMonthsGenitiveShortReverse[locale], knownMonthsGenitiveLongReverse[locale], knownPeriodsReverse[locale])
	}
}

func mapToLowerCase(source map[string]string) map[string]string {
	result := make(map[string]string, len(source))
	for k, v := range source {
		result[strings.ToLower(k)] = v
	}
	return result
}
