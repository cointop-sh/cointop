package cointop

import (
	"fmt"

	"github.com/miguelmota/cointop/cointop/common/open"
	"github.com/miguelmota/cointop/cointop/common/pad"
)

// StatusbarView is structure for statusbar view
type StatusbarView struct {
	*View
}

// NewStatusbarView returns a new statusbar view
func NewStatusbarView() *StatusbarView {
	return &StatusbarView{NewView("statusbar")}
}

// Update updates the content of the statusbar
func (statusbar *StatusbarView) Update(str string) error {
	if statusbar.Backing() == nil {
		return nil
	}

	statusbar.Backing().Clear()
	fmt.Fprintln(statusbar.Backing(), str)

	return nil
}

// UpdateStatusbar updates the statusbar view
func (ct *Cointop) UpdateStatusbar(s string) error {
	ct.debuglog("UpdateStatusbar()")
	currpage := ct.CurrentDisplayPage()
	totalpages := ct.TotalPagesDisplay()
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

	base := fmt.Sprintf("%s%s %sHelp %sChart %sRange %sSearch %sConvert %s %s %sSave", "[Q]", quitText, "[?]", "[Enter]", "[[ ]]", "[/]", "[C]", favoritesText, portfolioText, "[CTRL-S]")
	str := pad.Right(fmt.Sprintf("%v %sPage %v/%v %s", base, "[← →]", currpage, totalpages, s), ct.maxTableWidth, " ")
	v := fmt.Sprintf("v%s", ct.Version())
	end := len(str) - len(v) + 2
	if end > len(str) {
		end = len(str)
	}

	str = str[:end] + v

	ct.Update(func() error {
		ct.Views.Statusbar.Update(str)
		return nil
	})

	return nil
}

// RefreshRowLink updates the row link in the statusbar
func (ct *Cointop) RefreshRowLink() error {
	ct.debuglog("RefreshRowLink()")
	var shortcut string
	if !open.CommandExists() {
		shortcut = "[O]Open "
	}

	url := ct.RowLinkShort()
	ct.UpdateStatusbar(fmt.Sprintf("%s%s", shortcut, url))

	return nil
}
