package cointop

import (
	"fmt"
	"strings"
)

// TableHeaderView is structure for table header view
type TableHeaderView struct {
	*View
}

// NewTableHeaderView returns a new table header view
func NewTableHeaderView() *TableHeaderView {
	return &TableHeaderView{NewView("header")}
}

// UpdateTableHeader renders the table header
func (ct *Cointop) UpdateTableHeader() {
	ct.debuglog("UpdateTableHeader()")
	var cols []string

	type t struct {
		colorfn     func(a ...interface{}) string
		displaytext string
		padleft     int
		padright    int
		arrow       string
	}

	baseColor := ct.colorscheme.TableHeaderSprintf()
	cm := map[string]*t{
		"rank":            &t{baseColor, "[r]ank", 0, 1, " "},
		"name":            &t{baseColor, "[n]ame", 0, 11, " "},
		"symbol":          &t{baseColor, "[s]ymbol", 4, 0, " "},
		"price":           &t{baseColor, "[p]rice", 2, 0, " "},
		"holdings":        &t{baseColor, "[h]oldings", 5, 0, " "},
		"balance":         &t{baseColor, "[b]alance", 5, 0, " "},
		"marketcap":       &t{baseColor, "[m]arket cap", 5, 0, " "},
		"24hvolume":       &t{baseColor, "24H [v]olume", 3, 0, " "},
		"1hchange":        &t{baseColor, "[1]H%", 5, 0, " "},
		"24hchange":       &t{baseColor, "[2]4H%", 3, 0, " "},
		"7dchange":        &t{baseColor, "[7]D%", 4, 0, " "},
		"totalsupply":     &t{baseColor, "[t]otal supply", 7, 0, " "},
		"availablesupply": &t{baseColor, "[a]vailable supply", 0, 0, " "},
		"percentholdings": &t{baseColor, "%holdings", 2, 0, " "},
		"lastupdated":     &t{baseColor, "last [u]pdated", 3, 0, " "},
	}

	for k := range cm {
		cm[k].arrow = " "
		if ct.State.sortBy == k {
			cm[k].colorfn = ct.colorscheme.TableHeaderColumnActiveSprintf()
			if ct.State.sortDesc {
				cm[k].arrow = "▼"
			} else {
				cm[k].arrow = "▲"
			}
		}
	}

	if ct.State.portfolioVisible {
		cols = []string{"rank", "name", "symbol", "price",
			"holdings", "balance", "24hchange", "percentholdings", "lastupdated"}
	} else {
		cols = []string{"rank", "name", "symbol", "price",
			"marketcap", "24hvolume", "1hchange", "24hchange",
			"7dchange", "totalsupply", "availablesupply", "lastupdated"}
	}

	var headers []string
	for _, v := range cols {
		s, ok := cm[v]
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

	ct.Update(func() error {
		if ct.Views.TableHeader.Backing() == nil {
			return nil
		}

		ct.Views.TableHeader.Backing().Clear()
		fmt.Fprintln(ct.Views.TableHeader.Backing(), strings.Join(headers, ""))
		return nil
	})
}
