package cointop

import (
	"fmt"
	"sort"

	"github.com/miguelmota/cointop/pkg/color"
	"github.com/miguelmota/cointop/pkg/pad"
)

var supportedfiatconversions = map[string]string{
	"AUD": "Australian Dollar",
	"BRL": "Brazilian Real",
	"CAD": "Canadian Dollar",
	"CFH": "Swiss Franc",
	"CLP": "Chilean Peso",
	"CNY": "Chinese Yuan",
	"CZK": "Czech Koruna",
	"DKK": "Danish Krone",
	"EUR": "Euro",
	"GBP": "British Pound",
	"HKD": "Hong Kong Dollar",
	"HUF": "Hungarian Forint",
	"IDR": "Indonesian Rupiah",
	"ILS": "Israeli New Shekel",
	"INR": "Indian Rupee",
	"JPY": "Japanese Yen",
	"KRW": "South Korean Won",
	"MXN": "Mexican Peso",
	"MYR": "Malaysian Ringgit",
	"NOK": "Norwegian Krone",
	"NZD": "New Zealand Dollar",
	"PLN": "Polish złoty",
	"PHP": "Philippine Piso",
	"PKR": "Pakistani Rupe",
	"RUB": "Russian Ruble",
	"SEK": "Swedish Krona",
	"SGD": "Singapore Dollar",
	"THB": "Thai Baht",
	"TRY": "Turkish lira",
	"TWD": "New Taiwan Dollar",
	"USD": "US Dollar",
	"ZAR": "South African Rand",
}

var supportedcryptoconversion = map[string]string{
	"BTC": "Bitcoin",
	"ETH": "Ethereum",
}

var currencysymbols = map[string]string{
	"AUD": "$",
	"BRL": "R$",
	"BTC": "Ƀ",
	"CAD": "$",
	"CFH": "₣",
	"CLP": "$",
	"CNY": "¥",
	"CZK": "Kč",
	"DKK": "Kr",
	"EUR": "€",
	"ETH": "Ξ",
	"GBP": "£",
	"HKD": "$",
	"HUF": "Ft",
	"IDR": "Rp.",
	"ILS": "₪",
	"INR": "₹",
	"JPY": "¥",
	"KRW": "₩",
	"MXN": "$",
	"MYR": "RM",
	"NOK": "kr",
	"NZD": "$",
	"PLN": "zł",
	"PHP": "₱",
	"PKR": "₨",
	"RUB": "Ꝑ",
	"SEK": "kr",
	"SGD": "S$",
	"THB": "฿",
	"TRY": "₺",
	"TWD": "NT$",
	"USD": "$",
	"ZAR": "R",
}

var alphanumericcharacters = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

func (ct *Cointop) supportedCurrencyConversions() map[string]string {
	all := map[string]string{}
	for k, v := range supportedfiatconversions {
		all[k] = v
	}
	for k, v := range supportedcryptoconversion {
		all[k] = v
	}
	return all
}

func (ct *Cointop) supportedFiatCurrencyConversions() map[string]string {
	return supportedfiatconversions
}

func (ct *Cointop) supportedCryptoCurrencyConversions() map[string]string {
	return supportedfiatconversions
}

func (ct *Cointop) sortedSupportedCurrencyConversions() []string {
	currencies := ct.supportedCurrencyConversions()
	keys := make([]string, 0, len(currencies))
	for k := range currencies {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (ct *Cointop) toggleConvertMenu() error {
	ct.convertmenuvisible = !ct.convertmenuvisible
	if ct.convertmenuvisible {
		return ct.showConvertMenu()
	}
	return ct.hideConvertMenu()
}

func (ct *Cointop) updateConvertMenu() {
	header := color.GreenBg(fmt.Sprintf(" Currency Conversion %s\n\n", pad.Left("[q] close menu ", ct.maxtablewidth-20, " ")))
	helpline := " Press the corresponding key to select currency for conversion\n\n"
	cnt := 0
	h := ct.viewHeight(ct.convertmenuviewname)
	percol := h - 5
	cols := make([][]string, percol)
	for i := range cols {
		cols[i] = make([]string, 20)
	}

	keys := ct.sortedSupportedCurrencyConversions()
	currencies := ct.supportedCurrencyConversions()
	for i, key := range keys {
		currency := currencies[key]
		if cnt%percol == 0 {
			cnt = 0
		}
		shortcut := string(alphanumericcharacters[i])
		if key == ct.currencyconversion {
			shortcut = color.YellowBold("*")
			key = color.WhiteBold(key)
			currency = color.YellowBold(currency)
		} else {
			key = color.White(key)
			currency = color.Yellow(currency)
		}
		item := fmt.Sprintf(" [ %1s ] %4s %-34s", shortcut, key, currency)
		cols[cnt] = append(cols[cnt], item)
		cnt = cnt + 1
	}

	var body string
	for i := 0; i < percol; i++ {
		var row string
		for j := 0; j < len(cols[i]); j++ {
			item := cols[i][j]
			row = fmt.Sprintf("%s%s", row, item)
		}
		body = fmt.Sprintf("%s%s\n", body, row)
	}

	content := fmt.Sprintf("%s%s%s", header, helpline, body)
	ct.update(func() {
		ct.convertmenuview.Clear()
		ct.convertmenuview.Frame = true
		fmt.Fprintln(ct.convertmenuview, content)
	})
}

func (ct *Cointop) showConvertMenu() error {
	ct.convertmenuvisible = true
	ct.updateConvertMenu()
	ct.setActiveView(ct.convertmenuviewname)
	return nil
}

func (ct *Cointop) hideConvertMenu() error {
	ct.convertmenuvisible = false
	ct.setViewOnBottom(ct.convertmenuviewname)
	ct.setActiveView(ct.tableviewname)
	ct.update(func() {
		ct.convertmenuview.Clear()
		ct.convertmenuview.Frame = false
		fmt.Fprintln(ct.convertmenuview, "")
	})
	return nil
}

func (ct *Cointop) setCurrencyConverstion(convert string) func() error {
	return func() error {
		ct.currencyconversion = convert
		ct.hideConvertMenu()
		go ct.refreshAll()
		return nil
	}
}

func (ct *Cointop) currencySymbol() string {
	return currencysymbols[ct.currencyconversion]
}
