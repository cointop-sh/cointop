package cointop

import (
	"fmt"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/miguelmota/cointop/pkg/pad"
	"github.com/miguelmota/cointop/pkg/ui"
)

// ArrowUp is up arrow unicode character
var ArrowUp = "▲"

// ArrowDown is down arrow unicode character
var ArrowDown = "▼"

// HeaderColumn is header column struct
type HeaderColumn struct {
	Slug       string
	Label      string
	PlainLabel string
}

// HeaderColumns are the header column widths
var HeaderColumns = map[string]*HeaderColumn{
	"rank": &HeaderColumn{
		Slug:       "rank",
		Label:      "[r]ank",
		PlainLabel: "rank",
	},
	"name": &HeaderColumn{
		Slug:       "name",
		Label:      "[n]ame",
		PlainLabel: "name",
	},
	"symbol": &HeaderColumn{
		Slug:       "symbol",
		Label:      "[s]ymbol",
		PlainLabel: "symbol",
	},
	"target_price": &HeaderColumn{
		Slug:       "target_price",
		Label:      "[t]target price",
		PlainLabel: "target price",
	},
	"price": &HeaderColumn{
		Slug:       "price",
		Label:      "[p]rice",
		PlainLabel: "price",
	},
	"frequency": &HeaderColumn{
		Slug:       "frequency",
		Label:      "frequency",
		PlainLabel: "frequency",
	},
	"holdings": &HeaderColumn{
		Slug:       "holdings",
		Label:      "[h]oldings",
		PlainLabel: "holdings",
	},
	"balance": &HeaderColumn{
		Slug:       "balance",
		Label:      "[b]alance",
		PlainLabel: "balance",
	},
	"market_cap": &HeaderColumn{
		Slug:       "market_cap",
		Label:      "[m]arket cap",
		PlainLabel: "market cap",
	},
	"24h_volume": &HeaderColumn{
		Slug:       "24h_volume",
		Label:      "24H [v]olume",
		PlainLabel: "24H volume",
	},
	"1h_change": &HeaderColumn{
		Slug:       "1h_change",
		Label:      "[1]H%",
		PlainLabel: "1H%",
	},
	"24h_change": &HeaderColumn{
		Slug:       "24h_change",
		Label:      "[2]4H%",
		PlainLabel: "24H%",
	},
	"7d_change": &HeaderColumn{
		Slug:       "7d_change",
		Label:      "[7]D%",
		PlainLabel: "7D%",
	},
	"30d_change": &HeaderColumn{
		Slug:       "30d_change",
		Label:      "[3]0D%",
		PlainLabel: "30D%",
	},
	"total_supply": &HeaderColumn{
		Slug:       "total_supply",
		Label:      "[t]otal supply",
		PlainLabel: "total supply",
	},
	"available_supply": &HeaderColumn{
		Slug:       "available_supply",
		Label:      "[a]vailable supply",
		PlainLabel: "available supply",
	},
	"percent_holdings": &HeaderColumn{
		Slug:       "percent_holdings",
		Label:      "[%]holdings",
		PlainLabel: "%holdings",
	},
	"last_updated": &HeaderColumn{
		Slug:       "last_updated",
		Label:      "last [u]pdated",
		PlainLabel: "last updated",
	},
}

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

	baseColor := ct.colorscheme.TableHeaderSprintf()
	noSort := ct.IsPriceAlertsVisible()
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
	for i, col := range cols {
		hc, ok := HeaderColumns[col]
		if !ok {
			continue
		}
		width := ct.GetTableColumnWidth(col)
		if width == 0 {
			continue
		}
		arrow := " "
		colorfn := baseColor
		if !noSort {
			if ct.State.sortBy == col {
				colorfn = ct.colorscheme.TableHeaderColumnActiveSprintf()
				if ct.State.sortDesc {
					arrow = ArrowDown
				} else {
					arrow = ArrowUp
				}
			}
		}
		label := hc.Label
		if noSort {
			label = hc.PlainLabel
		}
		leftAlign := ct.GetTableColumnAlignLeft(col)
		switch col {
		case "price", "balance":
			label = ct.CurrencySymbol() + label
		}
		if leftAlign {
			label = label + arrow
		} else {
			label = arrow + label
		}
		padfn := pad.Left
		padLeft := 1
		if !noSort && i == 0 {
			padLeft = 0
		}
		if leftAlign {
			padfn = pad.Right
		}
		colStr := fmt.Sprintf(
			"%s%s%s",
			strings.Repeat(" ", padLeft),
			colorfn(padfn(label, width+(1-padLeft), " ")),
			strings.Repeat(" ", 1),
		)
		headers = append(headers, colStr)
	}

	ct.UpdateUI(func() error {
		return ct.Views.TableHeader.Update(strings.Join(headers, ""))
	})

	return nil
}

// SetTableColumnAlignLeft sets the column alignment direction for header
func (ct *Cointop) SetTableColumnAlignLeft(header string, alignLeft bool) {
	ct.State.tableColumnAlignLeft.Store(header, alignLeft)
}

// GetTableColumnAlignLeft gets the column alignment direction for header
func (ct *Cointop) GetTableColumnAlignLeft(header string) bool {
	ifc, ok := ct.State.tableColumnAlignLeft.Load(header)
	if ok {
		return ifc.(bool)
	}
	return false
}

// SetTableColumnWidth sets the column width for header
func (ct *Cointop) SetTableColumnWidth(header string, width int) {
	prevIfc, ok := ct.State.tableColumnWidths.Load(header)
	var prev int
	if ok {
		prev = prevIfc.(int)
	} else {
		hc := HeaderColumns[header]
		prev = utf8.RuneCountInString(hc.Label) + 1
		switch header {
		case "price", "balance":
			prev++
		}
	}

	ct.State.tableColumnWidths.Store(header, int(math.Max(float64(width), float64(prev))))
}

// SetTableColumnWidthFromString sets the column width for header given size of string
func (ct *Cointop) SetTableColumnWidthFromString(header string, text string) {
	ct.SetTableColumnWidth(header, utf8.RuneCountInString(text))
}

// GetTableColumnWidth gets the column width for header
func (ct *Cointop) GetTableColumnWidth(header string) int {
	ifc, ok := ct.State.tableColumnWidths.Load(header)
	if ok {
		return ifc.(int)
	}
	return 0
}
