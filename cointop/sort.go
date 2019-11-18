package cointop

import (
	"sort"
	"sync"

	"github.com/jroimartin/gocui"
)

var sortlock sync.Mutex

func (ct *Cointop) sort(sortBy string, desc bool, list []*Coin, renderHeaders bool) {
	ct.debuglog("sort()")
	sortlock.Lock()
	defer sortlock.Unlock()
	if list == nil {
		return
	}
	if len(list) < 2 {
		return
	}
	ct.State.sortBy = sortBy
	ct.State.sortDesc = desc
	sort.Slice(list[:], func(i, j int) bool {
		if ct.State.sortDesc {
			i, j = j, i
		}
		a := list[i]
		b := list[j]
		if a == nil {
			return true
		}
		if b == nil {
			return false
		}
		switch sortBy {
		case "rank":
			return a.Rank < b.Rank
		case "name":
			return a.Name < b.Name
		case "symbol":
			return a.Symbol < b.Symbol
		case "price":
			return a.Price < b.Price
		case "holdings":
			return a.Holdings < b.Holdings
		case "balance":
			return a.Balance < b.Balance
		case "marketcap":
			return a.MarketCap < b.MarketCap
		case "24hvolume":
			return a.Volume24H < b.Volume24H
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

	if renderHeaders {
		ct.UpdateTableHeader()
	}
}

func (ct *Cointop) sortAsc() error {
	ct.debuglog("sortAsc()")
	ct.State.sortDesc = false
	ct.UpdateTable()
	return nil
}

func (ct *Cointop) sortDesc() error {
	ct.debuglog("sortDesc()")
	ct.State.sortDesc = true
	ct.UpdateTable()
	return nil
}

func (ct *Cointop) sortPrevCol() error {
	ct.debuglog("sortPrevCol()")
	nextsortBy := ct.TableColumnOrder[0]
	i := ct.getSortColIndex()
	k := i - 1
	if k < 0 {
		k = 0
	}

	nextsortBy = ct.TableColumnOrder[k]
	ct.sort(nextsortBy, ct.State.sortDesc, ct.State.coins, true)
	ct.UpdateTable()
	return nil
}

func (ct *Cointop) sortNextCol() error {
	ct.debuglog("sortNextCol()")
	nextsortBy := ct.TableColumnOrder[0]
	l := len(ct.TableColumnOrder)
	i := ct.getSortColIndex()
	k := i + 1
	if k > l-1 {
		k = l - 1
	}

	nextsortBy = ct.TableColumnOrder[k]
	ct.sort(nextsortBy, ct.State.sortDesc, ct.State.coins, true)
	ct.UpdateTable()
	return nil
}

func (ct *Cointop) sortToggle(sortBy string, desc bool) error {
	ct.debuglog("sortToggle()")
	if ct.State.sortBy == sortBy {
		desc = !ct.State.sortDesc
	}

	ct.sort(sortBy, desc, ct.State.coins, true)
	ct.UpdateTable()
	return nil
}

func (ct *Cointop) sortfn(sortBy string, desc bool) func(g *gocui.Gui, v *gocui.View) error {
	ct.debuglog("sortfn()")
	return func(g *gocui.Gui, v *gocui.View) error {
		return ct.sortToggle(sortBy, desc)
	}
}

func (ct *Cointop) getSortColIndex() int {
	ct.debuglog("getSortColIndex()")
	for i, col := range ct.TableColumnOrder {
		if ct.State.sortBy == col {
			return i
		}
	}
	return 0
}
