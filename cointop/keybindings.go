package cointop

import (
	"log"

	"github.com/jroimartin/gocui"
)

func (ct *Cointop) keybindings(g *gocui.Gui) error {
	ct.setKeybinding(gocui.KeyArrowUp, ct.cursorUp)
	ct.setKeybinding(gocui.KeyArrowDown, ct.cursorDown)
	ct.setKeybinding(gocui.KeyArrowLeft, ct.prevPage)
	ct.setKeybinding(gocui.KeyArrowRight, ct.nextPage)
	ct.setKeybinding(gocui.KeyPgdn, ct.pageUp)
	ct.setKeybinding(gocui.KeyPgup, ct.pageUp)
	ct.setKeybinding(gocui.KeyHome, ct.navigateFirstLine)
	ct.setKeybinding(gocui.KeyEnd, ct.navigateLastLine)
	ct.setKeybinding(gocui.KeyEnter, ct.openLink)
	ct.setKeybinding(gocui.KeyEsc, ct.quit)
	ct.setKeybinding(gocui.KeySpace, ct.openLink)
	ct.setKeybinding(gocui.KeyCtrlC, ct.quit)
	ct.setKeybinding(gocui.KeyCtrlD, ct.pageDown)
	ct.setKeybinding(gocui.KeyCtrlN, ct.nextPage)
	ct.setKeybinding(gocui.KeyCtrlP, ct.prevPage)
	ct.setKeybinding(gocui.KeyCtrlR, ct.refresh)
	ct.setKeybinding(gocui.KeyCtrlU, ct.pageUp)
	ct.setKeybindingMod(gocui.KeyArrowUp, gocui.ModAlt, ct.sortAsc)
	ct.setKeybindingMod(gocui.KeyArrowDown, gocui.ModAlt, ct.sortDesc)
	ct.setKeybindingMod(gocui.KeyArrowLeft, gocui.ModAlt, ct.sortPrevCol)
	ct.setKeybindingMod(gocui.KeyArrowRight, gocui.ModAlt, ct.sortNextCol)
	ct.setKeybinding(gocui.KeyF1, ct.openHelp)
	ct.setKeybinding('0', ct.firstPage)
	ct.setKeybinding('1', ct.sortfn("1hchange", true))
	ct.setKeybinding('2', ct.sortfn("24hchange", true))
	ct.setKeybinding('7', ct.sortfn("7dchange", true))
	ct.setKeybinding('a', ct.sortfn("availablesupply", true))
	ct.setKeybinding('c', ct.toggleCoinChart)
	ct.setKeybinding('g', ct.navigateFirstLine)
	ct.setKeybinding('G', ct.navigateLastLine)
	ct.setKeybinding('h', ct.prevPage)
	ct.setKeybinding('H', ct.navigatePageFirstLine)
	ct.setKeybinding('j', ct.cursorDown)
	ct.setKeybinding('k', ct.cursorUp)
	ct.setKeybinding('l', ct.nextPage)
	ct.setKeybinding('L', ct.navigatePageLastLine)
	ct.setKeybinding('m', ct.sortfn("marketcap", true))
	ct.setKeybinding('M', ct.navigatePageMiddleLine)
	ct.setKeybinding('n', ct.sortfn("name", true))
	ct.setKeybinding('p', ct.sortfn("price", true))
	ct.setKeybinding('r', ct.sortfn("rank", false))
	ct.setKeybinding('s', ct.sortfn("symbol", false))
	ct.setKeybinding('t', ct.sortfn("totalsupply", true))
	ct.setKeybinding('u', ct.sortfn("lastupdated", true))
	ct.setKeybinding('v', ct.sortfn("24hvolume", true))
	ct.setKeybinding('q', ct.quit)
	ct.setKeybinding('$', ct.lastPage)
	ct.setKeybinding('?', ct.openHelp)
	return nil
}

func (ct *Cointop) setKeybinding(key interface{}, callback func(g *gocui.Gui, v *gocui.View) error) {
	var err error
	switch t := key.(type) {
	case gocui.Key:
		err = ct.g.SetKeybinding("", t, gocui.ModNone, callback)
	case rune:
		err = ct.g.SetKeybinding("", t, gocui.ModNone, callback)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func (ct *Cointop) setKeybindingMod(key interface{}, mod gocui.Modifier, callback func(g *gocui.Gui, v *gocui.View) error) {
	var err error
	switch t := key.(type) {
	case gocui.Key:
		err = ct.g.SetKeybinding("", t, mod, callback)
	case rune:
		err = ct.g.SetKeybinding("", t, mod, callback)
	}
	if err != nil {
		log.Fatal(err)
	}
}
