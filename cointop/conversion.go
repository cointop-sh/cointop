package cointop

import (
	"fmt"
	"sort"
	"strings"

	color "github.com/miguelmota/cointop/cointop/common/color"
	"github.com/miguelmota/cointop/cointop/common/pad"
)

// keep these in alphabetical order
var fiatCurrencyNames = map[string]string{
	"AUD": "Australian Dollar",
	"BGN": "Bulgarian lev",
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
	"HRK": "Croatian kuna",
	"HUF": "Hungarian Forint",
	"IDR": "Indonesian Rupiah",
	"ILS": "Israeli New Shekel",
	"INR": "Indian Rupee",
	"ISK": "Icelandic króna",
	"JPY": "Japanese Yen",
	"KRW": "South Korean Won",
	"MXN": "Mexican Peso",
	"MYR": "Malaysian Ringgit",
	"NOK": "Norwegian Krone",
	"NZD": "New Zealand Dollar",
	"PHP": "Philippine Peso",
	"PKR": "Pakistani Rupe",
	"PLN": "Polish złoty",
	"RON": "Romanian leu",
	"RUB": "Russian Ruble",
	"SEK": "Swedish Krona",
	"SGD": "Singapore Dollar",
	"THB": "Thai Baht",
	"TRY": "Turkish lira",
	"TWD": "New Taiwan Dollar",
	"USD": "US Dollar",
	"VND": "Vietnamese Dong",
	"ZAR": "South African Rand",
}

var cryptocurrencyNames = map[string]string{
	"BTC": "Bitcoin",
	"ETH": "Ethereum",
}

// keep these in alphabetical order
var currencySymbolMap = map[string]string{
	"AUD": "$",
	"BGN": "Лв.",
	"BRL": "R$",
	"BTC": "Ƀ",
	"CAD": "$",
	"CFH": "₣",
	"CLP": "$",
	"CNY": "¥",
	"CZK": "Kč",
	"DKK": "Kr",
	"ETH": "Ξ",
	"EUR": "€",
	"GBP": "£",
	"HKD": "$",
	"HRK": "kn",
	"HUF": "Ft",
	"IDR": "Rp.",
	"ILS": "₪",
	"INR": "₹",
	"ISK": "kr",
	"JPY": "¥",
	"KRW": "₩",
	"MXN": "$",
	"MYR": "RM",
	"NOK": "kr",
	"NZD": "$",
	"PHP": "₱",
	"PKR": "₨",
	"PLN": "zł",
	"RON": "lei",
	"RUB": "Ꝑ",
	"SEK": "kr",
	"SGD": "S$",
	"THB": "฿",
	"TRY": "₺",
	"TWD": "NT$",
	"USD": "$",
	"VND": "₫",
	"ZAR": "R",
}

var alphanumericcharacters = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

// ConvertMenuView is structure for convert menu view
type ConvertMenuView struct {
	*View
}

// NewConvertMenuView returns a new convert menu view
func NewConvertMenuView() *ConvertMenuView {
	return &ConvertMenuView{NewView("convertmenu")}
}

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
	ct.debuglog("updateConvertMenu()")
	header := ct.colorscheme.MenuHeader(fmt.Sprintf(" Currency Conversion %s\n\n", pad.Left("[q] close menu ", ct.maxTableWidth-20, " ")))
	helpline := " Press the corresponding key to select currency for conversion\n\n"
	cnt := 0
	h := ct.Views.ConvertMenu.Height()
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
	ct.Update(func() error {
		if ct.Views.ConvertMenu.Backing() == nil {
			return nil
		}

		ct.Views.ConvertMenu.Backing().Clear()
		ct.Views.ConvertMenu.Backing().Frame = true
		fmt.Fprintln(ct.Views.ConvertMenu.Backing(), content)
		return nil
	})
}

func (ct *Cointop) setCurrencyConverstionFn(convert string) func() error {
	ct.debuglog("setCurrencyConverstionFn()")
	return func() error {
		ct.hideConvertMenu()

		// NOTE: return if the currency selection wasn't changed
		if ct.State.currencyConversion == convert {
			return nil
		}

		ct.State.currencyConversion = convert

		if err := ct.Save(); err != nil {
			return err
		}

		go ct.refreshAll()
		return nil
	}
}

// currencySymbol returns the symbol for the currency
func (ct *Cointop) currencySymbol() string {
	ct.debuglog("currencySymbol()")
	symbol, ok := currencySymbolMap[strings.ToUpper(ct.State.currencyConversion)]
	if ok {
		return symbol
	}

	return "$"
}

func (ct *Cointop) showConvertMenu() error {
	ct.debuglog("showConvertMenu()")
	ct.State.convertMenuVisible = true
	ct.updateConvertMenu()
	ct.SetActiveView(ct.Views.ConvertMenu.Name())
	return nil
}

func (ct *Cointop) hideConvertMenu() error {
	ct.debuglog("hideConvertMenu()")
	ct.State.convertMenuVisible = false
	ct.SetViewOnBottom(ct.Views.ConvertMenu.Name())
	ct.SetActiveView(ct.Views.Table.Name())
	ct.Update(func() error {
		if ct.Views.ConvertMenu.Backing() == nil {
			return nil
		}

		ct.Views.ConvertMenu.Backing().Clear()
		ct.Views.ConvertMenu.Backing().Frame = false
		fmt.Fprintln(ct.Views.ConvertMenu.Backing(), "")
		return nil
	})
	return nil
}

func (ct *Cointop) toggleConvertMenu() error {
	ct.debuglog("toggleConvertMenu()")
	ct.State.convertMenuVisible = !ct.State.convertMenuVisible
	if ct.State.convertMenuVisible {
		return ct.showConvertMenu()
	}
	return ct.hideConvertMenu()
}

// currencySymbol returns the symbol for the currency
func currencySymbol(currency string) string {
	symbol, ok := currencySymbolMap[strings.ToUpper(currency)]
	if ok {
		return symbol
	}

	return "$"
}
