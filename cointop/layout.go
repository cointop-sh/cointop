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
	ct.rowChanged()
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
