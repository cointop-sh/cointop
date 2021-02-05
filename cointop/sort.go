package cointop

import (
	"sort"
	"sync"

	"github.com/miguelmota/gocui"
)

var sortlock sync.Mutex

// Sort sorts the list of coins
func (ct *Cointop) Sort(sortBy string, desc bool, list []*Coin, renderHeaders bool) {
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
		case "market_cap":
			return a.MarketCap < b.MarketCap
		case "24h_volume":
			return a.Volume24H < b.Volume24H
		case "1h_change":
			return a.PercentChange1H < b.PercentChange1H
		case "24h_change":
			return a.PercentChange24H < b.PercentChange24H
		case "7d_change":
			return a.PercentChange7D < b.PercentChange7D
		case "total_supply":
			return a.TotalSupply < b.TotalSupply
		case "available_supply":
			return a.AvailableSupply < b.AvailableSupply
		case "last_updated":
			return a.LastUpdated < b.LastUpdated
		default:
			return a.Rank < b.Rank
		}
	})

	if renderHeaders {
		ct.UpdateTableHeader()
	}
}

// SortAsc sorts list of coins in ascending order
func (ct *Cointop) SortAsc() error {
	ct.debuglog("sortAsc()")
	ct.State.sortDesc = false
	ct.UpdateTable()
	return nil
}

// SortDesc sorts list of coins in descending order
func (ct *Cointop) SortDesc() error {
	ct.debuglog("sortDesc()")
	ct.State.sortDesc = true
	ct.UpdateTable()
	return nil
}

// SortPrevCol sorts the previous column
func (ct *Cointop) SortPrevCol() error {
	ct.debuglog("sortPrevCol()")
	i := ct.GetSortColIndex()
	k := i - 1
	if k < 0 {
		k = 0
	}

	nextsortBy := ct.TableColumnOrder[k]
	ct.Sort(nextsortBy, ct.State.sortDesc, ct.State.coins, true)
	ct.UpdateTable()
	return nil
}

// SortNextCol sorts the next column
func (ct *Cointop) SortNextCol() error {
	ct.debuglog("sortNextCol()")
	l := len(ct.TableColumnOrder)
	i := ct.GetSortColIndex()
	k := i + 1
	if k > l-1 {
		k = l - 1
	}

	nextsortBy := ct.TableColumnOrder[k]
	ct.Sort(nextsortBy, ct.State.sortDesc, ct.State.coins, true)
	ct.UpdateTable()
	return nil
}

// SortToggle toggles the sort order
func (ct *Cointop) SortToggle(sortBy string, desc bool) error {
	ct.debuglog("sortToggle()")
	if ct.State.sortBy == sortBy {
		desc = !ct.State.sortDesc
	}

	ct.Sort(sortBy, desc, ct.State.coins, true)
	ct.UpdateTable()
	return nil
}

// Sortfn returns the sort function as a wrapped gocui keybinding function
func (ct *Cointop) Sortfn(sortBy string, desc bool) func(g *gocui.Gui, v *gocui.View) error {
	ct.debuglog("sortfn()")
	return func(g *gocui.Gui, v *gocui.View) error {
		return ct.SortToggle(sortBy, desc)
	}
}

// GetSortColIndex gets the sort column index
func (ct *Cointop) GetSortColIndex() int {
	ct.debuglog("getSortColIndex()")
	for i, col := range ct.TableColumnOrder {
		if ct.State.sortBy == col {
			return i
		}
	}
	return 0
}
