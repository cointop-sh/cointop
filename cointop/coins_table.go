package cointop

import (
	"fmt"
	"strconv"
	"time"

	"github.com/miguelmota/cointop/pkg/humanize"
	"github.com/miguelmota/cointop/pkg/table"
)

// DefaultCoinTableHeaders are the default coin table header columns
var DefaultCoinTableHeaders = []string{
	"rank",
	"name",
	"symbol",
	"price",
	"marketcap",
	"24h_volume",
	"1h_change",
	"24h_change",
	"7d_change",
	"total_supply",
	"available_supply",
	"last_updated",
}

// ValidCoinsTableHeader returns true if it's a valid table header name
func (ct *Cointop) ValidCoinsTableHeader(name string) bool {
	for _, v := range DefaultCoinTableHeaders {
		if v == name {
			return true
		}
	}

	return false
}

// GetCoinsTableHeaders returns the coins table headers
func (ct *Cointop) GetCoinsTableHeaders() []string {
	return ct.State.coinsTableColumns
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

		headers := ct.GetCoinsTableHeaders()
		var rowCells []*table.RowCell
		for _, header := range headers {
			switch header {
			case "rank":
				rowCells = append(rowCells, &table.RowCell{
					LeftMargin: 0,
					Width:      6,
					LeftAlign:  false,
					Color:      ct.colorscheme.Default,
					Text:       rank,
				})
			case "name":
				rowCells = append(rowCells, &table.RowCell{
					LeftMargin: 1,
					Width:      22,
					LeftAlign:  true,
					Color:      namecolor,
					Text:       name,
				})
			case "symbol":
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin: 1,
						Width:      symbolpadding,
						LeftAlign:  true,
						Color:      ct.colorscheme.TableRow,
						Text:       symbol,
					})
			case "price":
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin: 1,
						Width:      12,
						LeftAlign:  false,
						Color:      ct.colorscheme.TableColumnPrice,
						Text:       humanize.Commaf(coin.Price),
					})
			case "marketcap":
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin: 1,
						Width:      18,
						LeftAlign:  false,
						Color:      ct.colorscheme.TableRow,
						Text:       humanize.Commaf(coin.MarketCap),
					})
			case "24h_volume":
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin: 1,
						Width:      16,
						LeftAlign:  false,
						Color:      ct.colorscheme.TableRow,
						Text:       humanize.Commaf(coin.Volume24H),
					})
			case "1h_change":
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin: 1,
						Width:      11,
						LeftAlign:  false,
						Color:      color1h,
						Text:       fmt.Sprintf("%.2f%%", coin.PercentChange1H),
					})
			case "24h_change":
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin: 1,
						Width:      10,
						LeftAlign:  false,
						Color:      color24h,
						Text:       fmt.Sprintf("%.2f%%", coin.PercentChange24H),
					})
			case "7d_change":
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin: 1,
						Width:      10,
						LeftAlign:  false,
						Color:      color7d,
						Text:       fmt.Sprintf("%.2f%%", coin.PercentChange7D),
					})
			case "total_supply":
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin: 1,
						Width:      22,
						LeftAlign:  false,
						Color:      ct.colorscheme.TableRow,
						Text:       humanize.Commaf(coin.TotalSupply),
					})
			case "available_supply":
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin: 1,
						Width:      19,
						LeftAlign:  false,
						Color:      ct.colorscheme.TableRow,
						Text:       humanize.Commaf(coin.AvailableSupply),
					})
			case "last_updated":
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin: 1,
						Width:      18,
						LeftAlign:  false,
						Color:      ct.colorscheme.TableRow,
						Text:       lastUpdated,
					})
			}
		}

		t.AddRowCells(
			rowCells...,
		// TODO: add %percent of cap
		)
	}

	return t
}
