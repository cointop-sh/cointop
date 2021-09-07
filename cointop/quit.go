package cointop

import (
	"os"

	"github.com/miguelmota/gocui"
	log "github.com/sirupsen/logrus"
)

// Quit quits the program
func (ct *Cointop) Quit() error {
	ct.logfile.Close()
	return gocui.ErrQuit
}

// QuitView exists the current view
func (ct *Cointop) QuitView() error {
	log.Debug("QuitView()")
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
	log.Debug("Exit()")
	ct.logfile.Close()
	if ct.g != nil {
		ct.g.Close()
	} else {
		os.Exit(0)
	}
}
