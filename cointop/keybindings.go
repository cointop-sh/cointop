package cointop

import (
	"strings"
	"unicode"

	"github.com/cointop-sh/cointop/pkg/gocui"
	"github.com/gdamore/tcell/v2"
	log "github.com/sirupsen/logrus"
)

// keyMap translates key alternative names to a canonical version
func keyMap(k string) string {
	key := k
	switch strings.ToLower(k) {
	case "lsqrbracket", "leftsqrbracket", "leftsquarebracket":
		key = "["
	case "rsqrbracket", "rightsqrbracket", "rightsquarebracket":
		key = "]"
	case "space", "spacebar":
		key = " " // with meta should be "space"
	case "\\\\", "backslash":
		key = "\\"
	case "underscore":
		key = "_"
	case "arrowup", "uparrow":
		key = "Up"
	case "arrowdown", "downarrow":
		key = "Down"
	case "arrowleft", "leftarrow":
		key = "Left"
	case "arrowright", "rightarrow":
		key = "Right"
	case "return":
		key = "Enter"
	case "escape":
		key = "Esc"
	case "pageup":
		key = "PgUp"
	case "pagedown", "pgdown":
		key = "PgDn"
	}
	return key
}

// ParseKeys returns string keyboard key as gocui key type
func (ct *Cointop) ParseKeys(s string) (interface{}, tcell.ModMask) {
	// TODO: change file convention to match tcell (no aliases, dash between mod and key)
	// TODO: change to return EventKey?
	var key interface{}
	mod := tcell.ModNone

	// translate legacy and special names for keys
	keyName := keyMap(s)
	if len(s) > 1 {
		keyName := strings.TrimSpace(strings.Replace(s, "+", "-", -1))
		split := strings.Split(keyName, "-")

		if len(split) > 1 {
			m := strings.ToLower(strings.TrimSpace(split[0]))
			k := strings.TrimSpace(split[1])
			k = keyMap(k)
			if k == " " {
				k = "Space" // fix mod+space
			}

			if m == "alt" {
				mod = tcell.ModAlt
				keyName = k
			} else if m == "ctrl" {
				// let the lookup handle it
				keyName = m + "-" + k
			} else {
				keyName = m + "-" + k
			}
			// TODO: other mods?
		}
	}

	// First try looking up keyname directly
	lcKeyName := strings.ToLower(keyName)
	for key, name := range tcell.KeyNames {
		if strings.ToLower(name) == lcKeyName {
			if strings.HasPrefix(name, "Ctrl-") {
				mod = tcell.ModCtrl
			}
			return key, mod
		}
	}

	// Then try one-rune variants
	if len(keyName) == 1 {
		r := []rune(keyName)
		key = r[0]

		return key, mod
	}

	if key == nil {
		log.Debugf("Could not map key descriptio '%s' to key", s)
	}
	return key, mod
}

// SetKeybindingAction maps a shortcut key to an action
func (ct *Cointop) SetKeybindingAction(shortcutKey string, action string) error {
	if shortcutKey == "" {
		return nil
	}
	action = strings.TrimSpace(strings.ToLower(action))
	var fn func(g *gocui.Gui, v *gocui.View) error
	key, mod := ct.ParseKeys(shortcutKey)
	view := "table"
	switch action {
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
	case "help", "toggle_show_help":
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
	case "move_to_first_page_first_row":
		fn = ct.Keyfn(ct.NavigateToFirstPageFirstRow)
	case "sort_column_1h_change":
		fn = ct.Sortfn("1h_change", true)
	case "sort_column_24h_change":
		fn = ct.Sortfn("24h_change", true)
	case "sort_column_7d_change":
		fn = ct.Sortfn("7d_change", true)
	case "sort_column_30d_change":
		fn = ct.Sortfn("30d_change", true)
	case "sort_column_1y_change":
		fn = ct.Sortfn("1y_change", true)
	case "sort_column_available_supply":
		fn = ct.Sortfn("available_supply", true)
	case "toggle_row_chart":
		fn = ct.Keyfn(ct.ToggleCoinChart)
	case "move_to_page_visible_first_row":
		fn = ct.Keyfn(ct.NavigatePageFirstLine)
	case "move_to_page_visible_last_row":
		fn = ct.Keyfn(ct.navigatePageLastLine)
	case "sort_column_market_cap":
		fn = ct.Sortfn("market_cap", true)
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
		fn = ct.Sortfn("total_supply", true)
	case "sort_column_last_updated":
		fn = ct.Sortfn("last_updated", true)
	case "sort_column_24h_volume":
		fn = ct.Sortfn("24h_volume", true)
	case "sort_column_balance":
		fn = ct.Sortfn("balance", true)
	case "sort_column_holdings":
		fn = ct.Sortfn("holdings", true)
	case "sort_column_percent_holdings":
		fn = ct.Sortfn("percent_holdings", true)
	case "last_page":
		fn = ct.Keyfn(ct.LastPage)
	case "open_search":
		fn = ct.Keyfn(ct.OpenSearch)
		view = ""
	case "toggle_price_alerts":
		fn = ct.Keyfn(ct.TogglePriceAlerts)
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
	case "toggle_portfolio_balances":
		fn = ct.Keyfn(ct.TogglePortfolioBalances)
	case "show_portfolio_edit_menu":
		fn = ct.Keyfn(ct.TogglePortfolioUpdateMenu)
	case "show_price_alert_edit_menu":
		fn = ct.Keyfn(ct.ShowPriceAlertsUpdateMenu)
	case "show_price_alert_add_menu":
		fn = ct.Keyfn(ct.ShowPriceAlertsAddMenu)
	case "toggle_table_fullscreen":
		fn = ct.Keyfn(ct.ToggleTableFullscreen)
		view = ""
	case "toggle_chart_fullscreen":
		fn = ct.Keyfn(ct.ToggleChartFullscreen)
		view = ""
	case "enlarge_chart":
		fn = ct.Keyfn(ct.EnlargeChart)
	case "shorten_chart":
		fn = ct.Keyfn(ct.ShortenChart)
	case "move_down_or_next_page":
		fn = ct.Keyfn(ct.CursorDownOrNextPage)
	case "move_up_or_previous_page":
		fn = ct.Keyfn(ct.CursorUpOrPreviousPage)
	case "sort_column_cost":
		fn = ct.Sortfn("cost", true)
	case "sort_column_pnl":
		fn = ct.Sortfn("pnl", true)
	case "sort_column_pnl_percent":
		fn = ct.Sortfn("pnl_percent", true)
	default:
		fn = ct.Keyfn(ct.Noop)
	}

	ct.SetKeybindingMod(key, mod, fn, view)

	return nil
}

// SetKeybindings sets keyboard shortcut key bindings
func (ct *Cointop) SetKeybindings() error {
	for k, v := range ct.State.shortcutKeys {
		if err := ct.SetKeybindingAction(k, v); err != nil {
			return err
		}
	}

	// keys to force quit
	ct.SetKeybindingMod(tcell.KeyCtrlC, tcell.ModNone, ct.Keyfn(ct.Quit), "")
	ct.SetKeybindingMod(tcell.KeyCtrlZ, tcell.ModNone, ct.Keyfn(ct.Quit), "")

	// searchfield keys
	ct.SetKeybindingMod(tcell.KeyEnter, tcell.ModNone, ct.Keyfn(ct.DoSearch), ct.Views.SearchField.Name())
	ct.SetKeybindingMod(tcell.KeyEsc, tcell.ModNone, ct.Keyfn(ct.CancelSearch), ct.Views.SearchField.Name())

	// keys to quit help when open
	ct.SetKeybindingMod(tcell.KeyEsc, tcell.ModNone, ct.Keyfn(ct.HideHelp), ct.Views.Menu.Name())
	ct.SetKeybindingMod('q', tcell.ModNone, ct.Keyfn(ct.HideHelp), ct.Views.Menu.Name())

	// keys to quit portfolio update menu when open
	ct.SetKeybindingMod(tcell.KeyEsc, tcell.ModNone, ct.Keyfn(ct.HidePortfolioUpdateMenu), ct.Views.Input.Name())
	ct.SetKeybindingMod('q', tcell.ModNone, ct.Keyfn(ct.HidePortfolioUpdateMenu), ct.Views.Input.Name())

	// keys to quit convert menu when open
	ct.SetKeybindingMod(tcell.KeyEsc, tcell.ModNone, ct.Keyfn(ct.HideConvertMenu), ct.Views.Menu.Name())
	ct.SetKeybindingMod('q', tcell.ModNone, ct.Keyfn(ct.HideConvertMenu), ct.Views.Menu.Name())

	// keys to update portfolio holdings
	ct.SetKeybindingMod(tcell.KeyEnter, tcell.ModNone, ct.Keyfn(ct.EnterKeyPressHandler), ct.Views.Input.Name())

	// Work around issue with key-binding for '/' interfering with expressions
	key, mod := ct.ParseKeys("/")
	ct.DeleteKeybindingMod(key, mod, "")

	// mouse events
	ct.SetMousebindingMod(tcell.Button1, tcell.ModNone, ct.Keyfn(ct.MouseLeftClick), ct.Views.Table.Name()) // click to focus

	// clicking table headers sorts table
	ct.SetMousebindingMod(tcell.Button1, tcell.ModNone, ct.Keyfn(ct.TableHeaderMouseLeftClick), ct.Views.TableHeader.Name())
	ct.SetMousebindingMod(tcell.Button1, tcell.ModNone, ct.Keyfn(ct.StatusbarMouseLeftClick), ct.Views.Statusbar.Name())
	// debug mouse clicks
	ct.SetMousebindingMod(tcell.Button1, tcell.ModNone, ct.Keyfn(ct.MouseDebug), "")

	ct.SetMousebindingMod(tcell.WheelUp, tcell.ModNone, ct.Keyfn(ct.CursorUpOrPreviousPage), ct.Views.Table.Name())
	ct.SetMousebindingMod(tcell.WheelDown, tcell.ModNone, ct.Keyfn(ct.CursorDownOrNextPage), ct.Views.Table.Name())

	// character key press to select option
	// TODO: use scrolling table
	keys := ct.SortedSupportedCurrencyConversions()
	for i, k := range keys {
		ct.SetKeybindingMod(alphanumericcharacters[i], tcell.ModNone, ct.Keyfn(ct.SetCurrencyConverstionFn(k)), ct.Views.Menu.Name())
	}
	ct.SetMousebindingMod(tcell.Button1, tcell.ModNone, ct.Keyfn(ct.ConversionMouseLeftClick), ct.Views.Menu.Name())
	return nil
}

// MouseDebug emit a debug message about which View and coordinates are in MouseClick
func (ct *Cointop) MouseDebug() error {
	v, x, y, err := ct.g.GetViewRelativeMousePosition(ct.g.CurrentEvent)
	if err != nil {
		return err
	}
	log.Debugf("XXX MouseDebug view=%s %d,%d", v.Name(), x, y)
	return nil
}

// SetKeybindingMod sets the keybinding modifier key
func (ct *Cointop) SetKeybindingMod(key interface{}, mod tcell.ModMask, callback func(g *gocui.Gui, v *gocui.View) error, view string) error {
	// TODO: take EventKey?
	var err error
	switch t := key.(type) {
	case tcell.Key:
		err = ct.g.SetKeybinding(view, t, 0, mod, callback)
	case rune:
		err = ct.g.SetKeybinding(view, tcell.KeyRune, t, mod, callback)
		if err != nil {
			return err
		}

		// Binding Shift+[key] if key is uppercase and modifiers missing Shift
		// to support using on Windows
		if unicode.ToUpper(t) == t && (tcell.ModShift&mod == 0) {
			err = ct.g.SetKeybinding(view, tcell.KeyRune, t, mod|tcell.ModShift, callback)
		}
	}
	return err
}

// SetMousebindingMod adds a binding for a mouse eventdef
func (ct *Cointop) SetMousebindingMod(btn tcell.ButtonMask, mod tcell.ModMask, callback func(g *gocui.Gui, v *gocui.View) error, view string) error {
	return ct.g.SetMousebinding(view, btn, mod, callback)
}

// DeleteKeybinding ...
func (ct *Cointop) DeleteKeybinding(shortcutKey string) error {
	key, mod := ct.ParseKeys(shortcutKey)
	return ct.DeleteKeybindingMod(key, mod, "")
}

// DeleteKeybindingMod ...
func (ct *Cointop) DeleteKeybindingMod(key interface{}, mod tcell.ModMask, view string) error {
	// TODO: take EventKey
	var err error
	switch t := key.(type) {
	case tcell.Key:
		err = ct.g.DeleteKeybinding(view, t, 0, mod)
	case rune:
		err = ct.g.DeleteKeybinding(view, tcell.KeyRune, t, mod)
	}
	return err
}

// Keyfn returns the keybinding function as a wrapped gocui view function
func (ct *Cointop) Keyfn(fn func() error) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if fn != nil {
			return fn()
		}
		return nil
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
