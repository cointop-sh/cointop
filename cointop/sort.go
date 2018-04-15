package cointop

import (
	"github.com/bradfitz/slice"
	"github.com/jroimartin/gocui"
)

var colorder = []string{
	"rank",
	"name",
	"symbol",
	"price",
	"marketcap",
	"24hvolume",
	"1hchange",
	"7dchange",
	"totalsupply",
	"availablesupply",
	"lastupdated",
}

func (ct *Cointop) sort(sortby string, desc bool, list []*coin) {
	ct.sortby = sortby
	ct.sortdesc = desc
	slice.Sort(list[:], func(i, j int) bool {
		if ct.sortdesc {
			i, j = j, i
		}
		a := list[i]
		b := list[j]
		switch sortby {
		case "rank":
			return a.Rank < b.Rank
		case "name":
			return a.Name < b.Name
		case "symbol":
			return a.Symbol < b.Symbol
		case "price":
			return a.PriceUSD < b.PriceUSD
		case "marketcap":
			return a.MarketCapUSD < b.MarketCapUSD
		case "24hvolume":
			return a.USD24HVolume < b.USD24HVolume
		case "1hchange":
			return a.PercentChange1H < b.PercentChange1H
		case "24hchange":
			return a.PercentChange24H < b.PercentChange24H
		case "7dchange":
			return a.PercentChange7D < b.PercentChange7D
		case "totalsupply":
			return a.TotalSupply < b.TotalSupply
		case "availablesupply":
			return a.AvailableSupply < b.AvailableSupply
		case "lastupdated":
			return a.LastUpdated < b.LastUpdated
		default:
			return a.Rank < b.Rank
		}
	})
	ct.updateHeaders()
}

func (ct *Cointop) sortfn(sortby string, desc bool) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if ct.sortby == sortby {
			desc = !desc
		}

		ct.sort(sortby, desc, ct.coins)
		ct.update(func() {
			ct.tableview.Clear()
			ct.updateTable()
		})

		ct.rowChanged()
		return nil
	}
}

func (ct *Cointop) getSortColIndex() int {
	for i, col := range colorder {
		if ct.sortby == col {
			return i
		}
	}
	return 0
}

func (ct *Cointop) sortAsc() error {
	ct.sortdesc = false
	ct.sort(ct.sortby, ct.sortdesc, ct.coins)
	ct.update(func() {
		ct.tableview.Clear()
		ct.updateTable()
	})
	ct.rowChanged()
	return nil
}

func (ct *Cointop) sortDesc() error {
	ct.sortdesc = true
	ct.sort(ct.sortby, ct.sortdesc, ct.coins)
	ct.update(func() {
		ct.tableview.Clear()
		ct.updateTable()
	})
	ct.rowChanged()
	return nil
}

func (ct *Cointop) sortPrevCol() error {
	nextsortby := colorder[0]
	i := ct.getSortColIndex()
	k := i - 1
	if k < 0 {
		k = 0
	}
	nextsortby = colorder[k]
	ct.sort(nextsortby, ct.sortdesc, ct.coins)
	ct.update(func() {
		ct.tableview.Clear()
		ct.updateTable()
	})
	ct.rowChanged()
	return nil
}

func (ct *Cointop) sortNextCol() error {
	nextsortby := colorder[0]
	l := len(colorder)
	i := ct.getSortColIndex()
	k := i + 1
	if k > l-1 {
		k = l - 1
	}
	nextsortby = colorder[k]
	ct.sort(nextsortby, ct.sortdesc, ct.coins)
	ct.update(func() {
		ct.tableview.Clear()
		ct.updateTable()
	})
	ct.rowChanged()
	return nil
}
