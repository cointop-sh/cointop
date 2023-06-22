package cointop

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/cointop-sh/cointop/pkg/ui"
	log "github.com/sirupsen/logrus"
)

// TableView is structure for table view
type TableView = ui.View

// NewTableView returns a new table view
func NewTableView() *TableView {
	return ui.NewView("table")
}

const dots = "..."

// RefreshTable refreshes the table
func (ct *Cointop) RefreshTable() error {
	log.Debug("RefreshTable()")

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
	log.Debug("UpdateTable()")
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
		ct.State.coins = ct.GetTableCoinsSlice()
	}

	ct.Sort(ct.State.viewSorts[ct.State.selectedView], ct.State.coins, true)
	go ct.RefreshTable()
	return nil
}

// GetTableCoinsSlice returns a slice of the table rows
func (ct *Cointop) GetTableCoinsSlice() []*Coin {
	log.Debug("GetTableCoinsSlice()")
	var sliced []*Coin
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
	log.Debug("HighlightedRowIndex()")
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
	log.Debug("HighlightedRowCoin()")
	idx := ct.HighlightedRowIndex()
	coins := ct.State.coins
	if ct.IsPriceAlertsVisible() {
		rows := ct.ActivePriceAlerts()
		for i, row := range rows {
			if i == idx {
				return ct.CoinByName(row.CoinName)
			}
		}
	}
	if len(coins) == 0 {
		return nil
	}
	return coins[idx]
}

// HighlightedPageRowIndex returns the index of page row of the highlighted row
func (ct *Cointop) HighlightedPageRowIndex() int {
	log.Debug("HighlightedPageRowIndex()")
	cy := ct.Views.Table.CursorY()
	idx := cy
	if idx < 0 {
		idx = 0
	}

	return idx
}

// GetLastSelectedRowCoinIndex returns the index of the last selected row coin
func (ct *Cointop) GetLastSelectedRowCoinIndex() int {
	return ct.State.lastSelectedRowIndex
}

// RowLink returns the row url link
func (ct *Cointop) RowLink() string {
	log.Debug("RowLink()")
	coin := ct.HighlightedRowCoin()
	if coin == nil {
		return ""
	}

	// TODO: Can remove this one after some releases
	// because it is a way to force old client refresh coin to have a slug
	if coin.Slug == "" {
		if err := ct.UpdateCoin(coin); err != nil {
			log.Debugf("RowLink() Update coin got err %s", err.Error())
			return ""
		}
	}

	return ct.api.CoinLink(coin.Slug)
}

// RowLink returns the row url link
func (ct *Cointop) RowAltLink() string {
	log.Debug("RowAltLink()")
	coin := ct.HighlightedRowCoin()
	if coin == nil {
		return ""
	}

	return ct.GetAltCoinLink(coin)
}

// RowLinkShort returns a shortened version of the row url link
func (ct *Cointop) RowLinkShort() string {
	log.Debug("RowLinkShort()")
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

func (ct *Cointop) GetAltCoinLink(coin *Coin) string {
	if ct.State.altCoinLink == "" {
		return ct.api.CoinLink(coin.Slug)
	}

	url := ct.State.altCoinLink
	url = strings.Replace(url, "{{ID}}", coin.ID, -1)
	url = strings.Replace(url, "{{NAME}}", coin.Name, -1)
	url = strings.Replace(url, "{{RANK}}", strconv.Itoa(coin.Rank), -1)
	url = strings.Replace(url, "{{SLUG}}", coin.Slug, -1)
	url = strings.Replace(url, "{{SYMBOL}}", coin.Symbol, -1)
	return url
}

// ToggleTableFullscreen toggles the table fullscreen mode
func (ct *Cointop) ToggleTableFullscreen() error {
	log.Debug("ToggleTableFullscreen()")
	ct.State.onlyTable = !ct.State.onlyTable
	ct.State.onlyChart = false
	if !ct.State.onlyTable {
		// NOTE: cached values are initial config settings.
		// If the only-table config was set then toggle
		// all other initial hidden views.
		onlyTable, _ := ct.cache.Get("onlyTable")

		if onlyTable.(bool) {
			ct.State.hideMarketbar = false
			ct.State.hideChart = false
			ct.State.hideTable = false
			ct.State.hideStatusbar = false
		} else {
			// NOTE: cached values store initial hidden views preferences.
			hideMarketbar, _ := ct.cache.Get("hideMarketbar")
			ct.State.hideMarketbar = hideMarketbar.(bool)
			hideChart, _ := ct.cache.Get("hideChart")
			ct.State.hideChart = hideChart.(bool)
			hideTable, _ := ct.cache.Get("hideTable")
			ct.State.hideTable = hideTable.(bool)
			hideStatusbar, _ := ct.cache.Get("hideStatusbar")
			ct.State.hideStatusbar = hideStatusbar.(bool)
		}
	}

	go func() {
		ct.UpdateTable()
		ct.UpdateChart()
	}()

	return nil
}

// SetSelectedView sets the active table view
func (ct *Cointop) SetSelectedView(viewName string) {
	ct.State.lastSelectedView = ct.State.selectedView
	ct.State.selectedView = viewName

	// init sort constraint for the view if it hasn't been seen before
	if _, found := ct.State.viewSorts[viewName]; !found {
		ct.State.viewSorts[viewName] = &sortConstraint{DefaultSortBy, false}
	}
}

// ToggleSelectedView toggles between current table view and last selected table view
func (ct *Cointop) ToggleSelectedView(viewName string) {
	if !(ct.IsPortfolioVisible() || ct.IsFavoritesVisible()) {
		ct.State.lastSelectedRowIndex = ct.HighlightedRowIndex()
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
		ct.HighlightRow(ct.State.lastSelectedRowIndex)
	}
}
