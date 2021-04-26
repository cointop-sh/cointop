package cointop

import (
	"fmt"
	"strings"
)

// TODO: break up into small functions

var lastWidth int

// layout sets initial layout
func (ct *Cointop) layout() error {
	ct.debuglog("layout()")
	maxY := ct.height()
	maxX := ct.width()

	topOffset := 0
	headerHeight := 1
	marketbarHeight := ct.State.marketBarHeight
	chartHeight := ct.State.lastChartHeight
	statusbarHeight := 1

	if ct.State.onlyTable {
		ct.State.hideMarketbar = true
		ct.State.hideChart = true
		ct.State.hideTable = false
		ct.State.hideStatusbar = true
		ct.State.onlyChart = false
		marketbarHeight = 0
	} else if ct.State.onlyChart {
		ct.State.hideMarketbar = true
		ct.State.hideChart = false
		ct.State.hideTable = true
		ct.State.hideStatusbar = true
		ct.State.onlyTable = false
		marketbarHeight = 0
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

	if ct.State.hideMarketbar {
		if ct.Views.Marketbar.Backing() != nil {
			if err := ct.g.DeleteView(ct.Views.Marketbar.Name()); err != nil {
				return err
			}
			ct.Views.Marketbar.SetBacking(nil)
		}
	} else {
		if err := ct.ui.SetView(ct.Views.Marketbar, 0, topOffset-1, maxX, marketbarHeight+1); err != nil {
			ct.Views.Marketbar.SetFrame(false)
			ct.Views.Marketbar.SetFgColor(ct.colorscheme.gocuiFgColor(ct.Views.Marketbar.Name()))
			ct.Views.Marketbar.SetBgColor(ct.colorscheme.gocuiBgColor(ct.Views.Marketbar.Name()))
			go func() {
				ct.UpdateMarketbar()
				_, found := ct.cache.Get(ct.Views.Marketbar.Name())
				if found {
					ct.cache.Delete(ct.Views.Marketbar.Name())
					ct.UpdateMarketbar()
				}
			}()
		}
	}

	topOffset = topOffset + marketbarHeight

	if ct.State.hideChart {
		if ct.Views.Chart.Backing() != nil {
			if err := ct.g.DeleteView(ct.Views.Chart.Name()); err != nil {
				return err
			}
			ct.Views.Chart.SetBacking(nil)
		}
	} else {
		chartTopOffset := topOffset - 1
		if ct.State.hideStatusbar {
			chartTopOffset = topOffset
		}
		if ct.State.onlyChart {
			chartHeight = maxY - topOffset
		}
		ct.State.chartHeight = chartHeight
		if err := ct.ui.SetView(ct.Views.Chart, 0, chartTopOffset, maxX, topOffset+chartHeight); err != nil {
			ct.Views.Chart.Clear()
			ct.Views.Chart.SetFrame(false)
			ct.Views.Chart.SetFgColor(ct.colorscheme.gocuiFgColor(ct.Views.Chart.Name()))
			ct.Views.Chart.SetBgColor(ct.colorscheme.gocuiBgColor(ct.Views.Chart.Name()))
			go func() {
				ct.UpdateChart()
				cachekey := strings.ToLower(fmt.Sprintf("%s_%s", "globaldata", strings.Replace(ct.State.selectedChartRange, " ", "", -1)))
				_, found := ct.cache.Get(cachekey)
				if found {
					ct.cache.Delete(cachekey)
					ct.UpdateChart()
				}
			}()
		}
	}

	if ct.State.hideTable {
		if ct.Views.TableHeader.Backing() != nil {
			if err := ct.g.DeleteView(ct.Views.TableHeader.Name()); err != nil {
				return err
			}
			ct.Views.TableHeader.SetBacking(nil)
		}
		if ct.Views.Table.Backing() != nil {
			if err := ct.g.DeleteView(ct.Views.Table.Name()); err != nil {
				return err
			}
			ct.Views.Table.SetBacking(nil)
		}
	} else {
		tableOffsetX := ct.State.tableOffsetX
		topOffset = topOffset + chartHeight
		if err := ct.ui.SetView(ct.Views.TableHeader, tableOffsetX, topOffset-1, maxX, topOffset+1); err != nil {
			ct.Views.TableHeader.SetFrame(false)
			ct.Views.TableHeader.SetFgColor(ct.colorscheme.gocuiFgColor(ct.Views.TableHeader.Name()))
			ct.Views.TableHeader.SetBgColor(ct.colorscheme.gocuiBgColor(ct.Views.TableHeader.Name()))
			go ct.UpdateTableHeader()
		}

		topOffset = topOffset + headerHeight
		if err := ct.ui.SetView(ct.Views.Table, tableOffsetX, topOffset-1, maxX, maxY-statusbarHeight); err != nil {
			ct.Views.Table.SetFrame(false)
			ct.Views.Table.SetHighlight(true)
			ct.Views.Table.SetSelFgColor(ct.colorscheme.gocuiFgColor("table_row_active"))
			ct.Views.Table.SetSelBgColor(ct.colorscheme.gocuiBgColor("table_row_active"))
			_, found := ct.cache.Get("allCoinsSlugMap")
			if found {
				ct.cache.Delete("allCoinsSlugMap")
			}
			go func() {
				ct.UpdateCoins()
				ct.UpdateTable()
			}()
		}
	}

	if !ct.State.hideStatusbar {
		if err := ct.ui.SetView(ct.Views.Statusbar, 0, maxY-statusbarHeight-1, maxX, maxY); err != nil {
			ct.Views.Statusbar.SetFrame(false)
			ct.Views.Statusbar.SetFgColor(ct.colorscheme.gocuiFgColor(ct.Views.Statusbar.Name()))
			ct.Views.Statusbar.SetBgColor(ct.colorscheme.gocuiBgColor(ct.Views.Statusbar.Name()))
			go ct.UpdateStatusbar("")
		}
	} else {
		if ct.Views.Statusbar.Backing() != nil {
			if err := ct.g.DeleteView(ct.Views.Statusbar.Name()); err != nil {
				return err
			}
			ct.Views.Statusbar.SetBacking(nil)
		}
	}

	if err := ct.ui.SetView(ct.Views.SearchField, 0, maxY-2, maxX, maxY); err != nil {
		ct.Views.SearchField.SetEditable(true)
		ct.Views.SearchField.SetWrap(true)
		ct.Views.SearchField.SetFrame(false)
		ct.Views.SearchField.SetFgColor(ct.colorscheme.gocuiFgColor("searchbar"))
		ct.Views.SearchField.SetBgColor(ct.colorscheme.gocuiBgColor("searchbar"))
	}

	if err := ct.ui.SetView(ct.Views.Menu, 1, 1, maxX-1, maxY-1); err != nil {
		ct.Views.Menu.SetFrame(false)
		ct.Views.Menu.SetFgColor(ct.colorscheme.gocuiFgColor("menu"))
		ct.Views.Menu.SetBgColor(ct.colorscheme.gocuiBgColor("menu"))
	}

	if err := ct.ui.SetView(ct.Views.Input, 3, 6, 30, 8); err != nil {
		ct.Views.Input.SetFrame(true)
		ct.Views.Input.SetEditable(true)
		ct.Views.Input.SetWrap(true)
		ct.Views.Input.SetFgColor(ct.colorscheme.gocuiFgColor("menu"))
		ct.Views.Input.SetBgColor(ct.colorscheme.gocuiBgColor("menu"))

		// run only once on init.
		// this bit of code should be at the bottom
		ct.ui.SetViewOnBottom(ct.Views.SearchField) // hide
		ct.ui.SetViewOnBottom(ct.Views.Input)       // hide
		ct.ui.SetViewOnBottom(ct.Views.Menu)        // hide
		ct.SetActiveView(ct.Views.Table.Name())
		ct.intervalFetchData()
	}

	if lastWidth != maxX {
		lastWidth = maxX
		ct.Refresh()
	}

	return nil
}
