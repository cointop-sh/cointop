package cointop

import (
	"fmt"
	"strconv"
	"time"

	"github.com/miguelmota/cointop/pkg/humanize"
	"github.com/miguelmota/cointop/pkg/table"
)

// GetCoinsTableHeaders returns the coins table headers
func (ct *Cointop) GetCoinsTableHeaders() []string {
	return []string{
		"rank",
		"name",
		"symbol",
		"price",
		"marketcap",
		"24hvolume",
		"1hchange",
		"24hchange",
		"7dchange",
		"totalsupply",
		"availablesupply",
		"lastupdated",
	}
}

// GetCoinsTable returns the table for diplaying the coins
func (ct *Cointop) GetCoinsTable() *table.Table {
	maxX := ct.width()
	t := table.NewTable().SetWidth(maxX)
	for _, coin := range ct.State.coins {
		if coin == nil {
			continue
		}
		star := ct.colorscheme.TableRow(" ")
		if coin.Favorite {
			star = ct.colorscheme.TableRowFavorite("*")
		}
		rank := fmt.Sprintf("%s%v", star, ct.colorscheme.TableRow(fmt.Sprintf("%6v ", coin.Rank)))
		name := TruncateString(coin.Name, 20)
		symbol := TruncateString(coin.Symbol, 6)
		symbolpadding := 8
		// NOTE: this is to adjust padding by 1 because when all name rows are
		// yellow it messes the spacing (need to debug)
		if ct.IsFavoritesVisible() {
			symbolpadding++
		}
		namecolor := ct.colorscheme.TableRow
		color1h := ct.colorscheme.TableColumnChange
		color24h := ct.colorscheme.TableColumnChange
		color7d := ct.colorscheme.TableColumnChange
		if coin.Favorite {
			namecolor = ct.colorscheme.TableRowFavorite
		}
		if coin.PercentChange1H > 0 {
			color1h = ct.colorscheme.TableColumnChangeUp
		}
		if coin.PercentChange1H < 0 {
			color1h = ct.colorscheme.TableColumnChangeDown
		}
		if coin.PercentChange24H > 0 {
			color24h = ct.colorscheme.TableColumnChangeUp
		}
		if coin.PercentChange24H < 0 {
			color24h = ct.colorscheme.TableColumnChangeDown
		}
		if coin.PercentChange7D > 0 {
			color7d = ct.colorscheme.TableColumnChangeUp
		}
		if coin.PercentChange7D < 0 {
			color7d = ct.colorscheme.TableColumnChangeDown
		}
		unix, _ := strconv.ParseInt(coin.LastUpdated, 10, 64)
		lastUpdated := time.Unix(unix, 0).Format("15:04:05 Jan 02")
		t.AddRowCells(
			&table.RowCell{
				LeftMargin: 0,
				Width:      6,
				LeftAlign:  false,
				Color:      ct.colorscheme.Default,
				Text:       rank,
			},
			&table.RowCell{
				LeftMargin: 1,
				Width:      22,
				LeftAlign:  true,
				Color:      namecolor,
				Text:       name,
			},
			&table.RowCell{
				LeftMargin: 1,
				Width:      symbolpadding,
				LeftAlign:  true,
				Color:      ct.colorscheme.TableRow,
				Text:       symbol,
			},
			&table.RowCell{
				LeftMargin: 1,
				Width:      12,
				LeftAlign:  false,
				Color:      ct.colorscheme.TableColumnPrice,
				Text:       humanize.Commaf(coin.Price),
			},
			&table.RowCell{
				LeftMargin: 1,
				Width:      18,
				LeftAlign:  false,
				Color:      ct.colorscheme.TableRow,
				Text:       humanize.Commaf(coin.MarketCap),
			},
			&table.RowCell{
				LeftMargin: 1,
				Width:      16,
				LeftAlign:  false,
				Color:      ct.colorscheme.TableRow,
				Text:       humanize.Commaf(coin.Volume24H),
			},
			&table.RowCell{
				LeftMargin: 1,
				Width:      11,
				LeftAlign:  false,
				Color:      color1h,
				Text:       fmt.Sprintf("%.2f%%", coin.PercentChange1H),
			},
			&table.RowCell{
				LeftMargin: 1,
				Width:      10,
				LeftAlign:  false,
				Color:      color24h,
				Text:       fmt.Sprintf("%.2f%%", coin.PercentChange24H),
			},
			&table.RowCell{
				LeftMargin: 1,
				Width:      10,
				LeftAlign:  false,
				Color:      color7d,
				Text:       fmt.Sprintf("%.2f%%", coin.PercentChange7D),
			},
			&table.RowCell{
				LeftMargin: 1,
				Width:      22,
				LeftAlign:  false,
				Color:      ct.colorscheme.TableRow,
				Text:       humanize.Commaf(coin.TotalSupply),
			},
			&table.RowCell{
				LeftMargin: 1,
				Width:      19,
				LeftAlign:  false,
				Color:      ct.colorscheme.TableRow,
				Text:       humanize.Commaf(coin.AvailableSupply),
			},
			&table.RowCell{
				LeftMargin: 1,
				Width:      18,
				LeftAlign:  false,
				Color:      ct.colorscheme.TableRow,
				Text:       lastUpdated,
			},
			// TODO: add %percent of cap
		)
	}

	return t
}
