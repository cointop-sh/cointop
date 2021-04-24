package cointop

import (
	"os"

	"github.com/miguelmota/gocui"
)

// Quit quits the program
func (ct *Cointop) Quit() error {
	ct.logfile.Close()
	return gocui.ErrQuit
}

// QuitView exists the current view
func (ct *Cointop) QuitView() error {
	ct.debuglog("quitView()")
	if ct.State.selectedView != CoinsView {
		ct.SetSelectedView(CoinsView)
		return ct.UpdateTable()
	}
	if ct.ui.ActiveViewName() == ct.Views.Table.Name() {
		return ct.Quit()
	}

	return nil
}

// Exit safely exits the program
func (ct *Cointop) Exit() {
	ct.debuglog("exit()")
	ct.logfile.Close()
	if ct.g != nil {
		ct.g.Close()
	} else {
		os.Exit(0)
	}
}
