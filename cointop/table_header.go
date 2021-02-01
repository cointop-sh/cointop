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
	offset := 0
	lb := "["
	rb := "]"
	noSort := ct.IsPriceAlertsVisible()
	if noSort {
		offset = 2
		lb = ""
		rb = ""
	}
	possibleHeaders := map[string]*t{
		"rank":             {baseColor, fmt.Sprintf("%sr%sank", lb, rb), 0, 1 + offset, " "},
		"name":             {baseColor, fmt.Sprintf("%sn%same", lb, rb), 0, 11 + offset, " "},
		"symbol":           {baseColor, fmt.Sprintf("%ss%symbol", lb, rb), 4, 0 + offset, " "},
		"target_price":     {baseColor, fmt.Sprintf("%st%sarget price", lb, rb), 2, 0 + offset, " "},
		"price":            {baseColor, fmt.Sprintf("%sp%srice", lb, rb), 2, 0 + offset, " "},
		"frequency":        {baseColor, "frequency", 1, 0, " "},
		"holdings":         {baseColor, fmt.Sprintf("%sh%soldings", lb, rb), 5, 0 + offset, " "},
		"balance":          {baseColor, fmt.Sprintf("%sb%salance", lb, rb), 5, 0, " "},
		"marketcap":        {baseColor, fmt.Sprintf("%sm%sarket cap", lb, rb), 5, 0 + offset, " "},
		"24h_volume":       {baseColor, fmt.Sprintf("24H %sv%solume", lb, rb), 3, 0 + offset, " "},
		"1h_change":        {baseColor, fmt.Sprintf("%s1%sH%%", lb, rb), 5, 0 + offset, " "},
		"24h_change":       {baseColor, fmt.Sprintf("%s2%s4H%%", lb, rb), 3, 0 + offset, " "},
		"7d_change":        {baseColor, fmt.Sprintf("%s7%sD%%", lb, rb), 4, 0 + offset, " "},
		"total_supply":     {baseColor, fmt.Sprintf("%st%sotal supply", lb, rb), 7, 0 + offset, " "},
		"available_supply": {baseColor, fmt.Sprintf("%sa%svailable supply", lb, rb), 1, 0 + offset, " "},
		"percent_holdings": {baseColor, fmt.Sprintf("%s%%%sholdings", lb, rb), 2, 0 + offset, " "},
		"last_updated":     {baseColor, fmt.Sprintf("last %su%spdated", lb, rb), 3, 0, " "},
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
