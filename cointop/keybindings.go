package cointop

import (
	"strings"

	"github.com/jroimartin/gocui"
)

// defaults
var shortcutkeys = map[string]string{
	"arrowup":        "moveup",
	"arrowdown":      "movedown",
	"arrowleft":      "prevpage",
	"arrowright":     "nextpage",
	"pagedown":       "pagedown",
	"pageup":         "pageup",
	"home":           "movepagefirstrow",
	"end":            "movepagelastrow",
	"enter":          "openlink",
	"esc":            "quit",
	"space":          "openlink",
	"ctrl+c":         "quit",
	"ctrl+d":         "pagedown",
	"ctrl+n":         "nextpage",
	"ctrl+p":         "prevpage",
	"ctrl+r":         "refresh",
	"ctrl+u":         "pageup",
	"alt+arrowup":    "sortcolasc",
	"alt+arrowdown":  "sortcoldesc",
	"alt+arrowleft":  "sortleftcol",
	"alt+arrowright": "sortrightcol",
	"f1":             "help",
	"0":              "movefirstpage",
	"1":              "sortcol1hchange",
	"2":              "sortcol24hchange",
	"7":              "sortcol7dchange",
	"a":              "sortcolavailablesupply",
	"c":              "togglerowchart",
	"g":              "movepagefirstrow",
	"G":              "movepagelastrow",
	"h":              "prevpage",
	"H":              "movepagevisiblefirstrow",
	"j":              "movedown",
	"k":              "moveup",
	"l":              "nextpage",
	"L":              "movepagevisiblelastrow",
	"m":              "sortcolmarketcap",
	"M":              "movepagevisiblemiddlerow",
	"n":              "sortcolname",
	"p":              "sortcolprice",
	"r":              "sortcolrank",
	"s":              "sortcolsymbol",
	"t":              "sortcoltotalsupply",
	"u":              "sortcollastupdated",
	"v":              "sortcol24hvolume",
	"q":              "quit",
	"$":              "movelastpage",
	"?":              "help",
}

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
		key = gocui.KeyPgdn
	case "home":
		key = gocui.KeyHome
	case "end":
		key = gocui.KeyEnd
	}

	return key, mod
}

func (ct *Cointop) keybindings(g *gocui.Gui) error {
	for k, v := range shortcutkeys {
		var fn func(g *gocui.Gui, v *gocui.View) error
		key, mod := ct.parseKeys(k)
		switch v {
		case "moveup":
			fn = ct.cursorUp
		case "movedown":
			fn = ct.cursorDown
		case "prevpage":
			fn = ct.prevPage
		case "nextpage":
			fn = ct.nextPage
		case "pagedown":
			fn = ct.pageDown
		case "pageup":
			fn = ct.pageUp
		case "sortcolsymbol":
			fn = ct.sortfn("symbol", false)
		case "movepagefirstrow":
			fn = ct.navigateFirstLine
		case "movepagelastrow":
			fn = ct.navigateLastLine
		case "openlink":
			fn = ct.openLink
		case "refresh":
			fn = ct.refresh
		case "sortcolasc":
			fn = ct.sortAsc
		case "sortcoldesc":
			fn = ct.sortDesc
		case "sortleftcol":
			fn = ct.sortPrevCol
		case "sortrightcol":
			fn = ct.sortNextCol
		case "help":
			fn = ct.openHelp
		case "movefirstpage":
			fn = ct.firstPage
		case "sortcol1hchange":
			fn = ct.sortfn("1hchange", true)
		case "sortcol24hchange":
			fn = ct.sortfn("24hchange", true)
		case "sortcol7dchange":
			fn = ct.sortfn("7dchange", true)
		case "sortcolavailablesupply":
			fn = ct.sortfn("availablesupply", true)
		case "togglerowchart":
			fn = ct.toggleCoinChart
		case "movepagevisiblefirstrow":
			fn = ct.navigatePageFirstLine
		case "movepagevisiblelastrow":
			fn = ct.navigatePageLastLine
		case "sortcolmarketcap":
			fn = ct.sortfn("marketcap", true)
		case "movepagevisiblemiddlerow":
			fn = ct.navigatePageMiddleLine
		case "sortcolname":
			fn = ct.sortfn("name", true)
		case "sortcolprice":
			fn = ct.sortfn("price", true)
		case "sortcolrank":
			fn = ct.sortfn("rank", false)
		case "sortcoltotalsupply":
			fn = ct.sortfn("totalsupply", true)
		case "sortcollastupdated":
			fn = ct.sortfn("lastupdated", true)
		case "sortcol24hvolume":
			fn = ct.sortfn("24hvolume", true)
		case "movelastpage":
			fn = ct.lastPage
		case "quit":
			fn = ct.quit
		default:
			fn = keynoop
		}

		ct.setKeybindingMod(key, mod, fn)
	}

	return nil
}

func (ct *Cointop) setKeybindingMod(key interface{}, mod gocui.Modifier, callback func(g *gocui.Gui, v *gocui.View) error) error {
	var err error
	switch t := key.(type) {
	case gocui.Key:
		err = ct.g.SetKeybinding("", t, mod, callback)
	case rune:
		err = ct.g.SetKeybinding("", t, mod, callback)
	}
	return err
}

func keynoop(g *gocui.Gui, v *gocui.View) error {
	return nil
}
