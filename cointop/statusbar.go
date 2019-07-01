package cointop

import (
	"fmt"

	"github.com/miguelmota/cointop/cointop/common/open"
	"github.com/miguelmota/cointop/cointop/common/pad"
)

func (ct *Cointop) updateStatusbar(s string) error {
	if ct.Views.Statusbar.Backing == nil {
		return nil
	}

	currpage := ct.currentDisplayPage()
	totalpages := ct.totalPagesDisplay()
	var quitText string
	var favoritesText string
	var portfolioText string
	if ct.State.portfolioVisible || ct.State.filterByFavorites {
		quitText = "Return"
	} else {
		quitText = "Quit"
	}
	if ct.State.portfolioVisible {
		portfolioText = "[E]Edit"
	} else {
		portfolioText = "[P]Portfolio"
	}
	if ct.State.filterByFavorites {
		favoritesText = "[Space]Unfavorite"
	} else {
		favoritesText = "[F]Favorites"
	}

	ct.update(func() {
		if ct.Views.Statusbar.Backing == nil {
			return
		}

		ct.Views.Statusbar.Backing.Clear()
		base := fmt.Sprintf("%s%s %sHelp %sChart %sRange %sSearch %sConvert %s %s %sSave", "[Q]", quitText, "[?]", "[Enter]", "[[ ]]", "[/]", "[C]", favoritesText, portfolioText, "[CTRL-S]")
		str := pad.Right(fmt.Sprintf("%v %sPage %v/%v %s", base, "[← →]", currpage, totalpages, s), ct.maxTableWidth, " ")
		v := fmt.Sprintf("v%s", ct.version())
		str = str[:len(str)-len(v)+2] + v
		fmt.Fprintln(ct.Views.Statusbar.Backing, str)
	})

	return nil
}

func (ct *Cointop) refreshRowLink() error {
	var shortcut string
	if !open.CommandExists() {
		shortcut = "[O]Open "
	}

	url := ct.rowLinkShort()
	ct.updateStatusbar(fmt.Sprintf("%s%s", shortcut, url))

	return nil
}
