package cointop

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cointop-sh/cointop/pkg/humanize"
	"github.com/cointop-sh/cointop/pkg/table"
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
	"1y_change",
	"24h_volume",
	"market_cap",
	"available_supply",
	"total_supply",
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
	"available_supply",
	"total_supply",
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
	maxX := ct.Width()
	t := table.NewTable().SetWidth(maxX)
	var rows [][]*table.RowCell
	headers := ct.GetCoinsTableHeaders()
	if ct.IsFavoritesVisible() {
		headers = ct.GetFavoritesTableHeaders()
	}
	ct.ClearSyncMap(&ct.State.tableColumnWidths)
	ct.ClearSyncMap(&ct.State.tableColumnAlignLeft)
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
				star := " "
				rankcolor := ct.colorscheme.TableRow
				if coin.Favorite {
					star = ct.State.favoriteChar
					rankcolor = ct.colorscheme.TableRowFavorite
				}
				rank := fmt.Sprintf("%s%6v ", star, coin.Rank)
				ct.SetTableColumnWidth(header, 8)
				ct.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells, &table.RowCell{
					LeftMargin:  leftMargin,
					RightMargin: rightMargin,
					LeftAlign:   false,
					Color:       rankcolor,
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
				text := ct.FormatPrice(coin.Price)
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
				text := humanize.Monetaryf(coin.Volume24H, 0)
				if ct.State.scaleNumbers {
					volScale, volSuffix := humanize.Scale(coin.Volume24H)
					text = humanize.Numericf(volScale, 1) + volSuffix
				}
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
				text := fmt.Sprintf("%v%%", humanize.Numericf(coin.PercentChange1H, 2))
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
				text := fmt.Sprintf("%v%%", humanize.Numericf(coin.PercentChange24H, 2))
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
				text := fmt.Sprintf("%v%%", humanize.Numericf(coin.PercentChange7D, 2))
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
				text := fmt.Sprintf("%v%%", humanize.Numericf(coin.PercentChange30D, 2))
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
			case "1y_change":
				color1y := ct.colorscheme.TableColumnChange
				if coin.PercentChange1Y > 0 {
					color1y = ct.colorscheme.TableColumnChangeUp
				}
				if coin.PercentChange1Y < 0 {
					color1y = ct.colorscheme.TableColumnChangeDown
				}
				text := fmt.Sprintf("%v%%", humanize.Numericf(coin.PercentChange1Y, 2))
				ct.SetTableColumnWidthFromString(header, text)
				ct.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin:  leftMargin,
						RightMargin: rightMargin,
						LeftAlign:   false,
						Color:       color1y,
						Text:        text,
					})
			case "market_cap":
				text := humanize.Monetaryf(coin.MarketCap, 0)
				if ct.State.scaleNumbers {
					volScale, volSuffix := humanize.Scale(coin.MarketCap)
					text = humanize.Numericf(volScale, 1) + volSuffix
				}
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
				text := humanize.Numericf(coin.TotalSupply, 0)
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
				text := humanize.Numericf(coin.AvailableSupply, 0)
				if ct.State.scaleNumbers {
					volScale, volSuffix := humanize.Scale(coin.AvailableSupply)
					text = humanize.Numericf(volScale, 1) + volSuffix
				}
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

// TableCoinsLen returns the number of coins in coins table
func (ct *Cointop) TableCoinsLen() int {
	return len(ct.GetTableCoinsSlice())
}
