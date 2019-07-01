package cointop

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
)

// TODO: break up into small functions

// layout sets initial layout
func (ct *Cointop) layout(g *gocui.Gui) error {
	maxX, maxY := ct.size()
	topOffset := 0

	headerHeight := 1
	marketbarHeight := 1
	chartHeight := 10
	statusbarHeight := 1

	if ct.State.onlyTable {
		ct.State.hideMarketbar = true
		ct.State.hideChart = true
		ct.State.hideStatusbar = true
	}

	if ct.State.hideMarketbar {
		marketbarHeight = 0
	}

	if ct.State.hideChart {
		chartHeight = 0
	}

	if ct.State.hideStatusbar {
		statusbarHeight = 0
	}

	if !ct.State.hideMarketbar {
		if v, err := g.SetView(ct.Views.Marketbar.Name, 0, topOffset, maxX, 2); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			ct.Views.Marketbar.Backing = v
			ct.Views.Marketbar.Backing.Frame = false
			ct.colorscheme.SetViewColor(ct.Views.Marketbar.Backing, "marketbar")
			go func() {
				ct.updateMarketbar()
				_, found := ct.cache.Get(ct.Views.Marketbar.Name)
				if found {
					ct.cache.Delete(ct.Views.Marketbar.Name)
					ct.updateMarketbar()
				}
			}()
		}
	} else {
		if ct.Views.Marketbar.Backing != nil {
			if err := g.DeleteView(ct.Views.Marketbar.Name); err != nil {
				return err
			}
			ct.Views.Marketbar.Backing = nil
		}
	}

	topOffset = topOffset + marketbarHeight

	if !ct.State.hideChart {
		if v, err := g.SetView(ct.Views.Chart.Name, 0, topOffset, maxX, topOffset+chartHeight+marketbarHeight); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			ct.Views.Chart.Backing = v
			ct.Views.Chart.Backing.Frame = false
			ct.colorscheme.SetViewColor(ct.Views.Chart.Backing, "chart")
			go func() {
				ct.updateChart()
				cachekey := strings.ToLower(fmt.Sprintf("%s_%s", "globaldata", strings.Replace(ct.State.selectedChartRange, " ", "", -1)))
				_, found := ct.cache.Get(cachekey)
				if found {
					ct.cache.Delete(cachekey)
					ct.updateChart()
				}
			}()
		}
	} else {
		if ct.Views.Chart.Backing != nil {
			if err := g.DeleteView(ct.Views.Chart.Name); err != nil {
				return err
			}
			ct.Views.Chart.Backing = nil
		}
	}

	topOffset = topOffset + chartHeight
	if v, err := g.SetView(ct.Views.Header.Name, 0, topOffset, ct.maxTableWidth, topOffset+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ct.Views.Header.Backing = v
		ct.Views.Header.Backing.Frame = false
		ct.colorscheme.SetViewColor(ct.Views.Header.Backing, "table_header")
		go ct.updateHeaders()
	}

	topOffset = topOffset + headerHeight
	if v, err := g.SetView(ct.Views.Table.Name, 0, topOffset, ct.maxTableWidth, maxY-statusbarHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ct.Views.Table.Backing = v
		ct.Views.Table.Backing.Frame = false
		ct.Views.Table.Backing.Highlight = true
		ct.colorscheme.SetViewActiveColor(ct.Views.Table.Backing, "table_row_active")
		_, found := ct.cache.Get("allCoinsSlugMap")
		if found {
			ct.cache.Delete("allCoinsSlugMap")
		}
		go func() {
			ct.updateCoins()
			ct.updateTable()
		}()
	}

	if !ct.State.hideStatusbar {
		if v, err := g.SetView(ct.Views.Statusbar.Name, 0, maxY-statusbarHeight-1, ct.maxTableWidth, maxY); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			ct.Views.Statusbar.Backing = v
			ct.Views.Statusbar.Backing.Frame = false
			ct.colorscheme.SetViewColor(ct.Views.Statusbar.Backing, "statusbar")
			go ct.updateStatusbar("")
		}
	} else {
		if ct.Views.Statusbar.Backing != nil {
			if err := g.DeleteView(ct.Views.Statusbar.Name); err != nil {
				return err
			}
			ct.Views.Statusbar.Backing = nil
		}
	}

	if v, err := g.SetView(ct.Views.SearchField.Name, 0, maxY-2, ct.maxTableWidth, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ct.Views.SearchField.Backing = v
		ct.Views.SearchField.Backing.Editable = true
		ct.Views.SearchField.Backing.Wrap = true
		ct.Views.SearchField.Backing.Frame = false
		ct.colorscheme.SetViewColor(ct.Views.SearchField.Backing, "searchbar")
	}

	if v, err := g.SetView(ct.Views.Help.Name, 1, 1, ct.maxTableWidth-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ct.Views.Help.Backing = v
		ct.Views.Help.Backing.Frame = false
		ct.colorscheme.SetViewColor(ct.Views.Help.Backing, "menu")
	}

	if v, err := g.SetView(ct.Views.PortfolioUpdateMenu.Name, 1, 1, ct.maxTableWidth-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ct.Views.PortfolioUpdateMenu.Backing = v
		ct.Views.PortfolioUpdateMenu.Backing.Frame = false
		ct.colorscheme.SetViewColor(ct.Views.PortfolioUpdateMenu.Backing, "menu")
	}

	if v, err := g.SetView(ct.Views.Input.Name, 3, 6, 30, 8); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ct.Views.Input.Backing = v
		ct.Views.Input.Backing.Frame = true
		ct.Views.Input.Backing.Editable = true
		ct.Views.Input.Backing.Wrap = true
		ct.colorscheme.SetViewColor(ct.Views.Input.Backing, "menu")
	}

	if v, err := g.SetView(ct.Views.ConvertMenu.Name, 1, 1, ct.maxTableWidth-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ct.Views.ConvertMenu.Backing = v
		ct.Views.ConvertMenu.Backing.Frame = false
		ct.colorscheme.SetViewColor(ct.Views.ConvertMenu.Backing, "menu")

		// run only once on init.
		// this bit of code should be at the bottom
		ct.g = g
		g.SetViewOnBottom(ct.Views.SearchField.Name)         // hide
		g.SetViewOnBottom(ct.Views.Help.Name)                // hide
		g.SetViewOnBottom(ct.Views.ConvertMenu.Name)         // hide
		g.SetViewOnBottom(ct.Views.PortfolioUpdateMenu.Name) // hide
		g.SetViewOnBottom(ct.Views.Input.Name)               // hide
		ct.setActiveView(ct.Views.Table.Name)
		ct.intervalFetchData()
	}

	return nil
}
