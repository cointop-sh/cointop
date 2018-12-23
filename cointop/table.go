package cointop

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/miguelmota/cointop/pkg/color"
	"github.com/miguelmota/cointop/pkg/humanize"
	"github.com/miguelmota/cointop/pkg/pad"
	"github.com/miguelmota/cointop/pkg/table"
)

func (ct *Cointop) refreshTable() error {
	maxX := ct.width()
	ct.table = table.New().SetWidth(maxX)
	ct.table.HideColumHeaders = true

	if ct.portfoliovisible {
		ct.table.AddCol("")
		ct.table.AddCol("")
		ct.table.AddCol("")
		ct.table.AddCol("")
		ct.table.AddCol("")
		ct.table.AddCol("")
		ct.table.AddCol("")
		ct.table.AddCol("")
		for _, coin := range ct.coins {
			unix, _ := strconv.ParseInt(coin.LastUpdated, 10, 64)
			lastUpdated := time.Unix(unix, 0).Format("15:04:05 Jan 02")
			namecolor := color.White
			colorprice := color.White
			colorbalance := color.Cyan
			color24h := color.White
			if coin.PercentChange24H > 0 {
				color24h = color.Green
			}
			if coin.PercentChange24H < 0 {
				color24h = color.Red
			}
			name := coin.Name
			dots := "..."
			star := " "
			rank := fmt.Sprintf("%s%v", color.Yellow(star), color.White(fmt.Sprintf("%6v ", coin.Rank)))
			if len(name) > 20 {
				name = fmt.Sprintf("%s%s", name[0:18], dots)
			}

			ct.table.AddRow(
				rank,
				namecolor(pad.Right(fmt.Sprintf("%.22s", name), 21, " ")),
				color.White(pad.Right(fmt.Sprintf("%.6s", coin.Symbol), 5, " ")),
				colorprice(fmt.Sprintf("%13s", humanize.Commaf(coin.Price))),
				color.White(fmt.Sprintf("%15s", humanize.Commaf(coin.Holdings))),
				colorbalance(fmt.Sprintf("%15s", humanize.Commaf(coin.Balance))),
				color24h(fmt.Sprintf("%8.2f%%", coin.PercentChange24H)),
				color.White(pad.Right(fmt.Sprintf("%17s", lastUpdated), 80, " ")),
			)
		}
	} else {
		ct.table.AddCol("")
		ct.table.AddCol("")
		ct.table.AddCol("")
		ct.table.AddCol("")
		ct.table.AddCol("")
		ct.table.AddCol("")
		ct.table.AddCol("")
		ct.table.AddCol("")
		ct.table.AddCol("")
		ct.table.AddCol("")
		ct.table.AddCol("")
		ct.table.AddCol("")
		for _, coin := range ct.coins {
			unix, _ := strconv.ParseInt(coin.LastUpdated, 10, 64)
			lastUpdated := time.Unix(unix, 0).Format("15:04:05 Jan 02")
			namecolor := color.White
			colorprice := color.Cyan
			color1h := color.White
			color24h := color.White
			color7d := color.White
			if coin.Favorite {
				namecolor = color.Yellow
			}
			if coin.PercentChange1H > 0 {
				color1h = color.Green
			}
			if coin.PercentChange1H < 0 {
				color1h = color.Red
			}
			if coin.PercentChange24H > 0 {
				color24h = color.Green
			}
			if coin.PercentChange24H < 0 {
				color24h = color.Red
			}
			if coin.PercentChange7D > 0 {
				color7d = color.Green
			}
			if coin.PercentChange7D < 0 {
				color7d = color.Red
			}
			name := coin.Name
			dots := "..."
			star := " "
			if coin.Favorite {
				star = "*"
			}
			rank := fmt.Sprintf("%s%v", color.Yellow(star), color.White(fmt.Sprintf("%6v ", coin.Rank)))
			if len(name) > 20 {
				name = fmt.Sprintf("%s%s", name[0:18], dots)
			}
			ct.table.AddRow(
				rank,
				namecolor(pad.Right(fmt.Sprintf("%.22s", name), 21, " ")),
				color.White(pad.Right(fmt.Sprintf("%.6s", coin.Symbol), 5, " ")),
				colorprice(fmt.Sprintf("%12s", humanize.Commaf(coin.Price))),
				color.White(fmt.Sprintf("%17s", humanize.Commaf(coin.MarketCap))),
				color.White(fmt.Sprintf("%15s", humanize.Commaf(coin.Volume24H))),
				color1h(fmt.Sprintf("%8.2f%%", coin.PercentChange1H)),
				color24h(fmt.Sprintf("%8.2f%%", coin.PercentChange24H)),
				color7d(fmt.Sprintf("%8.2f%%", coin.PercentChange7D)),
				color.White(fmt.Sprintf("%21s", humanize.Commaf(coin.TotalSupply))),
				color.White(fmt.Sprintf("%18s", humanize.Commaf(coin.AvailableSupply))),
				color.White(fmt.Sprintf("%18s", lastUpdated)),
				// TODO: add %percent of cap
			)
		}
	}

	// highlight last row if current row is out of bounds (can happen when switching views)
	currentrow := ct.highlightedRowIndex()
	if len(ct.coins) > currentrow {
		ct.highlightRow(currentrow)
	}

	ct.update(func() {
		ct.tableview.Clear()
		ct.table.Format().Fprint(ct.tableview)
		ct.rowChanged()
	})

	return nil
}

func (ct *Cointop) updateTable() error {
	sliced := []*coin{}

	for i := range ct.allcoinsslugmap {
		v := ct.allcoinsslugmap[i]
		if ct.favorites[v.Name] {
			v.Favorite = true
		}
	}

	if ct.filterByFavorites {
		for i := range ct.allcoins {
			coin := ct.allcoins[i]
			if coin.Favorite {
				sliced = append(sliced, coin)
			}
		}
		ct.coins = sliced
		ct.sort(ct.sortby, ct.sortdesc, ct.coins)
		ct.refreshTable()
		return nil
	}

	if ct.portfoliovisible {
		for i := range ct.allcoins {
			if len(ct.portfolio.Entries) == 0 {
				break
			}
			coin := ct.allcoins[i]
			var p *portfolioEntry
			var ok bool
			if p, ok = ct.portfolio.Entries[strings.ToLower(coin.Name)]; !ok {
				// NOTE: if not found then try the symbol
				if p, ok = ct.portfolio.Entries[strings.ToLower(coin.Symbol)]; !ok {
					continue
				}
			}
			holdingsstr := fmt.Sprintf("%.2f", p.Holdings)
			if ct.currencyconversion == "ETH" || ct.currencyconversion == "BTC" {
				holdingsstr = fmt.Sprintf("%.5f", p.Holdings)
			}
			holdings, _ := strconv.ParseFloat(holdingsstr, 64)
			coin.Holdings = holdings

			balance := coin.Price * p.Holdings
			balancestr := fmt.Sprintf("%.2f", balance)
			if ct.currencyconversion == "ETH" || ct.currencyconversion == "BTC" {
				balancestr = fmt.Sprintf("%.5f", balance)
			}
			balance, _ = strconv.ParseFloat(balancestr, 64)
			coin.Balance = balance
			sliced = append(sliced, coin)
		}

		sort.Slice(sliced, func(i, j int) bool {
			return sliced[i].Balance > sliced[j].Balance
		})

		for i, coin := range sliced {
			coin.Rank = i + 1
		}

		ct.coins = sliced
		ct.sort(ct.sortby, ct.sortdesc, ct.coins)
		ct.refreshTable()
		return nil
	}

	start := ct.page * ct.perpage
	end := start + ct.perpage
	allcoins := ct.allCoins()
	size := len(allcoins)
	if start < 0 {
		start = 0
	}
	if end >= size-1 {
		start = int(math.Floor(float64(start/100)) * 100)
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
		sliced = allcoins[start:end]
	}
	ct.coins = sliced
	ct.sort(ct.sortby, ct.sortdesc, ct.coins)
	ct.refreshTable()
	return nil
}

func (ct *Cointop) highlightedRowIndex() int {
	_, y := ct.tableview.Origin()
	_, cy := ct.tableview.Cursor()
	idx := y + cy
	if idx < 0 {
		idx = 0
	}
	if idx >= len(ct.coins) {
		idx = len(ct.coins) - 1
	}
	return idx
}

func (ct *Cointop) highlightedRowCoin() *coin {
	idx := ct.highlightedRowIndex()
	if len(ct.coins) == 0 {
		return nil
	}
	return ct.coins[idx]
}

func (ct *Cointop) rowLink() string {
	coin := ct.highlightedRowCoin()
	if coin == nil {
		return ""
	}
	slug := strings.ToLower(strings.Replace(coin.Name, " ", "-", -1))
	// TODO: dynamic
	return fmt.Sprintf("https://coinmarketcap.com/currencies/%s", slug)
}

func (ct *Cointop) rowLinkShort() string {
	coin := ct.highlightedRowCoin()
	if coin == nil {
		return ""
	}
	// TODO: dynamic
	slug := strings.ToLower(strings.Replace(coin.Name, " ", "-", -1))
	return fmt.Sprintf("http://coinmarketcap.com/.../%s", slug)
}

func (ct *Cointop) allCoins() []*coin {
	if ct.filterByFavorites {
		var list []*coin
		for i := range ct.allcoins {
			coin := ct.allcoins[i]
			if coin.Favorite {
				list = append(list, coin)
			}
		}
		return list
	}

	return ct.allcoins
}

func (ct *Cointop) coinBySymbol(symbol string) *coin {
	for i := range ct.allcoins {
		coin := ct.allcoins[i]
		if coin.Symbol == symbol {
			return coin
		}
	}

	return nil
}
