package cointop

import (
	"fmt"
	"strings"

	"github.com/miguelmota/cointop/pkg/color"
)

func (ct *Cointop) updateHeaders() {
	cm := map[string]func(a ...interface{}) string{
		"rank":            color.Black,
		"name":            color.Black,
		"symbol":          color.Black,
		"price":           color.Black,
		"marketcap":       color.Black,
		"24hvolume":       color.Black,
		"1hchange":        color.Black,
		"24hchange":       color.Black,
		"7dchange":        color.Black,
		"totalsupply":     color.Black,
		"availablesupply": color.Black,
		"lastupdated":     color.Black,
	}
	sm := map[string]string{
		"rank":            " ",
		"name":            " ",
		"symbol":          " ",
		"price":           " ",
		"marketcap":       " ",
		"24hvolume":       " ",
		"1hchange":        " ",
		"24hchange":       " ",
		"7dchange":        " ",
		"totalsupply":     " ",
		"availablesupply": " ",
		"lastupdated":     " ",
	}
	for k := range cm {
		if ct.sortby == k {
			cm[k] = color.CyanBg
			if ct.sortdesc {
				sm[k] = "▼"
			} else {
				sm[k] = "▲"
			}
		}
	}
	headers := []string{
		fmt.Sprintf("%s%s", cm["rank"](sm["rank"]+"[r]ank"), strings.Repeat(" ", 1)),
		fmt.Sprintf("%s%s", cm["name"](sm["name"]+"[n]ame"), strings.Repeat(" ", 15)),
		fmt.Sprintf("%s%s", cm["symbol"](sm["symbol"]+"[s]ymbol"), strings.Repeat(" ", 1)),
		fmt.Sprintf("%s%s", strings.Repeat(" ", 1), cm["price"](sm["price"]+"[p]rice")),
		fmt.Sprintf("%s%s", strings.Repeat(" ", 5), cm["marketcap"](sm["marketcap"]+"[m]arket cap")),
		fmt.Sprintf("%s%s", strings.Repeat(" ", 3), cm["24hvolume"](sm["24hvolume"]+"24H [v]olume")),
		fmt.Sprintf("%s%s", strings.Repeat(" ", 4), cm["1hchange"](sm["1hchange"]+"[1]H%")),
		fmt.Sprintf("%s%s", strings.Repeat(" ", 3), cm["24hchange"](sm["24hchange"]+"[2]4H%")),
		fmt.Sprintf("%s%s", strings.Repeat(" ", 3), cm["7dchange"](sm["7dchange"]+"[7]DH%")),
		fmt.Sprintf("%s%s", strings.Repeat(" ", 6), cm["totalsupply"](sm["totalsupply"]+"[t]otal supply")),
		fmt.Sprintf("%s%s", strings.Repeat(" ", 1), cm["availablesupply"](sm["availablesupply"]+"[a]vailable supply")),
		fmt.Sprintf("%s%s", strings.Repeat(" ", 4), cm["lastupdated"](sm["lastupdated"]+"last [u]pdated")),
	}

	ct.update(func() {
		ct.headersview.Clear()
		fmt.Fprintln(ct.headersview, strings.Join(headers, ""))
	})
}
