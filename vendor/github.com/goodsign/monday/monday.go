package monday

import (
	"fmt"
	"time"
)

// internalFormatFunc is a preprocessor for default time.Format func
type internalFormatFunc func(value, layout string) string

var internalFormatFuncs = map[Locale]internalFormatFunc{
	LocaleEnUS: createCommonFormatFunc(LocaleEnUS),
	LocaleEnGB: createCommonFormatFunc(LocaleEnGB),
	LocaleDaDK: createCommonFormatFunc(LocaleDaDK),
	LocaleNlBE: createCommonFormatFunc(LocaleNlBE),
	LocaleNlNL: createCommonFormatFunc(LocaleNlNL),
	LocaleFrFR: createCommonFormatFunc(LocaleFrFR),
	LocaleFrCA: createCommonFormatFunc(LocaleFrFR),
	LocaleFrGP: createCommonFormatFunc(LocaleFrFR),
	LocaleFrLU: createCommonFormatFunc(LocaleFrFR),
	LocaleFrMQ: createCommonFormatFunc(LocaleFrFR),
	LocaleFrGF: createCommonFormatFunc(LocaleFrFR),
	LocaleFrRE: createCommonFormatFunc(LocaleFrFR),
	LocaleRuRU: createCommonFormatFuncWithGenitive(LocaleRuRU),
	LocaleFiFI: createCommonFormatFuncWithGenitive(LocaleFiFI),
	LocaleDeDE: createCommonFormatFunc(LocaleDeDE),
	LocaleHuHU: createCommonFormatFunc(LocaleHuHU),
	LocaleItIT: createCommonFormatFunc(LocaleItIT),
	LocaleNnNO: createCommonFormatFunc(LocaleNnNO),
	LocaleNbNO: createCommonFormatFunc(LocaleNbNO),
	LocalePlPL: createCommonFormatFunc(LocalePlPL),
	LocalePtPT: createCommonFormatFunc(LocalePtPT),
	LocalePtBR: createCommonFormatFunc(LocalePtBR),
	LocaleRoRO: createCommonFormatFunc(LocaleRoRO),
	LocaleEsES: createCommonFormatFunc(LocaleEsES),
	LocaleCaES: createCommonFormatFunc(LocaleCaES),
	LocaleSvSE: createCommonFormatFunc(LocaleSvSE),
	LocaleTrTR: createCommonFormatFunc(LocaleTrTR),
	LocaleUkUA: createCommonFormatFuncWithGenitive(LocaleUkUA),
	LocaleBgBG: createCommonFormatFunc(LocaleBgBG),
	LocaleZhCN: createCommonFormatFunc(LocaleZhCN),
	LocaleZhTW: createCommonFormatFunc(LocaleZhTW),
	LocaleZhHK: createCommonFormatFunc(LocaleZhHK),
	LocaleKoKR: createCommonFormatFunc(LocaleKoKR),
	LocaleJaJP: createCommonFormatFunc(LocaleJaJP),
	LocaleElGR: createCommonFormatFuncWithGenitive(LocaleElGR),
	LocaleIdID: createCommonFormatFunc(LocaleIdID),
	LocaleCsCZ: createCommonFormatFunc(LocaleCsCZ),
	LocaleSlSI: createCommonFormatFunc(LocaleSlSI),
}

// internalParseFunc is a preprocessor for default time.ParseInLocation func
type internalParseFunc func(layout, value string) string

var internalParseFuncs = map[Locale]internalParseFunc{
	LocaleEnUS: createCommonParseFunc(LocaleEnUS),
	LocaleEnGB: createCommonParseFunc(LocaleEnGB),
	LocaleDaDK: createCommonParseFunc(LocaleDaDK),
	LocaleNlBE: createCommonParseFunc(LocaleNlBE),
	LocaleNlNL: createCommonParseFunc(LocaleNlNL),
	LocaleFrFR: createCommonParseFunc(LocaleFrFR),
	LocaleFrCA: createCommonParseFunc(LocaleFrFR),
	LocaleFrGP: createCommonParseFunc(LocaleFrFR),
	LocaleFrLU: createCommonParseFunc(LocaleFrFR),
	LocaleFrMQ: createCommonParseFunc(LocaleFrFR),
	LocaleFrGF: createCommonParseFunc(LocaleFrFR),
	LocaleFrRE: createCommonParseFunc(LocaleFrFR),
	LocaleRuRU: createCommonParsetFuncWithGenitive(LocaleRuRU),
	LocaleFiFI: createCommonParsetFuncWithGenitive(LocaleFiFI),
	LocaleDeDE: createCommonParseFunc(LocaleDeDE),
	LocaleHuHU: createCommonParseFunc(LocaleHuHU),
	LocaleItIT: createCommonParseFunc(LocaleItIT),
	LocaleNnNO: createCommonParseFunc(LocaleNnNO),
	LocaleNbNO: createCommonParseFunc(LocaleNbNO),
	LocalePlPL: parseFuncPtCommon(LocalePlPL),
	LocalePtPT: parseFuncPtCommon(LocalePtPT),
	LocalePtBR: parseFuncPtCommon(LocalePtBR),
	LocaleRoRO: createCommonParseFunc(LocaleRoRO),
	LocaleEsES: createCommonParseFunc(LocaleEsES),
	LocaleCaES: createCommonParseFunc(LocaleCaES),
	LocaleSvSE: createCommonParseFunc(LocaleSvSE),
	LocaleTrTR: createCommonParseFunc(LocaleTrTR),
	LocaleUkUA: createCommonParsetFuncWithGenitive(LocaleUkUA),
	LocaleBgBG: createCommonParseFunc(LocaleBgBG),
	LocaleZhCN: parseFuncZhCommon(LocaleZhCN),
	LocaleZhTW: parseFuncZhCommon(LocaleZhTW),
	LocaleZhHK: parseFuncZhCommon(LocaleZhHK),
	LocaleKoKR: parseFuncKoCommon(LocaleKoKR),
	LocaleJaJP: parseFuncJaCommon(LocaleJaJP),
	LocaleElGR: createCommonParsetFuncWithGenitive(LocaleElGR),
	LocaleIdID: createCommonParseFunc(LocaleIdID),
	LocaleCsCZ: createCommonParseFunc(LocaleCsCZ),
	LocaleSlSI: createCommonParseFunc(LocaleSlSI),
}

var knownDaysShort = map[Locale]map[string]string{}           // Mapping for 'Format', days of week, short form
var knownDaysLong = map[Locale]map[string]string{}            // Mapping for 'Format', days of week, long form
var knownMonthsLong = map[Locale]map[string]string{}          // Mapping for 'Format', months: long form
var knownMonthsShort = map[Locale]map[string]string{}         // Mapping for 'Format', months: short form
var knownMonthsGenitiveShort = map[Locale]map[string]string{} // Mapping for 'Format', special for names in genitive, short form
var knownMonthsGenitiveLong = map[Locale]map[string]string{}  // Mapping for 'Format', special for names in genitive, long form
var knownPeriods = map[Locale]map[string]string{}             // Mapping for 'Format', AM/PM

// Reverse maps for the same

var knownDaysShortReverse = map[Locale]map[string]string{}           // Mapping for 'Format', days of week, short form
var knownDaysLongReverse = map[Locale]map[string]string{}            // Mapping for 'Format', days of week, long form
var knownMonthsLongReverse = map[Locale]map[string]string{}          // Mapping for 'Format', months: long form
var knownMonthsShortReverse = map[Locale]map[string]string{}         // Mapping for 'Format', months: short form
var knownMonthsGenitiveShortReverse = map[Locale]map[string]string{} // Mapping for 'Format', special for names in genitive, short form
var knownMonthsGenitiveLongReverse = map[Locale]map[string]string{}  // Mapping for 'Format', special for names in genitive, long form
var knownPeriodsReverse = map[Locale]map[string]string{}             // Mapping for 'Format', AM/PM

func init() {
	fillKnownWords()
}

func fillKnownWords() {

	// En_US: English (United States)
	fillKnownDaysLong(longDayNamesEnUS, LocaleEnUS)
	fillKnownDaysShort(shortDayNamesEnUS, LocaleEnUS)
	fillKnownMonthsLong(longMonthNamesEnUS, LocaleEnUS)
	fillKnownMonthsShort(shortMonthNamesEnUS, LocaleEnUS)

	// En_GB: English (United Kingdom)
	fillKnownDaysLong(longDayNamesEnUS, LocaleEnGB)
	fillKnownDaysShort(shortDayNamesEnUS, LocaleEnGB)
	fillKnownMonthsLong(longMonthNamesEnUS, LocaleEnGB)
	fillKnownMonthsShort(shortMonthNamesEnUS, LocaleEnGB)

	// Da_DK: Danish (Denmark)
	fillKnownDaysLong(longDayNamesDaDK, LocaleDaDK)
	fillKnownDaysShort(shortDayNamesDaDK, LocaleDaDK)
	fillKnownMonthsLong(longMonthNamesDaDK, LocaleDaDK)
	fillKnownMonthsShort(shortMonthNamesDaDK, LocaleDaDK)

	// Nl_BE: Dutch (Belgium)
	fillKnownDaysLong(longDayNamesNlBE, LocaleNlBE)
	fillKnownDaysShort(shortDayNamesNlBE, LocaleNlBE)
	fillKnownMonthsLong(longMonthNamesNlBE, LocaleNlBE)
	fillKnownMonthsShort(shortMonthNamesNlBE, LocaleNlBE)

	// Nl_NL: Dutch (Netherlands)
	fillKnownDaysLong(longDayNamesNlBE, LocaleNlNL)
	fillKnownDaysShort(shortDayNamesNlBE, LocaleNlNL)
	fillKnownMonthsLong(longMonthNamesNlBE, LocaleNlNL)
	fillKnownMonthsShort(shortMonthNamesNlBE, LocaleNlNL)

	// Fi_FI: Finnish (Finland)
	fillKnownDaysLong(longDayNamesFiFI, LocaleFiFI)
	fillKnownDaysShort(shortDayNamesFiFI, LocaleFiFI)
	fillKnownMonthsLong(longMonthNamesFiFI, LocaleFiFI)
	fillKnownMonthsShort(shortMonthNamesFiFI, LocaleFiFI)
	fillKnownMonthsGenitiveLong(longMonthNamesGenitiveFiFI, LocaleFiFI)
	fillKnownMonthsGenitiveShort(shortMonthNamesFiFI, LocaleFiFI)

	// Fr_FR: French (France)
	fillKnownDaysLong(longDayNamesFrFR, LocaleFrFR)
	fillKnownDaysShort(shortDayNamesFrFR, LocaleFrFR)
	fillKnownMonthsLong(longMonthNamesFrFR, LocaleFrFR)
	fillKnownMonthsShort(shortMonthNamesFrFR, LocaleFrFR)

	// Fr_CA: French (Canada)
	fillKnownDaysLong(longDayNamesFrFR, LocaleFrCA)
	fillKnownDaysShort(shortDayNamesFrFR, LocaleFrCA)
	fillKnownMonthsLong(longMonthNamesFrFR, LocaleFrCA)
	fillKnownMonthsShort(shortMonthNamesFrFR, LocaleFrCA)

	// Fr_GP: French (Guadeloupe)
	fillKnownDaysLong(longDayNamesFrFR, LocaleFrGP)
	fillKnownDaysShort(shortDayNamesFrFR, LocaleFrGP)
	fillKnownMonthsLong(longMonthNamesFrFR, LocaleFrGP)
	fillKnownMonthsShort(shortMonthNamesFrFR, LocaleFrGP)

	// Fr_LU: French (Luxembourg)
	fillKnownDaysLong(longDayNamesFrFR, LocaleFrLU)
	fillKnownDaysShort(longDayNamesFrFR, LocaleFrLU)
	fillKnownMonthsLong(longDayNamesFrFR, LocaleFrLU)
	fillKnownMonthsShort(longDayNamesFrFR, LocaleFrLU)

	// Fr_MQ: French (Martinique)
	fillKnownDaysLong(longDayNamesFrFR, LocaleFrMQ)
	fillKnownDaysShort(longDayNamesFrFR, LocaleFrMQ)
	fillKnownMonthsLong(longDayNamesFrFR, LocaleFrMQ)
	fillKnownMonthsShort(longDayNamesFrFR, LocaleFrMQ)

	// Fr_GF: French (French Guiana)
	fillKnownDaysLong(longDayNamesFrFR, LocaleFrGF)
	fillKnownDaysShort(longDayNamesFrFR, LocaleFrGF)
	fillKnownMonthsLong(longDayNamesFrFR, LocaleFrGF)
	fillKnownMonthsShort(longDayNamesFrFR, LocaleFrGF)

	// Fr_RE: French (French Reunion)
	fillKnownDaysLong(longDayNamesFrFR, LocaleFrRE)
	fillKnownDaysShort(longDayNamesFrFR, LocaleFrRE)
	fillKnownMonthsLong(longDayNamesFrFR, LocaleFrRE)
	fillKnownMonthsShort(longDayNamesFrFR, LocaleFrRE)

	// De_DE: German (Germany)
	fillKnownDaysLong(longDayNamesDeDE, LocaleDeDE)
	fillKnownDaysShort(shortDayNamesDeDE, LocaleDeDE)
	fillKnownMonthsLong(longMonthNamesDeDE, LocaleDeDE)
	fillKnownMonthsShort(shortMonthNamesDeDE, LocaleDeDE)

	// Hu_HU: Hungarian (Hungary)
	fillKnownDaysLong(longDayNamesHuHU, LocaleHuHU)
	fillKnownDaysShort(shortDayNamesHuHU, LocaleHuHU)
	fillKnownMonthsLong(longMonthNamesHuHU, LocaleHuHU)
	fillKnownMonthsShort(shortMonthNamesHuHU, LocaleHuHU)

	// It_IT: Italian (Italy)
	fillKnownDaysLong(longDayNamesItIT, LocaleItIT)
	fillKnownDaysShort(shortDayNamesItIT, LocaleItIT)
	fillKnownMonthsLong(longMonthNamesItIT, LocaleItIT)
	fillKnownMonthsShort(shortMonthNamesItIT, LocaleItIT)

	// Nn_NO: Norwegian Nynorsk (Norway)
	fillKnownDaysLong(longDayNamesNnNO, LocaleNnNO)
	fillKnownDaysShort(shortDayNamesNnNO, LocaleNnNO)
	fillKnownMonthsLong(longMonthNamesNnNO, LocaleNnNO)
	fillKnownMonthsShort(shortMonthNamesNnNO, LocaleNnNO)

	// Nb_NO: Norwegian Bokmål (Norway)
	fillKnownDaysLong(longDayNamesNbNO, LocaleNbNO)
	fillKnownDaysShort(shortDayNamesNbNO, LocaleNbNO)
	fillKnownMonthsLong(longMonthNamesNbNO, LocaleNbNO)
	fillKnownMonthsShort(shortMonthNamesNbNO, LocaleNbNO)

	// Pl_PL: Polish (Poland)
	fillKnownDaysLong(longDayNamesPlPL, LocalePlPL)
	fillKnownDaysShort(shortDayNamesPlPL, LocalePlPL)
	fillKnownMonthsLong(longMonthNamesPlPL, LocalePlPL)
	fillKnownMonthsShort(shortMonthNamesPlPL, LocalePlPL)

	// Pt_PT: Portuguese (Portugal)
	fillKnownDaysLong(longDayNamesPtPT, LocalePtPT)
	fillKnownDaysShort(shortDayNamesPtPT, LocalePtPT)
	fillKnownMonthsLong(longMonthNamesPtPT, LocalePtPT)
	fillKnownMonthsShort(shortMonthNamesPtPT, LocalePtPT)

	// Pt_BR: Portuguese (Brazil)
	fillKnownDaysLong(longDayNamesPtBR, LocalePtBR)
	fillKnownDaysShort(shortDayNamesPtBR, LocalePtBR)
	fillKnownMonthsLong(longMonthNamesPtBR, LocalePtBR)
	fillKnownMonthsShort(shortMonthNamesPtBR, LocalePtBR)

	// Ro_RO: Romanian (Romania)
	fillKnownDaysLong(longDayNamesRoRO, LocaleRoRO)
	fillKnownDaysShort(shortDayNamesRoRO, LocaleRoRO)
	fillKnownMonthsLong(longMonthNamesRoRO, LocaleRoRO)
	fillKnownMonthsShort(shortMonthNamesRoRO, LocaleRoRO)

	// Ru_RU: Russian (Russia)
	fillKnownDaysLong(longDayNamesRuRU, LocaleRuRU)
	fillKnownDaysShort(shortDayNamesRuRU, LocaleRuRU)
	fillKnownMonthsLong(longMonthNamesRuRU, LocaleRuRU)
	fillKnownMonthsShort(shortMonthNamesRuRU, LocaleRuRU)
	fillKnownMonthsGenitiveLong(longMonthNamesGenitiveRuRU, LocaleRuRU)
	fillKnownMonthsGenitiveShort(shortMonthNamesGenitiveRuRU, LocaleRuRU)

	// Es_ES: Spanish (Spain)
	fillKnownDaysLong(longDayNamesEsES, LocaleEsES)
	fillKnownDaysShort(shortDayNamesEsES, LocaleEsES)
	fillKnownMonthsLong(longMonthNamesEsES, LocaleEsES)
	fillKnownMonthsShort(shortMonthNamesEsES, LocaleEsES)

	// Ca_ES: Catalan (Spain)
	fillKnownDaysLong(longDayNamesCaES, LocaleCaES)
	fillKnownDaysShort(shortDayNamesCaES, LocaleCaES)
	fillKnownMonthsLong(longMonthNamesCaES, LocaleCaES)
	fillKnownMonthsShort(shortMonthNamesCaES, LocaleCaES)

	// Sv_SE: Swedish (Sweden)
	fillKnownDaysLong(longDayNamesSvSE, LocaleSvSE)
	fillKnownDaysShort(shortDayNamesSvSE, LocaleSvSE)
	fillKnownMonthsLong(longMonthNamesSvSE, LocaleSvSE)
	fillKnownMonthsShort(shortMonthNamesSvSE, LocaleSvSE)

	// Tr_TR: Turkish (Turkey)
	fillKnownDaysLong(longDayNamesTrTR, LocaleTrTR)
	fillKnownDaysShort(shortDayNamesTrTR, LocaleTrTR)
	fillKnownMonthsLong(longMonthNamesTrTR, LocaleTrTR)
	fillKnownMonthsShort(shortMonthNamesTrTR, LocaleTrTR)

	// Uk_UA: Ukrainian (Ukraine)
	fillKnownDaysLong(longDayNamesUkUA, LocaleUkUA)
	fillKnownDaysShort(shortDayNamesUkUA, LocaleUkUA)
	fillKnownMonthsLong(longMonthNamesUkUA, LocaleUkUA)
	fillKnownMonthsShort(shortMonthNamesUkUA, LocaleUkUA)
	fillKnownMonthsGenitiveLong(longMonthNamesGenitiveUkUA, LocaleUkUA)
	fillKnownMonthsGenitiveShort(shortMonthNamesGenitiveUkUA, LocaleUkUA)

	// Bg_BG: Bulgarian (Bulgaria)
	fillKnownDaysLong(longDayNamesBgBG, LocaleBgBG)
	fillKnownDaysShort(shortDayNamesBgBG, LocaleBgBG)
	fillKnownMonthsLong(longMonthNamesBgBG, LocaleBgBG)
	fillKnownMonthsShort(shortMonthNamesBgBG, LocaleBgBG)

	// Zh_CN: Chinese (Mainland)
	fillKnownDaysLong(longDayNamesZhCN, LocaleZhCN)
	fillKnownDaysShort(shortDayNamesZhCN, LocaleZhCN)
	fillKnownMonthsLong(longMonthNamesZhCN, LocaleZhCN)
	fillKnownMonthsShort(shortMonthNamesZhCN, LocaleZhCN)

	// Zh_TW: Chinese (Taiwan)
	fillKnownDaysLong(longDayNamesZhTW, LocaleZhTW)
	fillKnownDaysShort(shortDayNamesZhTW, LocaleZhTW)
	fillKnownMonthsLong(longMonthNamesZhTW, LocaleZhTW)
	fillKnownMonthsShort(shortMonthNamesZhTW, LocaleZhTW)

	// Zh_HK: Chinese (Hong Kong)
	fillKnownDaysLong(longDayNamesZhHK, LocaleZhHK)
	fillKnownDaysShort(shortDayNamesZhHK, LocaleZhHK)
	fillKnownMonthsLong(longMonthNamesZhHK, LocaleZhHK)
	fillKnownMonthsShort(shortMonthNamesZhHK, LocaleZhHK)

	// Ko_KR: Korean (Korea)
	fillKnownDaysLong(longDayNamesKoKR, LocaleKoKR)
	fillKnownDaysShort(shortDayNamesKoKR, LocaleKoKR)
	fillKnownMonthsLong(longMonthNamesKoKR, LocaleKoKR)
	fillKnownMonthsShort(shortMonthNamesKoKR, LocaleKoKR)
	fillKnownPeriods(periodsKoKR, LocaleKoKR)

	// Ja_JP: Japanese (Japan)
	fillKnownDaysLong(longDayNamesJaJP, LocaleJaJP)
	fillKnownDaysShort(shortDayNamesJaJP, LocaleJaJP)
	fillKnownMonthsLong(longMonthNamesJaJP, LocaleJaJP)
	fillKnownMonthsShort(shortMonthNamesJaJP, LocaleJaJP)
	fillKnownPeriods(periodsJaJP, LocaleJaJP)

	// El_GR: Greek (Greece)
	fillKnownDaysLong(longDayNamesElGR, LocaleElGR)
	fillKnownDaysShort(shortDayNamesElGR, LocaleElGR)
	fillKnownMonthsLong(longMonthNamesElGR, LocaleElGR)
	fillKnownMonthsShort(shortMonthNamesElGR, LocaleElGR)
	fillKnownMonthsGenitiveLong(longMonthNamesGenitiveElGR, LocaleElGR)
	fillKnownMonthsGenitiveShort(shortMonthNamesElGR, LocaleElGR)
	fillKnownPeriods(periodsElGR, LocaleElGR)

	// Id_ID: Indonesia (Indonesia)
	fillKnownDaysLong(longDayNamesIdID, LocaleIdID)
	fillKnownDaysShort(shortDayNamesIdID, LocaleIdID)
	fillKnownMonthsLong(longMonthNamesIdID, LocaleIdID)
	fillKnownMonthsShort(shortMonthNamesIdID, LocaleIdID)

	// Cs_CZ: Czech (Czech Republic)
	fillKnownDaysLong(longDayNamesCsCZ, LocaleCsCZ)
	fillKnownDaysShort(shortDayNamesCsCZ, LocaleCsCZ)
	fillKnownMonthsLong(longMonthNamesCsCZ, LocaleCsCZ)
	fillKnownMonthsShort(shortMonthNamesCsCZ, LocaleCsCZ)

	// Sl_SI: Slovenian (Slovenia)
	fillKnownDaysLong(longDayNamesSlSI, LocaleSlSI)
	fillKnownDaysShort(shortDayNamesSlSI, LocaleSlSI)
	fillKnownMonthsLong(longMonthNamesSlSI, LocaleSlSI)
	fillKnownMonthsShort(shortMonthNamesSlSI, LocaleSlSI)
}

func fill(src map[string]string, dest map[Locale]map[string]string, locale Locale) {
	loc, ok := dest[locale]

	if !ok {
		loc = make(map[string]string, len(src))
		dest[locale] = loc
	}

	for k, v := range src {
		loc[k] = v
	}
}

func fillReverse(src map[string]string, dest map[Locale]map[string]string, locale Locale) {
	loc, ok := dest[locale]

	if !ok {
		loc = make(map[string]string, len(src))
		dest[locale] = loc
	}

	for k, v := range src {
		loc[v] = k
	}
}

func fillKnownMonthsGenitiveShort(src map[string]string, locale Locale) {
	fillReverse(src, knownMonthsGenitiveShortReverse, locale)
	fill(src, knownMonthsGenitiveShort, locale)
}

func fillKnownMonthsGenitiveLong(src map[string]string, locale Locale) {
	fillReverse(src, knownMonthsGenitiveLongReverse, locale)
	fill(src, knownMonthsGenitiveLong, locale)
}

func fillKnownDaysShort(src map[string]string, locale Locale) {
	fillReverse(src, knownDaysShortReverse, locale)
	fill(src, knownDaysShort, locale)
}

func fillKnownDaysLong(src map[string]string, locale Locale) {
	fillReverse(src, knownDaysLongReverse, locale)
	fill(src, knownDaysLong, locale)
}

func fillKnownMonthsShort(src map[string]string, locale Locale) {
	fillReverse(src, knownMonthsShortReverse, locale)
	fill(src, knownMonthsShort, locale)
}

func fillKnownMonthsLong(src map[string]string, locale Locale) {
	fillReverse(src, knownMonthsLongReverse, locale)
	fill(src, knownMonthsLong, locale)
}

func fillKnownPeriods(src map[string]string, locale Locale) {
	fillReverse(src, knownPeriodsReverse, locale)
	fill(src, knownPeriods, locale)
}

// Format is the standard time.Format wrapper, that replaces known standard 'time' package
// identifiers for months and days to their equivalents in the specified language.
//
// Values of variables 'longDayNames', 'shortDayNames', 'longMonthNames', 'shortMonthNames'
// from file 'time/format.go' (revision 'go1') are chosen as the 'known' words.
//
// Some languages have specific behavior, e.g. in Russian language
// month names have different suffix when they are presented stand-alone (i.e. in a list or something)
// and yet another one when they are part of a formatted date.
// So, even though March is "Март" in Russian, correctly formatted today's date would be: "7 марта 2007".
// Thus, some transformations for some languages may be a bit more complex than just plain replacements.
func Format(dt time.Time, layout string, locale Locale) string {
	fm := dt.Format(layout)
	intFunc, ok := internalFormatFuncs[locale]
	if !ok {
		return fm
	}
	return intFunc(fm, layout)
}

// ParseInLocation is the standard time.ParseInLocation wrapper, which replaces
// known month/day translations for a specified locale back to English before
// calling time.ParseInLocation. So, you can parse localized dates with this wrapper.
func ParseInLocation(layout, value string, loc *time.Location, locale Locale) (time.Time, error) {
	intFunc, ok := internalParseFuncs[locale]
	if ok {
		value = intFunc(layout, value)
	} else {
		return time.Now(), fmt.Errorf("unsupported locale: %v", locale)
	}

	return time.ParseInLocation(layout, value, loc)
}

// Parse is the standard time.Parse wrapper, which replaces
// known month/day translations for a specified locale back to English before
// calling time.Parse.
func Parse(layout, value string, locale Locale) (time.Time, error) {
	intFunc, ok := internalParseFuncs[locale]
	if ok {
		value = intFunc(layout, value)
	} else {
		return time.Now(), fmt.Errorf("unsupported locale: %v", locale)
	}

	return time.Parse(layout, value)
}

// GetShortDays retrieves the list of days for the given locale.
// "Short" days are abbreviated versions of the full day names. In English,
// for example, this might return "Tues" for "Tuesday". For certain locales,
// the long and short form of the days of the week may be the same.
//
// If the locale cannot be found, the resulting slice will be nil.
func GetShortDays(locale Locale) []string {
	days, ok := knownDaysShort[locale]
	if !ok {
		return nil
	}

	var dayOrder []string

	// according to https://www.timeanddate.com/calendar/days/monday.html
	// only Canada, USA and Japan use Sunday as first day of the week
	switch locale {
	case LocaleEnUS, LocaleJaJP:
		dayOrder = dayShortOrderSundayFirst
	default:
		dayOrder = dayShortOrderMondayFirst
	}

	ret := make([]string, 0, len(days))
	for _, day := range dayOrder {
		ret = append(ret, days[day])
	}
	return ret
}

// GetShortMonths retrieves the list of months for the given locale.
// "Short" months are abbreviated versions of the full month names. In
// English, for example, this might return "Jan" for "January". For
// certain locales, the long and short form of the months may be the same.
//
// If the locale cannot be found, the resulting slice will be nil.
func GetShortMonths(locale Locale) []string {
	months, ok := knownMonthsShort[locale]
	if !ok {
		return nil
	}

	ret := make([]string, 0, len(months))
	for _, m := range monthShortOrder {
		ret = append(ret, months[m])
	}
	return ret
}

// GetLongDays retrieves the list of days for the given locale. It will return
// the full name of the days of the week.
//
// If the locale cannot be found, the resulting slice will be nil.
func GetLongDays(locale Locale) []string {
	days, ok := knownDaysLong[locale]
	if !ok {
		return nil
	}

	var dayOrder []string

	// according to https://www.timeanddate.com/calendar/days/monday.html
	// only Canada, USA and Japan use Sunday as first day of the week
	switch locale {
	case LocaleEnUS, LocaleJaJP:
		dayOrder = dayLongOrderSundayFirst
	default:
		dayOrder = dayLongOrderMondayFirst
	}

	ret := make([]string, 0, len(days))
	for _, day := range dayOrder {
		ret = append(ret, days[day])
	}
	return ret
}

// GetLongMonths retrieves the list of months for the given locale. In
// contrast to the "short" version of this function, this functions returns
// the full name of the month.
//
// If the locale cannot be found, the resulting slice will be nil.
func GetLongMonths(locale Locale) []string {
	months, ok := knownMonthsLong[locale]
	if !ok {
		return nil
	}

	ret := make([]string, 0, len(months))
	for _, m := range monthLongOrder {
		ret = append(ret, months[m])
	}

	return ret
}
