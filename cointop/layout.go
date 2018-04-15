package cointop

import (
	"fmt"

	"github.com/jroimartin/gocui"
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
	if v, err := g.SetView("header", 0, topOffset, ct.maxtablewidth, topOffset+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ct.headersview = v
		ct.headersview.Frame = false
		ct.headersview.FgColor = gocui.ColorBlack
		ct.headersview.BgColor = gocui.ColorGreen
		ct.updateHeaders()
	}

	topOffset = topOffset + 1
	if v, err := g.SetView("table", 0, topOffset, ct.maxtablewidth, maxY-1); err != nil {
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

	if v, err := g.SetView("statusbar", 0, maxY-2, ct.maxtablewidth, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ct.statusbarview = v
		ct.statusbarview.Frame = false
		ct.statusbarview.BgColor = gocui.ColorCyan
		ct.statusbarview.FgColor = gocui.ColorBlack
		ct.updateStatusbar("")
	}

	if v, err := g.SetView("searchfield", 0, maxY-2, ct.maxtablewidth, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ct.searchfield = v
		ct.searchfield.Editable = true
		ct.searchfield.Wrap = true
		ct.searchfield.Frame = false
		ct.searchfield.FgColor = gocui.ColorWhite

		// run only once on init
		ct.g = g
		g.SetViewOnBottom("searchfield")
		ct.setActiveView("table")
		ct.intervalFetchData()
	}
	return nil
}

func (ct *Cointop) setActiveView(v string) error {
	ct.g.SetViewOnTop(v)
	ct.g.SetCurrentView(v)
	if v == "searchfield" {
		ct.searchfield.Clear()
		ct.searchfield.SetCursor(1, 0)
		fmt.Fprintf(ct.searchfield, "%s", "/")
	} else if v == "table" {
		ct.g.SetViewOnTop("statusbar")
	}
	return nil
}

func (ct *Cointop) updateCoins() error {
	list := []*coin{}
	allcoinsmap, err := ct.api.GetAllCoinData()
	if err != nil {
		return err
	}

	if len(ct.allcoinsmap) == 0 {
		ct.allcoinsmap = map[string]coin{}
	}

	for k, v := range allcoinsmap {
		ct.allcoinsmap[k] = coin{
			ID:               v.ID,
			Name:             v.Name,
			Symbol:           v.Symbol,
			Rank:             v.Rank,
			PriceUSD:         v.PriceUSD,
			PriceBTC:         v.PriceBTC,
			USD24HVolume:     v.USD24HVolume,
			MarketCapUSD:     v.MarketCapUSD,
			AvailableSupply:  v.AvailableSupply,
			TotalSupply:      v.TotalSupply,
			PercentChange1H:  v.PercentChange1H,
			PercentChange24H: v.PercentChange24H,
			PercentChange7D:  v.PercentChange7D,
			LastUpdated:      v.LastUpdated,
		}
	}
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
