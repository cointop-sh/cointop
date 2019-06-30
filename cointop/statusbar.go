package cointop

import (
	"fmt"

	"github.com/miguelmota/cointop/cointop/common/open"
	"github.com/miguelmota/cointop/cointop/common/pad"
)

func (ct *Cointop) updateStatusbar(s string) error {
	if ct.statusbarview == nil {
		return nil
	}

	currpage := ct.currentDisplayPage()
	totalpages := ct.totalPagesDisplay()
	var quitText string
	var favoritesText string
	var portfolioText string
	if ct.portfoliovisible || ct.filterByFavorites {
		quitText = "Return"
	} else {
		quitText = "Quit"
	}
	if ct.portfoliovisible {
		portfolioText = "[E]Edit"
	} else {
		portfolioText = "[P]Portfolio"
	}
	if ct.filterByFavorites {
		favoritesText = "[Space]Unfavorite"
	} else {
		favoritesText = "[F]Favorites"
	}

	ct.update(func() {
		if ct.statusbarview == nil {
			return
		}

		ct.statusbarview.Clear()
		base := fmt.Sprintf("%s%s %sHelp %sChart %sRange %sSearch %sConvert %s %s %sSave", "[Q]", quitText, "[?]", "[Enter]", "[[ ]]", "[/]", "[C]", favoritesText, portfolioText, "[CTRL-S]")
		str := pad.Right(fmt.Sprintf("%v %sPage %v/%v %s", base, "[← →]", currpage, totalpages, s), ct.maxtablewidth, " ")
		v := fmt.Sprintf("v%s", ct.version())
		str = str[:len(str)-len(v)+2] + v
		fmt.Fprintln(ct.statusbarview, str)
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
