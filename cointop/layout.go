package cointop

import (
	"math"

	"github.com/jroimartin/gocui"
	apt "github.com/miguelmota/cointop/pkg/api/types"
	"github.com/miguelmota/cointop/pkg/pad"
	"github.com/miguelmota/cointop/pkg/table"
)

func (ct *Cointop) layout(g *gocui.Gui) error {
	maxX, maxY := ct.Size()
	chartHeight := 10
	topOffset := 0

	if v, err := g.SetView("market", 0, topOffset, maxX, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ct.marketview = v
		ct.marketview.Frame = false
		ct.marketview.BgColor = gocui.ColorBlack
		ct.marketview.FgColor = gocui.ColorWhite
		ct.updateMarketbar()
	}

	topOffset = topOffset + 1
	if v, err := g.SetView("chart", 0, topOffset, maxX, topOffset+chartHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ct.chartview = v
		ct.chartview.Frame = false
		ct.updateChart()
	}

	topOffset = topOffset + chartHeight
	if v, err := g.SetView("header", 0, topOffset, maxX, topOffset+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		t := table.New().SetWidth(maxX)
		headers := []string{
			pad.Right("[r]ank", 7, " "),
			pad.Right("[n]ame", 18, " "),
			pad.Right("[s]ymbol", 8, " "),
			pad.Left("[p]rice", 13, " "),
			pad.Left("[m]arket cap", 17, " "),
			pad.Left("24H [v]olume", 15, " "),
			pad.Left("[1]H%", 9, " "),
			pad.Left("[2]4H%", 9, " "),
			pad.Left("[7]D%", 9, " "),
			pad.Left("[t]otal supply", 20, " "),
			pad.Left("[a]vailable supply", 18, " "),
			pad.Left("[l]ast updated", 18, " "),
		}
		for _, h := range headers {
			t.AddCol(h)
		}

		t.Format().Fprint(v)
		ct.headersview = v
		ct.headersview.Frame = false
		ct.headersview.Highlight = true
		ct.headersview.SelBgColor = gocui.ColorGreen
		ct.headersview.SelFgColor = gocui.ColorBlack
	}

	topOffset = topOffset + 1
	if v, err := g.SetView("table", 0, topOffset, maxX, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ct.tableview = v
		ct.tableview.Frame = false
		ct.tableview.Highlight = true
		ct.tableview.SelBgColor = gocui.ColorCyan
		ct.tableview.SelFgColor = gocui.ColorBlack
		ct.updateCoins()
		ct.updateTable()
		ct.rowChanged()
	}

	if v, err := g.SetView("statusbar", 0, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ct.statusbarview = v
		ct.statusbarview.Frame = false
		ct.statusbarview.BgColor = gocui.ColorCyan
		ct.statusbarview.FgColor = gocui.ColorBlack
		ct.updateStatusbar("")
	}

	ct.intervalFetchData()
	return nil
}

func (ct *Cointop) updateCoins() error {
	list := []*apt.Coin{}
	allcoinsmap, err := ct.api.GetAllCoinData()
	if err != nil {
		return err
	}

	ct.allcoinsmap = allcoinsmap
	if len(ct.allcoins) == 0 {
		for i := range ct.allcoinsmap {
			coin := ct.allcoinsmap[i]
			list = append(list, &coin)
		}
		ct.allcoins = list
		ct.sort(ct.sortby, ct.sortdesc, ct.allcoins)
	} else {
		// update list in place without changing order
		for i := range ct.allcoinsmap {
			cm := ct.allcoinsmap[i]
			for k := range ct.allcoins {
				c := ct.allcoins[k]
				if c.ID == cm.ID {
					// TODO: improve this
					c.ID = cm.ID
					c.Name = cm.Name
					c.Symbol = cm.Symbol
					c.Rank = cm.Rank
					c.PriceUSD = cm.PriceUSD
					c.PriceBTC = cm.PriceBTC
					c.USD24HVolume = cm.USD24HVolume
					c.MarketCapUSD = cm.MarketCapUSD
					c.AvailableSupply = cm.AvailableSupply
					c.TotalSupply = cm.TotalSupply
					c.PercentChange1H = cm.PercentChange1H
					c.PercentChange24H = cm.PercentChange24H
					c.PercentChange7D = cm.PercentChange7D
					c.LastUpdated = cm.LastUpdated
				}
			}
		}
	}
	return nil
}

func (ct *Cointop) updateTable() error {
	start := ct.page * ct.perpage
	end := start + ct.perpage
	if end >= len(ct.allcoins)-1 {
		start = int(math.Floor(float64(start/100)) * 100)

		end = len(ct.allcoins) - 1
	}

	sliced := ct.allcoins[start:end]
	ct.coins = sliced
	ct.sort(ct.sortby, ct.sortdesc, ct.coins)
	ct.refreshTable()
	return nil
}

func (ct *Cointop) intervalFetchData() {
	go func() {
		for {
			select {
			case <-ct.forcerefresh:
				ct.refreshAll()
			case <-ct.refreshticker.C:
				ct.refreshAll()
			}
		}
	}()
}
