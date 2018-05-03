package cointop

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/miguelmota/cointop/pkg/color"
	"github.com/miguelmota/cointop/pkg/humanize"
	"github.com/miguelmota/cointop/pkg/table"
)

func (ct *Cointop) refreshTable() error {
	maxX := ct.width()
	ct.table = table.New().SetWidth(maxX)
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
	ct.table.HideColumHeaders = true
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
			namecolor(fmt.Sprintf("%.22s", name)),
			color.White(fmt.Sprintf("%.6s", coin.Symbol)),
			colorprice(fmt.Sprintf("%12s", humanize.Commaf(coin.PriceUSD))),
			color.White(fmt.Sprintf("%17s", humanize.Commaf(coin.MarketCapUSD))),
			color.White(fmt.Sprintf("%15s", humanize.Commaf(coin.USD24HVolume))),
			color1h(fmt.Sprintf("%8.2f%%", coin.PercentChange1H)),
			color24h(fmt.Sprintf("%8.2f%%", coin.PercentChange24H)),
			color7d(fmt.Sprintf("%8.2f%%", coin.PercentChange7D)),
			color.White(fmt.Sprintf("%21s", humanize.Commaf(coin.TotalSupply))),
			color.White(fmt.Sprintf("%18s", humanize.Commaf(coin.AvailableSupply))),
			color.White(fmt.Sprintf("%18s", lastUpdated)),
			// add %percent of cap
		)
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

	for i := range ct.allcoinsmap {
		v := ct.allcoinsmap[i]
		if ct.favorites[v.Symbol] {
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
	return fmt.Sprintf("https://coinmarketcap.com/currencies/%s", slug)
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
