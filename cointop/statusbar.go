package cointop

import (
	"fmt"

	"github.com/miguelmota/cointop/cointop/common/pad"
)

func (ct *Cointop) updateStatusbar(s string) {
	currpage := ct.currentDisplayPage()
	totalpages := ct.totalPages()
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
		ct.statusbarview.Clear()
		base := fmt.Sprintf("%s%s %sHelp %sChart %sRange %sSearch %sConvert %s %s %sSave", "[Q]", quitText, "[?]", "[Enter]", "[[ ]]", "[/]", "[C]", favoritesText, portfolioText, "[CTRL-S]")
		str := pad.Right(fmt.Sprintf("%v %sPage %v/%v %s", base, "[← →]", currpage, totalpages, s), ct.maxtablewidth, " ")
		v := fmt.Sprintf("v%s", ct.version())
		str = str[:len(str)-len(v)+2] + v
		fmt.Fprintln(ct.statusbarview, str)
	})
}

func (ct *Cointop) refreshRowLink() {
	url := ct.rowLinkShort()
	ct.updateStatusbar(fmt.Sprintf("%sOpen %s", "[O]", url))
}
