package cointop

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/miguelmota/cointop/cointop/common/humanize"
	"github.com/miguelmota/cointop/cointop/common/pad"
	"github.com/miguelmota/cointop/cointop/common/table"
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
		ct.table.AddCol("")

		total := ct.getPortfolioTotal()

		for _, coin := range ct.coins {
			unix, _ := strconv.ParseInt(coin.LastUpdated, 10, 64)
			lastUpdated := time.Unix(unix, 0).Format("15:04:05 Jan 02")
			colorbalance := ct.colorscheme.TableColumnPrice
			color24h := ct.colorscheme.TableColumnChange
			if coin.PercentChange24H > 0 {
				color24h = ct.colorscheme.TableColumnChangeUp
			}
			if coin.PercentChange24H < 0 {
				color24h = ct.colorscheme.TableColumnChangeDown
			}
			name := coin.Name
			dots := "..."
			star := " "
			rank := fmt.Sprintf("%s%v", star, ct.colorscheme.TableRow(fmt.Sprintf("%6v ", coin.Rank)))
			if len(name) > 20 {
				name = fmt.Sprintf("%s%s", name[0:18], dots)
			}

			percentHoldings := (coin.Balance / total) * 1e2
			if math.IsNaN(percentHoldings) {
				percentHoldings = 0
			}

			ct.table.AddRow(
				rank,
				ct.colorscheme.TableRow(pad.Right(fmt.Sprintf("%.22s", name), 21, " ")),
				ct.colorscheme.TableRow(pad.Right(fmt.Sprintf("%.6s", coin.Symbol), 5, " ")),
				ct.colorscheme.TableRow(fmt.Sprintf("%13s", humanize.Commaf(coin.Price))),
				ct.colorscheme.TableRow(fmt.Sprintf("%15s", strconv.FormatFloat(coin.Holdings, 'f', -1, 64))),
				colorbalance(fmt.Sprintf("%15s", humanize.Commaf(coin.Balance))),
				color24h(fmt.Sprintf("%8.2f%%", coin.PercentChange24H)),
				ct.colorscheme.TableRow(fmt.Sprintf("%10.2f%%", percentHoldings)),
				ct.colorscheme.TableRow(pad.Right(fmt.Sprintf("%17s", lastUpdated), 80, " ")),
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
			if coin == nil {
				continue
			}
			unix, _ := strconv.ParseInt(coin.LastUpdated, 10, 64)
			lastUpdated := time.Unix(unix, 0).Format("15:04:05 Jan 02")
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
			name := coin.Name
			dots := "..."
			star := ct.colorscheme.TableRow(" ")
			if coin.Favorite {
				star = ct.colorscheme.TableRowFavorite("*")
			}
			rank := fmt.Sprintf("%s%v", star, ct.colorscheme.TableRow(fmt.Sprintf("%6v ", coin.Rank)))
			if len(name) > 20 {
				name = fmt.Sprintf("%s%s", name[0:18], dots)
			}

			symbolpadding := 5
			// NOTE: this is to adjust padding by 1 because when all name rows are
			// yellow it messes the spacing (need to debug)
			if ct.filterByFavorites {
				symbolpadding = 6
			}
			ct.table.AddRow(
				rank,
				namecolor(pad.Right(fmt.Sprintf("%.22s", name), 21, " ")),
				ct.colorscheme.TableRow(pad.Right(fmt.Sprintf("%.6s", coin.Symbol), symbolpadding, " ")),
				ct.colorscheme.TableColumnPrice(fmt.Sprintf("%12s", humanize.Commaf(coin.Price))),
				ct.colorscheme.TableRow(fmt.Sprintf("%18s", humanize.Commaf(coin.MarketCap))),
				ct.colorscheme.TableRow(fmt.Sprintf("%15s", humanize.Commaf(coin.Volume24H))),
				color1h(fmt.Sprintf("%8.2f%%", coin.PercentChange1H)),
				color24h(fmt.Sprintf("%8.2f%%", coin.PercentChange24H)),
				color7d(fmt.Sprintf("%8.2f%%", coin.PercentChange7D)),
				ct.colorscheme.TableRow(fmt.Sprintf("%21s", humanize.Commaf(coin.TotalSupply))),
				ct.colorscheme.TableRow(fmt.Sprintf("%18s", humanize.Commaf(coin.AvailableSupply))),
				ct.colorscheme.TableRow(fmt.Sprintf("%18s", lastUpdated)),
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
		go ct.rowChanged()
		go ct.updateHeaders()
		go ct.updateMarketbar()
		go ct.updateChart()
	})

	return nil
}

func (ct *Cointop) updateTable() error {
	sliced := []*Coin{}

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
		go ct.refreshTable()
		return nil
	}

	if ct.portfoliovisible {
		sliced = ct.getPortfolioSlice()
		ct.coins = sliced
		go ct.refreshTable()
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

	ct.sort(ct.sortby, ct.sortdesc, ct.coins, true)
	go ct.refreshTable()
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

func (ct *Cointop) highlightedRowCoin() *Coin {
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

	return ct.api.CoinLink(coin.Name)
}

func (ct *Cointop) rowLinkShort() string {
	link := ct.rowLink()
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

		return fmt.Sprintf("http://%s/.../%s", host, path)
	}

	return ""
}

func (ct *Cointop) allCoins() []*Coin {
	if ct.filterByFavorites {
		var list []*Coin
		for i := range ct.allcoins {
			coin := ct.allcoins[i]
			if coin.Favorite {
				list = append(list, coin)
			}
		}
		return list
	}

	if ct.portfoliovisible {
		var list []*Coin
		for i := range ct.allcoins {
			coin := ct.allcoins[i]
			if ct.portfolioEntryExists(coin) {
				list = append(list, coin)
			}
		}
		return list
	}

	return ct.allcoins
}

func (ct *Cointop) coinBySymbol(symbol string) *Coin {
	for i := range ct.allcoins {
		coin := ct.allcoins[i]
		if coin.Symbol == symbol {
			return coin
		}
	}

	return nil
}
