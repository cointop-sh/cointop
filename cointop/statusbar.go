package cointop

import (
	"fmt"
	"unicode/utf8"

	"github.com/miguelmota/cointop/pkg/open"
	"github.com/miguelmota/cointop/pkg/pad"
	"github.com/miguelmota/cointop/pkg/ui"
)

// StatusbarView is structure for statusbar view
type StatusbarView = ui.View

// NewStatusbarView returns a new statusbar view
func NewStatusbarView() *StatusbarView {
	var view *StatusbarView = ui.NewView("statusbar")
	return view
}

// UpdateStatusbar updates the statusbar view
func (ct *Cointop) UpdateStatusbar(s string) error {
	ct.debuglog("UpdateStatusbar()")
	currpage := ct.CurrentDisplayPage()
	totalpages := ct.TotalPagesDisplay()
	var quitText string
	var favoritesText string
	var portfolioText string
	if ct.IsPortfolioVisible() || ct.IsFavoritesVisible() {
		quitText = "Return"
	} else {
		quitText = "Quit"
	}
	if ct.IsPortfolioVisible() {
		portfolioText = "[E]Edit"
	} else {
		portfolioText = "[P]Portfolio"
	}
	if ct.IsFavoritesVisible() {
		favoritesText = "[Space]Unfavorite"
	} else {
		favoritesText = "[F]Favorites"
	}

	helpStr := fmt.Sprintf("%s%s %sHelp", "[Q]", quitText, "[?]")
	var content string
	if ct.IsPriceAlertsVisible() {
		var editStr string
		if ct.ActivePriceAlertsLen() > 0 {
			editStr = "[E]Edit "
		}
		content = fmt.Sprintf("%s %s[+]Add", helpStr, editStr)
	} else {
		base := fmt.Sprintf("%s %sChart %sRange %sSearch %sConvert %s %s", helpStr, "[Enter]", "[[ ]]", "[/]", "[C]", favoritesText, portfolioText)
		str := pad.Right(fmt.Sprintf("%v %sPage %v/%v %s", base, "[← →]", currpage, totalpages, s), ct.Width(), " ")
		v := ct.Version()
		size := utf8.RuneCountInString(str)
		end := size - utf8.RuneCountInString(v) + 2
		if end > size {
			end = size
		}

		content = str[:end] + v
	}

	ct.UpdateUI(func() error {
		return ct.Views.Statusbar.Update(content)
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
