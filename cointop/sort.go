package cointop

import (
	"github.com/bradfitz/slice"
	"github.com/jroimartin/gocui"
	apt "github.com/miguelmota/cointop/pkg/api/types"
)

func (ct *Cointop) sort(sortby string, desc bool, list []*apt.Coin) {
	ct.sortby = sortby
	ct.sortdesc = desc
	slice.Sort(list[:], func(i, j int) bool {
		if ct.sortdesc {
			i, j = j, i
		}
		a := list[i]
		b := list[j]
		switch sortby {
		case "rank":
			return a.Rank < b.Rank
		case "name":
			return a.Name < b.Name
		case "symbol":
			return a.Symbol < b.Symbol
		case "price":
			return a.PriceUSD < b.PriceUSD
		case "marketcap":
			return a.MarketCapUSD < b.MarketCapUSD
		case "24hvolume":
			return a.USD24HVolume < b.USD24HVolume
		case "1hchange":
			return a.PercentChange1H < b.PercentChange1H
		case "24hchange":
			return a.PercentChange24H < b.PercentChange24H
		case "7dchange":
			return a.PercentChange7D < b.PercentChange7D
		case "totalsupply":
			return a.TotalSupply < b.TotalSupply
		case "availablesupply":
			return a.AvailableSupply < b.AvailableSupply
		case "lastupdated":
			return a.LastUpdated < b.LastUpdated
		default:
			return a.Rank < b.Rank
		}
	})
}

func (ct *Cointop) sortfn(sortby string, desc bool) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if ct.sortby == sortby {
			desc = !desc
		}

		ct.sort(sortby, desc, ct.coins)
		ct.g.Update(func(g *gocui.Gui) error {
			ct.tableview.Clear()
			ct.updateTable()
			return nil
		})
		/*
			g.Update(func(g *gocui.Gui) error {
				ct.chartview.Clear()
				maxX, _ := g.Size()
				_, cy := ct.chartview.Cursor()
				coin := "ethereum"
				ct.chartPoints(maxX, coin)
				ct.updateChart()
				fmt.Fprint(ct.chartview, cy)
				return nil
			})
		*/

		ct.rowChanged()
		return nil
	}
}
