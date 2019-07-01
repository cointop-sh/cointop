package cointop

import (
	"os"

	"github.com/jroimartin/gocui"
)

func (ct *Cointop) quit() error {
	return gocui.ErrQuit
}

func (ct *Cointop) quitView() error {
	if ct.State.portfolioVisible {
		ct.State.portfolioVisible = false
		return ct.updateTable()
	}
	if ct.State.filterByFavorites {
		ct.State.filterByFavorites = false
		return ct.updateTable()
	}
	if ct.activeViewName() == ct.Views.Table.Name {
		return ct.quit()
	}

	return nil
}

// Exit safely exit application
func (ct *Cointop) Exit() {
	if ct.g != nil {
		ct.g.Close()
	} else {
		os.Exit(0)
	}
}
