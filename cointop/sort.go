package cointop

import (
	"sort"
	"sync"

	"github.com/cointop-sh/cointop/pkg/gocui"
	log "github.com/sirupsen/logrus"
)

var sortlock sync.Mutex

// Sort sorts the list of coins
func (ct *Cointop) Sort(sortCons *sortConstraint, list []*Coin, renderHeaders bool) {
	log.Debug("Sort()")
	sortlock.Lock()
	defer sortlock.Unlock()

	ct.State.viewSorts[ct.State.selectedView] = sortCons
	if list == nil {
		return
	}
	if len(list) < 2 {
		return
	}
	sort.SliceStable(list[:], func(i, j int) bool {
		if sortCons.sortDesc {
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
		switch sortCons.sortBy {
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
		case "30d_change":
			return a.PercentChange30D < b.PercentChange30D
		case "1y_change":
			return a.PercentChange1Y < b.PercentChange1Y
		case "total_supply":
			return a.TotalSupply < b.TotalSupply
		case "available_supply":
			return a.AvailableSupply < b.AvailableSupply
		case "last_updated":
			return a.LastUpdated < b.LastUpdated
		case "cost_price":
			return a.BuyPrice < b.BuyPrice
		case "cost":
			return (a.BuyPrice * a.Holdings) < (b.BuyPrice * b.Holdings) // TODO: convert?
		case "pnl":
			return (a.Price - a.BuyPrice) < (b.Price - b.BuyPrice)
		case "pnl_percent":
			return (a.Price - a.BuyPrice) < (b.Price - b.BuyPrice)
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
	log.Debug("SortAsc()")
	ct.State.viewSorts[ct.State.selectedView].sortDesc = false
	ct.UpdateTable()
	return nil
}

// SortDesc sorts list of coins in descending order
func (ct *Cointop) SortDesc() error {
	log.Debug("SortDesc()")
	ct.State.viewSorts[ct.State.selectedView].sortDesc = true
	ct.UpdateTable()
	return nil
}

// SortPrevCol sorts the previous column
func (ct *Cointop) SortPrevCol() error {
	log.Debug("SortPrevCol()")
	cols := ct.GetActiveTableHeaders()
	i := ct.GetSortColIndex()
	k := i - 1
	if k < 0 {
		k = 0
	}
	nextsortBy := cols[k]

	curSortConst := ct.State.viewSorts[ct.State.selectedView]
	curSortConst.sortBy = nextsortBy
	ct.Sort(curSortConst, ct.State.coins, true)
	ct.UpdateTable()
	return nil
}

// SortNextCol sorts the next column
func (ct *Cointop) SortNextCol() error {
	log.Debug("SortNextCol()")
	cols := ct.GetActiveTableHeaders()
	l := len(cols)
	i := ct.GetSortColIndex()
	k := i + 1
	if k > l-1 {
		k = l - 1
	}
	nextsortBy := cols[k]
	curSortCons := ct.State.viewSorts[ct.State.selectedView]
	curSortCons.sortBy = nextsortBy
	ct.Sort(curSortCons, ct.State.coins, true)
	ct.UpdateTable()
	return nil
}

// SortToggle toggles the sort order
func (ct *Cointop) SortToggle(sortBy string, desc bool) error {
	log.Debug("SortToggle()")
	curSortCons := ct.State.viewSorts[ct.State.selectedView]
	if curSortCons.sortBy == sortBy {
		curSortCons.sortDesc = !curSortCons.sortDesc
	} else {
		curSortCons.sortBy = sortBy
		curSortCons.sortDesc = desc
	}

	ct.Sort(curSortCons, ct.State.coins, true)
	ct.UpdateTable()
	return nil
}

// Sortfn returns the sort function as a wrapped gocui keybinding function
func (ct *Cointop) Sortfn(sortBy string, desc bool) func(g *gocui.Gui, v *gocui.View) error {
	log.Debug("Sortfn()")
	return func(g *gocui.Gui, v *gocui.View) error {
		coin := ct.HighlightedRowCoin()
		err := ct.SortToggle(sortBy, desc)
		if err != nil {
			return err
		}
		if ct.State.keepRowFocusOnSort {
			err = ct.GoToCoinRow(coin)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// GetSortColIndex gets the sort column index
func (ct *Cointop) GetSortColIndex() int {
	log.Debug("GetSortColIndex()")
	cols := ct.GetActiveTableHeaders()
	for i, col := range cols {
		if ct.State.viewSorts[ct.State.selectedView].sortBy == col {
			return i
		}
	}
	return 0
}
