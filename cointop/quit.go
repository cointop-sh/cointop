package cointop

import (
	"os"

	"github.com/jroimartin/gocui"
)

func (ct *Cointop) quit() error {
	return gocui.ErrQuit
}

func (ct *Cointop) quitView() error {
	if ct.portfoliovisible {
		ct.portfoliovisible = false
		return ct.updateTable()
	}
	if ct.filterByFavorites {
		ct.filterByFavorites = false
		return ct.updateTable()
	}
	if ct.activeViewName() == ct.tableviewname {
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
