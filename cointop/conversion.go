package cointop

import (
	"fmt"
	"sort"

	color "github.com/miguelmota/cointop/cointop/common/color"
	"github.com/miguelmota/cointop/cointop/common/pad"
)

var fiatCurrencyNames = map[string]string{
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

var cryptocurrencyNames = map[string]string{
	"BTC": "Bitcoin",
	"ETH": "Ethereum",
}

var currencySymbol = map[string]string{
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
	for _, symbol := range ct.api.SupportedCurrencies() {
		if v, ok := fiatCurrencyNames[symbol]; ok {
			all[symbol] = v
		}
		if v, ok := cryptocurrencyNames[symbol]; ok {
			all[symbol] = v
		}
	}

	return all
}

func (ct *Cointop) supportedFiatCurrencyConversions() map[string]string {
	return fiatCurrencyNames
}

func (ct *Cointop) supportedCryptoCurrencyConversions() map[string]string {
	return cryptocurrencyNames
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

func (ct *Cointop) updateConvertMenu() {
	header := ct.colorscheme.MenuHeader(fmt.Sprintf(" Currency Conversion %s\n\n", pad.Left("[q] close menu ", ct.maxTableWidth-20, " ")))
	helpline := " Press the corresponding key to select currency for conversion\n\n"
	cnt := 0
	h := ct.viewHeight(ct.Views.ConvertMenu.Name)
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
		if key == ct.State.currencyConversion {
			shortcut = ct.colorscheme.MenuLabelActive(color.Bold("*"))
			key = ct.colorscheme.Menu(color.Bold(key))
			currency = ct.colorscheme.MenuLabelActive(color.Bold(currency))
		} else {
			key = ct.colorscheme.Menu(key)
			currency = ct.colorscheme.MenuLabel(currency)
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
		if ct.Views.ConvertMenu.Backing == nil {
			return
		}

		ct.Views.ConvertMenu.Backing.Clear()
		ct.Views.ConvertMenu.Backing.Frame = true
		fmt.Fprintln(ct.Views.ConvertMenu.Backing, content)
	})
}

func (ct *Cointop) setCurrencyConverstion(convert string) func() error {
	return func() error {
		ct.State.currencyConversion = convert
		ct.hideConvertMenu()
		go ct.refreshAll()
		return nil
	}
}

func (ct *Cointop) currencySymbol() string {
	symbol, ok := currencySymbol[ct.State.currencyConversion]
	if ok {
		return symbol
	}

	return "$"
}

func (ct *Cointop) showConvertMenu() error {
	ct.State.convertMenuVisible = true
	ct.updateConvertMenu()
	ct.setActiveView(ct.Views.ConvertMenu.Name)
	return nil
}

func (ct *Cointop) hideConvertMenu() error {
	ct.State.convertMenuVisible = false
	ct.setViewOnBottom(ct.Views.ConvertMenu.Name)
	ct.setActiveView(ct.Views.Table.Name)
	ct.update(func() {
		if ct.Views.ConvertMenu.Backing == nil {
			return
		}

		ct.Views.ConvertMenu.Backing.Clear()
		ct.Views.ConvertMenu.Backing.Frame = false
		fmt.Fprintln(ct.Views.ConvertMenu.Backing, "")
	})
	return nil
}

func (ct *Cointop) toggleConvertMenu() error {
	ct.State.convertMenuVisible = !ct.State.convertMenuVisible
	if ct.State.convertMenuVisible {
		return ct.showConvertMenu()
	}
	return ct.hideConvertMenu()
}
