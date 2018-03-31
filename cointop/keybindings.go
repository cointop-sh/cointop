package cointop

import (
	"log"
	"os/exec"

	"github.com/jroimartin/gocui"
)

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

func (ct *Cointop) keybindings(g *gocui.Gui) error {
	ct.setKeybinding(gocui.KeyArrowDown, ct.cursorDown)
	ct.setKeybinding('j', ct.cursorDown)
	ct.setKeybinding(gocui.KeyArrowUp, ct.cursorUp)
	ct.setKeybinding('k', ct.cursorUp)
	ct.setKeybinding(gocui.KeyCtrlD, ct.pageDown)
	ct.setKeybinding(gocui.KeyCtrlU, ct.pageUp)
	ct.setKeybinding('r', ct.sortfn("rank", false))
	ct.setKeybinding('n', ct.sortfn("name", true))
	ct.setKeybinding('s', ct.sortfn("symbol", false))
	ct.setKeybinding('p', ct.sortfn("price", true))
	ct.setKeybinding('m', ct.sortfn("marketcap", true))
	ct.setKeybinding('v', ct.sortfn("24hvolume", true))
	ct.setKeybinding('1', ct.sortfn("1hchange", true))
	ct.setKeybinding('2', ct.sortfn("24hchange", true))
	ct.setKeybinding('7', ct.sortfn("7dchange", true))
	ct.setKeybinding('t', ct.sortfn("totalsupply", true))
	ct.setKeybinding('a', ct.sortfn("availablesupply", true))
	ct.setKeybinding('l', ct.sortfn("lastupdated", true))
	ct.setKeybinding(gocui.KeyEnter, ct.enter)
	ct.setKeybinding(gocui.KeySpace, ct.enter)
	ct.setKeybinding(gocui.KeyCtrlC, ct.quit)
	ct.setKeybinding('q', ct.quit)
	ct.setKeybinding(gocui.KeyEsc, ct.quit)
	return nil
}

func (ct *Cointop) enter(g *gocui.Gui, v *gocui.View) error {
	exec.Command("open", ct.rowLink()).Output()
	return nil
}

func (ct *Cointop) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
