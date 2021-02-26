package cointop

import (
	"fmt"
	"math"
	"net/url"
	"strings"

	"github.com/miguelmota/cointop/pkg/ui"
)

// TableView is structure for table view
type TableView = ui.View

// NewTableView returns a new table view
func NewTableView() *TableView {
	var view *TableView = ui.NewView("table")
	return view
}

const dots = "..."

// RefreshTable refreshes the table
func (ct *Cointop) RefreshTable() error {
	ct.debuglog("refreshTable()")

	statusText := ""
	switch ct.State.selectedView {
	case PortfolioView:
		ct.table = ct.GetPortfolioTable()
		if ct.table.RowCount() == 0 {
			statusText = "No holdings found. Press \"e\" on a coin to edit holdings."
		}
	case PriceAlertsView:
		ct.table = ct.GetPriceAlertsTable()
		if ct.table.RowCount() == 0 {
			statusText = "No price alerts found. Press \"+\" on a coin to add a price alert."
		}
	default:
		ct.table = ct.GetCoinsTable()
		if ct.table.RowCount() == 0 {
			statusText = "no coin data"
		}
	}
	ct.table.HideColumHeaders = true

	ct.UpdateUI(func() error {
		ct.Views.Table.Clear()
		if statusText == "" {
			if ct.Views.Table.HasBacking() {
				ct.table.Format().Fprint(ct.Views.Table.Backing())
			}
		} else {
			ct.Views.Table.Update(fmt.Sprintf("\n\n%s", statusText))
		}
		go ct.RowChanged()
		go ct.UpdateTableHeader()
		go ct.UpdateMarketbar()
		go ct.UpdateChart()
		return nil
	})

	return nil
}

// UpdateTable updates the table
func (ct *Cointop) UpdateTable() error {
	ct.debuglog("UpdateTable()")
	ct.State.allCoinsSlugMap.Range(func(key, value interface{}) bool {
		k := key.(string)
		if v, ok := value.(*Coin); ok {
			v.Favorite = ct.State.favorites[v.Name]
			ct.State.allCoinsSlugMap.Store(k, v)
		}

		return true
	})

	if ct.IsFavoritesVisible() {
		ct.State.coins = ct.GetFavoritesSlice()
	} else if ct.IsPortfolioVisible() {
		ct.State.coins = ct.GetPortfolioSlice()
	} else {
		// TODO: maintain state of previous sorting
		if ct.State.sortBy == "holdings" {
			ct.State.sortBy = "rank"
			ct.State.sortDesc = false
		}

		ct.State.coins = ct.GetTableCoinsSlice()
	}

	ct.Sort(ct.State.sortBy, ct.State.sortDesc, ct.State.coins, true)
	go ct.RefreshTable()
	return nil
}

// GetTableCoinsSlice returns a slice of the table rows
func (ct *Cointop) GetTableCoinsSlice() []*Coin {
	ct.debuglog("GetTableCoinsSlice()")
	sliced := []*Coin{}
	start := ct.State.page * ct.State.perPage
	end := start + ct.State.perPage
	allCoins := ct.AllCoins()
	size := len(allCoins)
	if start < 0 {
		start = 0
	}
	if end >= size-1 {
		start = int(float64(start/100) * 100)
		end = size - 1
	}
	if start < 0 {
		start = 0
	}
	if end >= size {
		end = size - 1
	}
	if end < 0 {
		end = 0
	}
	if start >= end {
		return nil
	}
	if end > 0 {
		sliced = allCoins[start:end]

		// NOTE: restore rank
		for _, coin := range sliced {
			icoin, _ := ct.State.allCoinsSlugMap.Load(coin.Name)
			if icoin != nil {
				c, _ := icoin.(*Coin)
				coin.Rank = c.Rank
			}
		}
	}

	return sliced
}

// HighlightedRowIndex returns the index of the highlighted row within the per-page limit
func (ct *Cointop) HighlightedRowIndex() int {
	ct.debuglog("HighlightedRowIndex()")
	oy := ct.Views.Table.OriginY()
	cy := ct.Views.Table.CursorY()
	idx := oy + cy
	l := ct.TableRowsLen()
	if idx >= l {
		idx = l - 1
	}
	if idx < 0 {
		idx = 0
	}
	return idx
}

// HighlightedRowCoin returns the coin at the index of the highlighted row
func (ct *Cointop) HighlightedRowCoin() *Coin {
	ct.debuglog("HighlightedRowCoin()")
	idx := ct.HighlightedRowIndex()
	if len(ct.State.coins) == 0 {
		return nil
	}
	return ct.State.coins[idx]
}

// HighlightedPageRowIndex returns the index of page row of the highlighted row
func (ct *Cointop) HighlightedPageRowIndex() int {
	ct.debuglog("HighlightedPageRowIndex()")
	cy := ct.Views.Table.CursorY()
	idx := cy
	if idx < 0 {
		idx = 0
	}

	return idx
}

// RowLink returns the row url link
func (ct *Cointop) RowLink() string {
	ct.debuglog("RowLink()")
	coin := ct.HighlightedRowCoin()
	if coin == nil {
		return ""
	}

	return ct.api.CoinLink(coin.Name)
}

// RowLinkShort returns a shortened version of the row url link
func (ct *Cointop) RowLinkShort() string {
	ct.debuglog("RowLinkShort()")
	link := ct.RowLink()
	if link != "" {
		u, err := url.Parse(link)
		if err != nil {
			return ""
		}

		host := u.Hostname()
		host = strings.Replace(host, "www.", "", -1)
		path := u.EscapedPath()
		parts := strings.Split(path, "/")
		if len(parts) > 0 {
			path = parts[len(parts)-1]
		}

		return fmt.Sprintf("http://%s/%s/%s", host, dots, path)
	}

	return ""
}

// ToggleTableFullscreen toggles the table fullscreen mode
func (ct *Cointop) ToggleTableFullscreen() error {
	ct.debuglog("ToggleTableFullscreen()")
	ct.State.onlyTable = !ct.State.onlyTable
	if !ct.State.onlyTable {
		// NOTE: cached values are initial config settings.
		// If the only-table config was set then toggle
		// all other initial hidden views.
		onlyTable, _ := ct.cache.Get("onlyTable")

		if onlyTable.(bool) {
			ct.State.hideMarketbar = false
			ct.State.hideChart = false
			ct.State.hideStatusbar = false
		} else {
			// NOTE: cached values store initial hidden views preferences.
			hideMarketbar, _ := ct.cache.Get("hideMarketbar")
			ct.State.hideMarketbar = hideMarketbar.(bool)
			hideChart, _ := ct.cache.Get("hideChart")
			ct.State.hideChart = hideChart.(bool)
			hideStatusbar, _ := ct.cache.Get("hideStatusbar")
			ct.State.hideStatusbar = hideStatusbar.(bool)
		}
	}

	return nil
}

// SetSelectedView sets the active table view
func (ct *Cointop) SetSelectedView(viewName string) {
	ct.State.lastSelectedView = ct.State.selectedView
	ct.State.selectedView = viewName
}

// ToggleSelectedView toggles between current table view and last selected table view
func (ct *Cointop) ToggleSelectedView(viewName string) {
	if !(ct.IsPortfolioVisible() || ct.IsFavoritesVisible()) {
		ct.State.lastSelectedRowIndex = ct.HighlightedPageRowIndex()
	}
	if ct.State.lastSelectedView == "" || ct.State.selectedView != viewName {
		ct.SetSelectedView(viewName)
	} else {
		ct.SetSelectedView(ct.State.lastSelectedView)
	}

	l := ct.TableRowsLen()
	if ct.IsPortfolioVisible() || ct.IsFavoritesVisible() {
		// highlight last row if current row is out of bounds (can happen when switching views).
		currentRowIdx := ct.HighlightedRowIndex()
		if currentRowIdx >= l-1 {
			ct.HighlightRow(l - 1)
		}
	} else {
		ct.GoToPageRowIndex(int(math.Min(float64(l-1), float64(ct.State.lastSelectedRowIndex))))
	}
}
