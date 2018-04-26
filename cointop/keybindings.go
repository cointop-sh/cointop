package cointop

import (
	"strings"

	"github.com/miguelmota/cointop/pkg/gocui"
)

func (ct *Cointop) parseKeys(s string) (interface{}, gocui.Modifier) {
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
			case "0": // not supported?
			case "1": // not supported?
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
			case "9": // not supported?
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
				key = gocui.KeyCtrlLsqBracket
			case "]":
				key = gocui.KeyCtrlRsqBracket
			case "space":
				key = gocui.KeyCtrlSpace
			case "backslash":
				key = gocui.KeyCtrlBackslash
			case "underscore":
				key = gocui.KeyCtrlUnderscore
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
	}

	return key, mod
}

func (ct *Cointop) keybindings(g *gocui.Gui) error {
	for k, v := range ct.shortcutkeys {
		v = strings.TrimSpace(strings.ToLower(v))
		var fn func(g *gocui.Gui, v *gocui.View) error
		key, mod := ct.parseKeys(k)
		view := "table"
		switch v {
		case "move_up":
			fn = ct.keyfn(ct.cursorUp)
		case "move_down":
			fn = ct.keyfn(ct.cursorDown)
		case "previous_page":
			fn = ct.keyfn(ct.prevPage)
		case "next_page":
			fn = ct.keyfn(ct.nextPage)
		case "page_down":
			fn = ct.keyfn(ct.pageDown)
		case "page_up":
			fn = ct.keyfn(ct.pageUp)
		case "sort_column_symbol":
			fn = ct.sortfn("symbol", false)
		case "move_to_page_first_row":
			fn = ct.keyfn(ct.navigateFirstLine)
		case "move_to_page_last_row":
			fn = ct.keyfn(ct.navigateLastLine)
		case "open_link":
			fn = ct.keyfn(ct.openLink)
		case "refresh":
			fn = ct.keyfn(ct.refresh)
		case "sort_column_asc":
			fn = ct.keyfn(ct.sortAsc)
		case "sort_column_desc":
			fn = ct.keyfn(ct.sortDesc)
		case "sort_left_column":
			fn = ct.keyfn(ct.sortPrevCol)
		case "sort_right_column":
			fn = ct.keyfn(ct.sortNextCol)
		case "help":
			fallthrough
		case "toggle_show_help":
			fn = ct.keyfn(ct.toggleHelp)
			view = ""
		case "hide_help":
			fn = ct.keyfn(ct.hideHelp)
			view = "help"
		case "first_page":
			fn = ct.keyfn(ct.firstPage)
		case "sort_column_1h_change":
			fn = ct.sortfn("1hchange", true)
		case "sort_column_24h_change":
			fn = ct.sortfn("24hchange", true)
		case "sort_column_7d_change":
			fn = ct.sortfn("7dchange", true)
		case "sort_column_available_supply":
			fn = ct.sortfn("availablesupply", true)
		case "toggle_row_chart":
			fn = ct.keyfn(ct.toggleCoinChart)
		case "move_to_page_visible_first_row":
			fn = ct.keyfn(ct.navigatePageFirstLine)
		case "move_to_page_visible_last_row":
			fn = ct.keyfn(ct.navigatePageLastLine)
		case "sort_column_market_cap":
			fn = ct.sortfn("marketcap", true)
		case "move_to_page_visible_middle_row":
			fn = ct.keyfn(ct.navigatePageMiddleLine)
		case "sort_column_name":
			fn = ct.sortfn("name", true)
		case "sort_column_price":
			fn = ct.sortfn("price", true)
		case "sort_column_rank":
			fn = ct.sortfn("rank", false)
		case "sort_column_total_supply":
			fn = ct.sortfn("totalsupply", true)
		case "sort_column_last_updated":
			fn = ct.sortfn("lastupdated", true)
		case "sort_column_24h_volume":
			fn = ct.sortfn("24hvolume", true)
		case "last_page":
			fn = ct.keyfn(ct.lastPage)
		case "open_search":
			fn = ct.keyfn(ct.openSearch)
			view = ""
		case "toggle_favorite":
			fn = ct.keyfn(ct.toggleFavorite)
		case "toggle_show_favorites":
			fn = ct.keyfn(ct.toggleShowFavorites)
		case "save":
			fn = ct.keyfn(ct.save)
		case "quit":
			fn = ct.keyfn(ct.quit)
			view = ""
		default:
			fn = ct.keyfn(ct.noop)
		}

		ct.setKeybindingMod(key, mod, fn, view)
	}

	// keys to force quit
	ct.setKeybindingMod(gocui.KeyCtrlC, gocui.ModNone, ct.keyfn(ct.forceQuit), "")
	ct.setKeybindingMod(gocui.KeyCtrlZ, gocui.ModNone, ct.keyfn(ct.forceQuit), "")

	// searchfield keys
	ct.setKeybindingMod(gocui.KeyEnter, gocui.ModNone, ct.keyfn(ct.doSearch), "searchfield")
	ct.setKeybindingMod(gocui.KeyEsc, gocui.ModNone, ct.keyfn(ct.cancelSearch), "searchfield")

	// keys to quit help when open
	ct.setKeybindingMod(gocui.KeyEsc, gocui.ModNone, ct.keyfn(ct.hideHelp), "help")
	ct.setKeybindingMod('q', gocui.ModNone, ct.keyfn(ct.hideHelp), "help")
	ct.setKeybindingMod('x', gocui.ModNone, ct.keyfn(ct.hideHelp), "help")
	ct.setKeybindingMod('c', gocui.ModNone, ct.keyfn(ct.hideHelp), "help")
	return nil
}

func (ct *Cointop) setKeybindingMod(key interface{}, mod gocui.Modifier, callback func(g *gocui.Gui, v *gocui.View) error, view string) error {
	var err error
	switch t := key.(type) {
	case gocui.Key:
		err = ct.g.SetKeybinding(view, t, mod, callback)
	case rune:
		err = ct.g.SetKeybinding(view, t, mod, callback)
	}
	return err
}

func (ct *Cointop) keyfn(fn func() error) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return fn()
	}
}

func (ct *Cointop) noop() error {
	return nil
}
