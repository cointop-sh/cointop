package cointop

import (
	"fmt"

	"github.com/miguelmota/cointop/pkg/gocui"
)

// layout sets initial layout
func (ct *Cointop) layout(g *gocui.Gui) error {
	maxX, maxY := ct.size()
	chartHeight := 10
	topOffset := 0

	if v, err := g.SetView(ct.marketbarviewname, 0, topOffset, maxX, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ct.marketbarview = v
		ct.marketbarview.Frame = false
		ct.marketbarview.BgColor = gocui.ColorBlack
		ct.marketbarview.FgColor = gocui.ColorWhite
		go func() {
			ct.updateMarketbar()
			_, found := ct.cache.Get(ct.marketbarviewname)
			if found {
				ct.cache.Delete(ct.marketbarviewname)
				ct.updateMarketbar()
			}
		}()
	}

	topOffset = topOffset + 1
	if v, err := g.SetView(ct.chartviewname, 0, topOffset, maxX, topOffset+chartHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ct.chartview = v
		ct.chartview.Frame = false
		go func() {
			ct.updateChart()
			_, found := ct.cache.Get("globaldata")
			if found {
				ct.cache.Delete("globaldata")
				ct.updateChart()
			}
		}()
	}

	topOffset = topOffset + chartHeight
	if v, err := g.SetView(ct.headerviewname, 0, topOffset, ct.maxtablewidth, topOffset+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ct.headersview = v
		ct.headersview.Frame = false
		ct.headersview.FgColor = gocui.ColorBlack
		ct.headersview.BgColor = gocui.ColorGreen
		go ct.updateHeaders()
	}

	topOffset = topOffset + 1
	if v, err := g.SetView(ct.tableviewname, 0, topOffset, ct.maxtablewidth, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ct.tableview = v
		ct.tableview.Frame = false
		ct.tableview.Highlight = true
		ct.tableview.SelBgColor = gocui.ColorCyan
		ct.tableview.SelFgColor = gocui.ColorBlack
		go func() {
			ct.updateCoins()
			ct.updateTable()
			_, found := ct.cache.Get("allcoinsmap")
			if found {
				ct.cache.Delete("allcoinsmap")
				ct.updateCoins()
				ct.updateTable()
			}
		}()
	}

	if v, err := g.SetView(ct.statusbarviewname, 0, maxY-2, ct.maxtablewidth, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ct.statusbarview = v
		ct.statusbarview.Frame = false
		ct.statusbarview.BgColor = gocui.ColorCyan
		ct.statusbarview.FgColor = gocui.ColorBlack
		go ct.updateStatusbar("")
	}

	if v, err := g.SetView(ct.searchfieldviewname, 0, maxY-2, ct.maxtablewidth, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ct.searchfield = v
		ct.searchfield.Editable = true
		ct.searchfield.Wrap = true
		ct.searchfield.Frame = false
		ct.searchfield.FgColor = gocui.ColorWhite
	}

	if v, err := g.SetView(ct.helpviewname, 1, 1, ct.maxtablewidth-2, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ct.helpview = v
		ct.helpview.Frame = false
		ct.helpview.BgColor = gocui.ColorBlack
		ct.helpview.FgColor = gocui.ColorWhite

		// run only once on init.
		// this bit of code should be at the bottom
		ct.g = g
		g.SetViewOnBottom(ct.searchfieldviewname) // hide
		g.SetViewOnBottom(ct.helpviewname)        // hide
		ct.setActiveView(ct.tableviewname)
		ct.intervalFetchData()
	}

	return nil
}

func (ct *Cointop) setActiveView(v string) error {
	ct.g.SetViewOnTop(v)
	ct.g.SetCurrentView(v)
	if v == ct.searchfieldviewname {
		ct.searchfield.Clear()
		ct.searchfield.SetCursor(1, 0)
		fmt.Fprintf(ct.searchfield, "%s", "/")
	} else if v == ct.tableviewname {
		ct.g.SetViewOnTop(ct.statusbarviewname)
	}
	return nil
}

func (ct *Cointop) setViewOnBottom(v string) error {
	_, err := ct.g.SetViewOnBottom(v)
	return err
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
