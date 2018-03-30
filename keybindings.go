package main

import (
	"log"

	"github.com/jroimartin/gocui"
)

func (ct *Cointop) setKeybinding(key gocui.Key, callback func(g *gocui.Gui, v *gocui.View) error) {
	if err := ct.g.SetKeybinding("", key, gocui.ModNone, callback); err != nil {
		log.Fatal(err)
	}
}

func (ct *Cointop) keybindings(g *gocui.Gui) error {
	ct.setKeybinding(gocui.KeyArrowDown, ct.cursorDown)
	ct.setKeybinding('j', ct.cursorDown)
	ct.setKeybinding(gocui.KeyArrowUp, ct.cursorUp)
	ct.setKeybinding('k', ct.cursorUp)
	ct.setKeybinding(gocui.KeyCtrlD, ct.pageDown)
	ct.setKeybinding(gocui.KeyCtrlU, ct.pageUp)
	ct.setKeybinding('r', ct.sort("rank", false))
	ct.setKeybinding('n', ct.sort("name", true))
	ct.setKeybinding('s', ct.sort("symbol", false))
	ct.setKeybinding('p', ct.sort("price", true))
	ct.setKeybinding('m', ct.sort("marketcap", true))
	ct.setKeybinding('v', ct.sort("24hvolume", true))
	ct.setKeybinding('1', ct.sort("1hchange", true))
	ct.setKeybinding('2', ct.sort("24hchange", true))
	ct.setKeybinding('7', ct.sort("7dchange", true))
	ct.setKeybinding('t', ct.sort("totalsupply", true))
	ct.setKeybinding('a', ct.sort("availablesupply", true))
	ct.setKeybinding('l', ct.sort("lastupdated", true))
	ct.setKeybinding(gocui.KeyCtrlC, ct.quit)
	ct.setKeybinding('q', ct.quit)
	ct.setKeybinding(gocui.KeyEsc, ct.quit)
	return nil
}
