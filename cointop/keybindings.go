package cointop

import (
	"strings"

	"github.com/miguelmota/gocui"
)

// ParseKeys returns string keyboard key as gocui key type
func (ct *Cointop) ParseKeys(s string) (interface{}, gocui.Modifier) {
	var key interface{}
	mod := gocui.ModNone
	split := strings.Split(s, "+")
	if len(split) > 1 {
		m := strings.ToLower(split[0])
		k := strings.ToLower(split[1])
		if m == "alt" {
			mod = gocui.ModAlt
			s = k
		} else if m == "ctrl" {
			switch k {
			case "0":
				key = '0'
			case "1":
				key = '1'
			case "2":
				key = gocui.KeyCtrl2
			case "3":
				key = gocui.KeyCtrl3
			case "4":
				key = gocui.KeyCtrl4
			case "5":
				key = gocui.KeyCtrl5
			case "6":
				key = gocui.KeyCtrl6
			case "7":
				key = gocui.KeyCtrl7
			case "8":
				key = gocui.KeyCtrl8
			case "9":
				key = '9'
			case "a":
				key = gocui.KeyCtrlA
			case "b":
				key = gocui.KeyCtrlB
			case "c":
				key = gocui.KeyCtrlC
			case "d":
				key = gocui.KeyCtrlD
			case "e":
				key = gocui.KeyCtrlE
			case "f":
				key = gocui.KeyCtrlF
			case "g":
				key = gocui.KeyCtrlG
			case "h":
				key = gocui.KeyCtrlH
			case "i":
				key = gocui.KeyCtrlI
			case "j":
				key = gocui.KeyCtrlJ
			case "k":
				key = gocui.KeyCtrlK
			case "l":
				key = gocui.KeyCtrlL
			case "m":
				key = gocui.KeyCtrlL
			case "n":
				key = gocui.KeyCtrlN
			case "o":
				key = gocui.KeyCtrlO
			case "p":
				key = gocui.KeyCtrlP
			case "q":
				key = gocui.KeyCtrlQ
			case "r":
				key = gocui.KeyCtrlR
			case "s":
				key = gocui.KeyCtrlS
			case "t":
				key = gocui.KeyCtrlT
			case "u":
				key = gocui.KeyCtrlU
			case "v":
				key = gocui.KeyCtrlV
			case "w":
				key = gocui.KeyCtrlW
			case "x":
				key = gocui.KeyCtrlX
			case "y":
				key = gocui.KeyCtrlY
			case "z":
				key = gocui.KeyCtrlZ
			case "~":
				key = gocui.KeyCtrlTilde
			case "[":
				fallthrough
			case "lsqrbracket":
				fallthrough
			case "leftsqrbracket":
				fallthrough
			case "leftsquarebracket":
				key = gocui.KeyCtrlLsqBracket
			case "]":
				fallthrough
			case "rsqrbracket":
				fallthrough
			case "rightsqrbracket":
				fallthrough
			case "rightsquarebracket":
				key = gocui.KeyCtrlRsqBracket
			case "space":
				key = gocui.KeyCtrlSpace
			case "backslash":
				key = gocui.KeyCtrlBackslash
			case "underscore":
				key = gocui.KeyCtrlUnderscore
			case "\\\\":
				key = '\\'
			}
			return key, mod
		}
	}

	if len(s) == 1 {
		r := []rune(s)
		key = r[0]
		return key, mod
	}

	s = strings.ToLower(s)
	switch s {
	case "arrowup":
		fallthrough
	case "uparrow":
		fallthrough
	case "up":
		key = gocui.KeyArrowUp
	case "arrowdown":
		fallthrough
	case "downarrow":
		fallthrough
	case "down":
		key = gocui.KeyArrowDown
	case "arrowleft":
		fallthrough
	case "leftarrow":
		fallthrough
	case "left":
		key = gocui.KeyArrowLeft
	case "arrowright":
		fallthrough
	case "rightarrow":
		fallthrough
	case "right":
		key = gocui.KeyArrowRight
	case "enter":
		fallthrough
	case "return":
		key = gocui.KeyEnter
	case "space":
		fallthrough
	case "spacebar":
		key = gocui.KeySpace
	case "esc":
		fallthrough
	case "escape":
		key = gocui.KeyEsc
	case "f1":
		key = gocui.KeyF1
	case "f2":
		key = gocui.KeyF2
	case "f3":
		key = gocui.KeyF3
	case "f4":
		key = gocui.KeyF4
	case "f5":
		key = gocui.KeyF5
	case "f6":
		key = gocui.KeyF6
	case "f7":
		key = gocui.KeyF7
	case "f8":
		key = gocui.KeyF8
	case "f9":
		key = gocui.KeyF9
	case "tab":
		key = gocui.KeyTab
	case "pageup":
		fallthrough
	case "pgup":
		key = gocui.KeyPgup
	case "pagedown":
		fallthrough
	case "pgdown":
		fallthrough
	case "pgdn":
		key = gocui.KeyPgdn
	case "home":
		key = gocui.KeyHome
	case "end":
		key = gocui.KeyEnd
	case "\\\\":
		key = '\\'
	}

	return key, mod
}

// Keybindings sets keyboard shortcut key bindings
func (ct *Cointop) Keybindings(g *gocui.Gui) error {
	for k, v := range ct.State.shortcutKeys {
		if k == "" {
			continue
		}
		v = strings.TrimSpace(strings.ToLower(v))
		var fn func(g *gocui.Gui, v *gocui.View) error
		key, mod := ct.ParseKeys(k)
		view := "table"
		switch v {
		case "move_up":
			fn = ct.Keyfn(ct.CursorUp)
		case "move_down":
			fn = ct.Keyfn(ct.CursorDown)
		case "previous_page":
			fn = ct.handleHkey(key)
		case "next_page":
			fn = ct.Keyfn(ct.NextPage)
		case "page_down":
			fn = ct.Keyfn(ct.PageDown)
		case "page_up":
			fn = ct.Keyfn(ct.PageUp)
		case "sort_column_symbol":
			fn = ct.Sortfn("symbol", false)
		case "move_to_page_first_row":
			fn = ct.Keyfn(ct.NavigateFirstLine)
		case "move_to_page_last_row":
			fn = ct.Keyfn(ct.NavigateLastLine)
		case "open_link":
			fn = ct.Keyfn(ct.OpenLink)
		case "refresh":
			fn = ct.Keyfn(ct.Refresh)
		case "sort_column_asc":
			fn = ct.Keyfn(ct.SortAsc)
		case "sort_column_desc":
			fn = ct.Keyfn(ct.SortDesc)
		case "sort_left_column":
			fn = ct.Keyfn(ct.SortPrevCol)
		case "sort_right_column":
			fn = ct.Keyfn(ct.SortNextCol)
		case "help":
			fallthrough
		case "toggle_show_help":
			fn = ct.Keyfn(ct.ToggleHelp)
			view = ""
		case "show_help":
			fn = ct.Keyfn(ct.ShowHelp)
			view = ""
		case "hide_help":
			fn = ct.Keyfn(ct.HideHelp)
			view = "help"
		case "first_page":
			fn = ct.Keyfn(ct.FirstPage)
		case "sort_column_1h_change":
			fn = ct.Sortfn("1hchange", true)
		case "sort_column_24h_change":
			fn = ct.Sortfn("24hchange", true)
		case "sort_column_7d_change":
			fn = ct.Sortfn("7dchange", true)
		case "sort_column_available_supply":
			fn = ct.Sortfn("availablesupply", true)
		case "toggle_row_chart":
			fn = ct.Keyfn(ct.ToggleCoinChart)
		case "move_to_page_visible_first_row":
			fn = ct.Keyfn(ct.NavigatePageFirstLine)
		case "move_to_page_visible_last_row":
			fn = ct.Keyfn(ct.navigatePageLastLine)
		case "sort_column_market_cap":
			fn = ct.Sortfn("marketcap", true)
		case "move_to_page_visible_middle_row":
			fn = ct.Keyfn(ct.NavigatePageMiddleLine)
		case "scroll_left":
			fn = ct.Keyfn(ct.TableScrollLeft)
		case "scroll_right":
			fn = ct.Keyfn(ct.TableScrollRight)
		case "sort_column_name":
			fn = ct.Sortfn("name", false)
		case "sort_column_price":
			fn = ct.Sortfn("price", true)
		case "sort_column_rank":
			fn = ct.Sortfn("rank", false)
		case "sort_column_total_supply":
			fn = ct.Sortfn("totalsupply", true)
		case "sort_column_last_updated":
			fn = ct.Sortfn("lastupdated", true)
		case "sort_column_24h_volume":
			fn = ct.Sortfn("24hvolume", true)
		case "sort_column_balance":
			fn = ct.Sortfn("balance", true)
		case "sort_column_holdings":
			fn = ct.Sortfn("holdings", true)
		case "sort_column_percent_holdings":
			fn = ct.Sortfn("percentholdings", true)
		case "last_page":
			fn = ct.Keyfn(ct.LastPage)
		case "open_search":
			fn = ct.Keyfn(ct.openSearch)
			view = ""
		case "toggle_favorite":
			fn = ct.Keyfn(ct.ToggleFavorite)
		case "toggle_favorites":
			fn = ct.Keyfn(ct.ToggleFavorites)
		case "toggle_show_favorites":
			fn = ct.Keyfn(ct.ToggleShowFavorites)
		case "save":
			fn = ct.Keyfn(ct.Save)
		case "quit":
			fn = ct.Keyfn(ct.Quit)
			view = ""
		case "quit_view":
			fn = ct.Keyfn(ct.QuitView)
		case "next_chart_range":
			fn = ct.Keyfn(ct.NextChartRange)
		case "previous_chart_range":
			fn = ct.Keyfn(ct.PrevChartRange)
		case "first_chart_range":
			fn = ct.Keyfn(ct.FirstChartRange)
		case "last_chart_range":
			fn = ct.Keyfn(ct.LastChartRange)
		case "toggle_show_currency_convert_menu":
			fn = ct.Keyfn(ct.ToggleConvertMenu)
		case "show_currency_convert_menu":
			fn = ct.Keyfn(ct.ShowConvertMenu)
		case "hide_currency_convert_menu":
			fn = ct.Keyfn(ct.HideConvertMenu)
			view = "convertmenu"
		case "toggle_portfolio":
			fn = ct.Keyfn(ct.TogglePortfolio)
		case "toggle_show_portfolio":
			fn = ct.Keyfn(ct.ToggleShowPortfolio)
		case "show_portfolio_edit_menu":
			fn = ct.Keyfn(ct.TogglePortfolioUpdateMenu)
		case "toggle_table_fullscreen":
			fn = ct.Keyfn(ct.ToggleTableFullscreen)
			view = ""
		case "enlarge_chart":
			fn = ct.Keyfn(ct.EnlargeChart)
		case "shorten_chart":
			fn = ct.Keyfn(ct.ShortenChart)
		case "move_down_or_next_page":
			fn = ct.Keyfn(ct.CursorDownOrNextPage)
		case "move_up_or_previous_page":
			fn = ct.Keyfn(ct.CursorUpOrPreviousPage)
		default:
			fn = ct.Keyfn(ct.Noop)
		}

		ct.SetKeybindingMod(key, mod, fn, view)
	}

	// keys to force quit
	ct.SetKeybindingMod(gocui.KeyCtrlC, gocui.ModNone, ct.Keyfn(ct.Quit), "")
	ct.SetKeybindingMod(gocui.KeyCtrlZ, gocui.ModNone, ct.Keyfn(ct.Quit), "")

	// searchfield keys
	ct.SetKeybindingMod(gocui.KeyEnter, gocui.ModNone, ct.Keyfn(ct.DoSearch), ct.Views.SearchField.Name())
	ct.SetKeybindingMod(gocui.KeyEsc, gocui.ModNone, ct.Keyfn(ct.CancelSearch), ct.Views.SearchField.Name())

	// keys to quit help when open
	ct.SetKeybindingMod(gocui.KeyEsc, gocui.ModNone, ct.Keyfn(ct.HideHelp), ct.Views.Help.Name())
	ct.SetKeybindingMod('q', gocui.ModNone, ct.Keyfn(ct.HideHelp), ct.Views.Help.Name())

	// keys to quit portfolio update menu when open
	ct.SetKeybindingMod(gocui.KeyEsc, gocui.ModNone, ct.Keyfn(ct.HidePortfolioUpdateMenu), ct.Views.Input.Name())
	ct.SetKeybindingMod('q', gocui.ModNone, ct.Keyfn(ct.HidePortfolioUpdateMenu), ct.Views.Input.Name())

	// keys to update portfolio holdings
	ct.SetKeybindingMod(gocui.KeyEnter, gocui.ModNone, ct.Keyfn(ct.SetPortfolioHoldings), ct.Views.Input.Name())

	// keys to quit convert menu when open
	ct.SetKeybindingMod(gocui.KeyEsc, gocui.ModNone, ct.Keyfn(ct.HideConvertMenu), ct.Views.ConvertMenu.Name())
	ct.SetKeybindingMod('q', gocui.ModNone, ct.Keyfn(ct.HideConvertMenu), ct.Views.ConvertMenu.Name())

	// mouse events
	ct.SetKeybindingMod(gocui.MouseRelease, gocui.ModNone, ct.Keyfn(ct.MouseRelease), "")
	ct.SetKeybindingMod(gocui.MouseLeft, gocui.ModNone, ct.Keyfn(ct.MouseLeftClick), "")
	ct.SetKeybindingMod(gocui.MouseMiddle, gocui.ModNone, ct.Keyfn(ct.MouseMiddleClick), "")
	ct.SetKeybindingMod(gocui.MouseRight, gocui.ModNone, ct.Keyfn(ct.MouseRightClick), "")
	ct.SetKeybindingMod(gocui.MouseWheelUp, gocui.ModNone, ct.Keyfn(ct.MouseWheelUp), "")
	ct.SetKeybindingMod(gocui.MouseWheelDown, gocui.ModNone, ct.Keyfn(ct.MouseWheelDown), "")

	// character key press to select option
	// TODO: use scrolling table
	keys := ct.SortedSupportedCurrencyConversions()
	for i, k := range keys {
		ct.SetKeybindingMod(rune(alphanumericcharacters[i]), gocui.ModNone, ct.Keyfn(ct.SetCurrencyConverstionFn(k)), ct.Views.ConvertMenu.Name())
	}

	return nil
}

// SetKeybindingMod sets the keybinding modifier key
func (ct *Cointop) SetKeybindingMod(key interface{}, mod gocui.Modifier, callback func(g *gocui.Gui, v *gocui.View) error, view string) error {
	var err error
	switch t := key.(type) {
	case gocui.Key:
		err = ct.g.SetKeybinding(view, t, mod, callback)
	case rune:
		err = ct.g.SetKeybinding(view, t, mod, callback)
	}
	return err
}

// Keyfn returns the keybinding function as a wrapped gocui view function
func (ct *Cointop) Keyfn(fn func() error) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return fn()
	}
}

// handleHkey handles the h key
func (ct *Cointop) handleHkey(key interface{}) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if k, ok := key.(rune); ok && k == 'h' && ct.IsPortfolioVisible() {
			ct.SortToggle("holdings", true)
		} else {
			ct.PrevPage()
		}
		return nil
	}
}

// Noop is a no-operation function
func (ct *Cointop) Noop() error {
	return nil
}
