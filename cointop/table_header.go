package cointop

import (
	"fmt"
	"strings"

	"github.com/miguelmota/cointop/pkg/ui"
)

// TableHeaderView is structure for table header view
type TableHeaderView = ui.View

// NewTableHeaderView returns a new table header view
func NewTableHeaderView() *TableHeaderView {
	var view *TableHeaderView = ui.NewView("table_header")
	return view
}

// UpdateTableHeader renders the table header
func (ct *Cointop) UpdateTableHeader() error {
	ct.debuglog("UpdateTableHeader()")

	type t struct {
		colorfn     func(a ...interface{}) string
		displaytext string
		padleft     int
		padright    int
		arrow       string
	}

	baseColor := ct.colorscheme.TableHeaderSprintf()
	possibleHeaders := map[string]*t{
		"rank":             {baseColor, "[r]ank", 0, 1, " "},
		"name":             {baseColor, "[n]ame", 0, 11, " "},
		"symbol":           {baseColor, "[s]ymbol", 4, 0, " "},
		"target_price":     {baseColor, "[t]arget price", 2, 0, " "},
		"price":            {baseColor, "[p]rice", 2, 0, " "},
		"frequency":        {baseColor, "frequency", 1, 0, " "},
		"holdings":         {baseColor, "[h]oldings", 5, 0, " "},
		"balance":          {baseColor, "[b]alance", 5, 0, " "},
		"marketcap":        {baseColor, "[m]arket cap", 5, 0, " "},
		"24h_volume":       {baseColor, "24H [v]olume", 3, 0, " "},
		"1h_change":        {baseColor, "[1]H%", 5, 0, " "},
		"24h_change":       {baseColor, "[2]4H%", 3, 0, " "},
		"7d_change":        {baseColor, "[7]D%", 4, 0, " "},
		"total_supply":     {baseColor, "[t]otal supply", 7, 0, " "},
		"available_supply": {baseColor, "[a]vailable supply", 0, 0, " "},
		"percent_holdings": {baseColor, "[%]holdings", 2, 0, " "},
		"last_updated":     {baseColor, "last [u]pdated", 3, 0, " "},
	}

	for k := range possibleHeaders {
		possibleHeaders[k].arrow = " "
		if ct.State.sortBy == k {
			possibleHeaders[k].colorfn = ct.colorscheme.TableHeaderColumnActiveSprintf()
			if ct.State.sortDesc {
				possibleHeaders[k].arrow = "▼"
			} else {
				possibleHeaders[k].arrow = "▲"
			}
		}
	}

	var cols []string
	switch ct.State.selectedView {
	case PortfolioView:
		cols = ct.GetPortfolioTableHeaders()
	case PriceAlertsView:
		cols = ct.GetPriceAlertsTableHeaders()
	default:
		cols = ct.GetCoinsTableHeaders()
	}

	var headers []string
	for _, v := range cols {
		s, ok := possibleHeaders[v]
		if !ok {
			continue
		}
		var str string
		d := s.arrow + s.displaytext
		if v == "price" || v == "balance" {
			d = s.arrow + ct.CurrencySymbol() + s.displaytext
		}

		str = fmt.Sprintf(
			"%s%s%s",
			strings.Repeat(" ", s.padleft),
			s.colorfn(d),
			strings.Repeat(" ", s.padright),
		)
		headers = append(headers, str)
	}

	ct.UpdateUI(func() error {
		return ct.Views.TableHeader.Update(strings.Join(headers, ""))
	})

	return nil
}
