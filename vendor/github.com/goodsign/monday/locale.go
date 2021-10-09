package monday

// Locale identifies locales supported by 'monday' package.
// Monday uses ICU locale identifiers. See http://userguide.icu-project.org/locale
type Locale string

// Locale constants represent all locales that are currently supported by
// this package.
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
	LocaleFrRE = "fr_RE" // French (Reunion)
	LocaleFrGF = "fr_GF" // French (French Guiana)
	LocaleCsCZ = "cs_CZ" // Czech (Czech Republic)
	LocaleSlSI = "sl_SI" // Slovenian (Slovenia)
)

// ListLocales returns all locales supported by the package.
func ListLocales() []Locale {
	return []Locale{
		LocaleEnUS,
		LocaleEnGB,
		LocaleDaDK,
		LocaleNlBE,
		LocaleNlNL,
		LocaleFiFI,
		LocaleFrFR,
		LocaleFrCA,
		LocaleDeDE,
		LocaleHuHU,
		LocaleItIT,
		LocaleNnNO,
		LocaleNbNO,
		LocalePlPL,
		LocalePtPT,
		LocalePtBR,
		LocaleRoRO,
		LocaleRuRU,
		LocaleEsES,
		LocaleCaES,
		LocaleSvSE,
		LocaleTrTR,
		LocaleUkUA,
		LocaleBgBG,
		LocaleZhCN,
		LocaleZhTW,
		LocaleZhHK,
		LocaleKoKR,
		LocaleJaJP,
		LocaleElGR,
		LocaleFrGP,
		LocaleFrLU,
		LocaleFrMQ,
		LocaleFrRE,
		LocaleFrGF,
		LocaleCsCZ,
		LocaleSlSI,
	}
}
