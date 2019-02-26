package cointop

import (
	"fmt"
	"strings"

	"github.com/miguelmota/cointop/cointop/common/color"
)

func (ct *Cointop) updateHeaders() {
	var cols []string

	type t struct {
		colorfn     func(a ...interface{}) string
		displaytext string
		padleft     int
		padright    int
		arrow       string
	}

	cm := map[string]*t{
		"rank":            &t{color.Black, "[r]ank", 0, 1, " "},
		"name":            &t{color.Black, "[n]ame", 0, 11, " "},
		"symbol":          &t{color.Black, "[s]ymbol", 4, 0, " "},
		"price":           &t{color.Black, "[p]rice", 2, 0, " "},
		"holdings":        &t{color.Black, "[h]oldings", 5, 0, " "},
		"balance":         &t{color.Black, "[b]alance", 5, 0, " "},
		"marketcap":       &t{color.Black, "[m]arket cap", 5, 0, " "},
		"24hvolume":       &t{color.Black, "24H [v]olume", 3, 0, " "},
		"1hchange":        &t{color.Black, "[1]H%", 5, 0, " "},
		"24hchange":       &t{color.Black, "[2]4H%", 3, 0, " "},
		"7dchange":        &t{color.Black, "[7]D%", 4, 0, " "},
		"totalsupply":     &t{color.Black, "[t]otal supply", 7, 0, " "},
		"availablesupply": &t{color.Black, "[a]vailable supply", 0, 0, " "},
		"percentholdings": &t{color.Black, "%holdings", 2, 0, " "},
		"lastupdated":     &t{color.Black, "last [u]pdated", 3, 0, " "},
	}

	for k := range cm {
		cm[k].arrow = " "
		if ct.sortby == k {
			cm[k].colorfn = color.CyanBg
			if ct.sortdesc {
				cm[k].arrow = "▼"
			} else {
				cm[k].arrow = "▲"
			}
		}
	}

	if ct.portfoliovisible {
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
			d = s.arrow + ct.currencySymbol() + s.displaytext
		}

		str = fmt.Sprintf(
			"%s%s%s",
			strings.Repeat(" ", s.padleft),
			s.colorfn(d),
			strings.Repeat(" ", s.padright),
		)
		headers = append(headers, str)
	}

	ct.update(func() {
		ct.headersview.Clear()
		fmt.Fprintln(ct.headersview, strings.Join(headers, ""))
	})
}
