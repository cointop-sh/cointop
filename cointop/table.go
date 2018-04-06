package cointop

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
	apt "github.com/miguelmota/cointop/pkg/api/types"
	"github.com/miguelmota/cointop/pkg/color"
	"github.com/miguelmota/cointop/pkg/humanize"
	"github.com/miguelmota/cointop/pkg/open"
	"github.com/miguelmota/cointop/pkg/table"
)

func (ct *Cointop) refreshTable() error {
	maxX := ct.Width()
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
		colorprice := color.Cyan
		color1h := color.White
		color24h := color.White
		color7d := color.White
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
		lastchar := len(name)
		if lastchar > 20 {
			lastchar = 20
			name = fmt.Sprintf("%s...", name[0:18])
		}
		ct.table.AddRow(
			color.White(fmt.Sprintf("%7v ", coin.Rank)),
			color.White(fmt.Sprintf("%.22s", name)),
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

	ct.Update(func() {
		ct.tableview.Clear()
		ct.table.Format().Fprint(ct.tableview)
	})

	return nil
}

func (ct *Cointop) selectedRowIndex() int {
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

func (ct *Cointop) selectedCoin() *apt.Coin {
	idx := ct.selectedRowIndex()
	return ct.coins[idx]
}

func (ct *Cointop) rowLink() string {
	slug := strings.ToLower(strings.Replace(ct.selectedCoin().Name, " ", "-", -1))
	return fmt.Sprintf("https://coinmarketcap.com/currencies/%s", slug)
}

func (ct *Cointop) openLink(g *gocui.Gui, v *gocui.View) error {
	open.URL(ct.rowLink())
	return nil
}
