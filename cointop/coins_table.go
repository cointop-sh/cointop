package cointop

import (
	"fmt"
	"strconv"
	"time"

	"github.com/miguelmota/cointop/pkg/humanize"
	"github.com/miguelmota/cointop/pkg/table"
)

// SupportedCoinTableHeaders are all the supported coin table header columns
var SupportedCoinTableHeaders = []string{
	"rank",
	"name",
	"symbol",
	"price",
	"1h_change",
	"24h_change",
	"7d_change",
	"30d_change",
	"24h_volume",
	"market_cap",
	"total_supply",
	"available_supply",
	"last_updated",
}

// DefaultCoinTableHeaders are the default coin table header columns
var DefaultCoinTableHeaders = []string{
	"rank",
	"name",
	"symbol",
	"price",
	"1h_change",
	"24h_change",
	"7d_change",
	"24h_volume",
	"market_cap",
	"total_supply",
	"available_supply",
	"last_updated",
}

// ValidCoinsTableHeader returns true if it's a valid table header name
func (ct *Cointop) ValidCoinsTableHeader(name string) bool {
	for _, v := range SupportedCoinTableHeaders {
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
	var rows [][]*table.RowCell
	headers := ct.GetCoinsTableHeaders()
	if ct.IsFavoritesVisible() {
		headers = ct.GetFavoritesTableHeaders()
	}
	ct.ClearSyncMap(ct.State.tableColumnWidths)
	ct.ClearSyncMap(ct.State.tableColumnAlignLeft)
	for _, coin := range ct.State.coins {
		if coin == nil {
			continue
		}
		var rowCells []*table.RowCell
		for _, header := range headers {
			leftMargin := 1
			rightMargin := 1
			switch header {
			case "rank":
				star := ct.colorscheme.TableRow(" ")
				if coin.Favorite {
					star = ct.colorscheme.TableRowFavorite("*")
				}
				rank := fmt.Sprintf("%s%v", star, ct.colorscheme.TableRow(fmt.Sprintf("%6v ", coin.Rank)))
				ct.SetTableColumnWidth(header, 8)
				ct.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells, &table.RowCell{
					LeftMargin:  leftMargin,
					RightMargin: rightMargin,
					LeftAlign:   false,
					Color:       ct.colorscheme.Default,
					Text:        rank,
				})
			case "name":
				name := TruncateString(coin.Name, 16)
				namecolor := ct.colorscheme.TableRow
				if coin.Favorite {
					namecolor = ct.colorscheme.TableRowFavorite
				}
				ct.SetTableColumnWidthFromString(header, name)
				ct.SetTableColumnAlignLeft(header, true)
				rowCells = append(rowCells, &table.RowCell{
					LeftMargin:  leftMargin,
					RightMargin: rightMargin,
					LeftAlign:   true,
					Color:       namecolor,
					Text:        name,
				})
			case "symbol":
				symbol := TruncateString(coin.Symbol, 6)
				ct.SetTableColumnWidthFromString(header, symbol)
				ct.SetTableColumnAlignLeft(header, true)
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin:  leftMargin,
						RightMargin: rightMargin,
						LeftAlign:   true,
						Color:       ct.colorscheme.TableRow,
						Text:        symbol,
					})
			case "price":
				text := humanize.Commaf(coin.Price)
				ct.SetTableColumnWidthFromString(header, text)
				ct.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin:  leftMargin,
						RightMargin: rightMargin,
						LeftAlign:   false,
						Color:       ct.colorscheme.TableColumnPrice,
						Text:        text,
					})
			case "24h_volume":
				text := humanize.Commaf(coin.Volume24H)
				ct.SetTableColumnWidthFromString(header, text)
				ct.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin:  leftMargin,
						RightMargin: rightMargin,
						LeftAlign:   false,
						Color:       ct.colorscheme.TableRow,
						Text:        text,
					})
			case "1h_change":
				color1h := ct.colorscheme.TableColumnChange
				if coin.PercentChange1H > 0 {
					color1h = ct.colorscheme.TableColumnChangeUp
				}
				if coin.PercentChange1H < 0 {
					color1h = ct.colorscheme.TableColumnChangeDown
				}
				text := fmt.Sprintf("%.2f%%", coin.PercentChange1H)
				ct.SetTableColumnWidthFromString(header, text)
				ct.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin:  leftMargin,
						RightMargin: rightMargin,
						LeftAlign:   false,
						Color:       color1h,
						Text:        text,
					})
			case "24h_change":
				color24h := ct.colorscheme.TableColumnChange
				if coin.PercentChange24H > 0 {
					color24h = ct.colorscheme.TableColumnChangeUp
				}
				if coin.PercentChange24H < 0 {
					color24h = ct.colorscheme.TableColumnChangeDown
				}
				text := fmt.Sprintf("%.2f%%", coin.PercentChange24H)
				ct.SetTableColumnWidthFromString(header, text)
				ct.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin:  leftMargin,
						RightMargin: rightMargin,
						LeftAlign:   false,
						Color:       color24h,
						Text:        text,
					})
			case "7d_change":
				color7d := ct.colorscheme.TableColumnChange
				if coin.PercentChange7D > 0 {
					color7d = ct.colorscheme.TableColumnChangeUp
				}
				if coin.PercentChange7D < 0 {
					color7d = ct.colorscheme.TableColumnChangeDown
				}
				text := fmt.Sprintf("%.2f%%", coin.PercentChange7D)
				ct.SetTableColumnWidthFromString(header, text)
				ct.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin:  leftMargin,
						RightMargin: rightMargin,
						LeftAlign:   false,
						Color:       color7d,
						Text:        text,
					})
			case "30d_change":
				color30d := ct.colorscheme.TableColumnChange
				if coin.PercentChange30D > 0 {
					color30d = ct.colorscheme.TableColumnChangeUp
				}
				if coin.PercentChange30D < 0 {
					color30d = ct.colorscheme.TableColumnChangeDown
				}
				text := fmt.Sprintf("%.2f%%", coin.PercentChange30D)
				ct.SetTableColumnWidthFromString(header, text)
				ct.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin:  leftMargin,
						RightMargin: rightMargin,
						LeftAlign:   false,
						Color:       color30d,
						Text:        text,
					})
			case "market_cap":
				text := humanize.Commaf(coin.MarketCap)
				ct.SetTableColumnWidthFromString(header, text)
				ct.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin:  leftMargin,
						RightMargin: rightMargin,
						LeftAlign:   false,
						Color:       ct.colorscheme.TableRow,
						Text:        text,
					})
			case "total_supply":
				text := humanize.Commaf(coin.TotalSupply)
				ct.SetTableColumnWidthFromString(header, text)
				ct.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin:  leftMargin,
						RightMargin: rightMargin,
						LeftAlign:   false,
						Color:       ct.colorscheme.TableRow,
						Text:        text,
					})
			case "available_supply":
				text := humanize.Commaf(coin.AvailableSupply)
				ct.SetTableColumnWidthFromString(header, text)
				ct.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin:  leftMargin,
						RightMargin: rightMargin,
						LeftAlign:   false,
						Color:       ct.colorscheme.TableRow,
						Text:        text,
					})
			case "last_updated":
				unix, _ := strconv.ParseInt(coin.LastUpdated, 10, 64)
				lastUpdated := time.Unix(unix, 0).Format("15:04:05 Jan 02")
				ct.SetTableColumnWidthFromString(header, lastUpdated)
				ct.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin:  leftMargin,
						RightMargin: rightMargin,
						LeftAlign:   false,
						Color:       ct.colorscheme.TableRow,
						Text:        lastUpdated,
					})
			}
		}
		rows = append(rows, rowCells)
	}

	for _, row := range rows {
		for i, header := range headers {
			row[i].Width = ct.GetTableColumnWidth(header)
		}
		t.AddRowCells(row...)
	}

	return t
}
