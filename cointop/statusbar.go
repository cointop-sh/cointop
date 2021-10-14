package cointop

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/cointop-sh/cointop/pkg/open"
	"github.com/cointop-sh/cointop/pkg/pad"
	"github.com/cointop-sh/cointop/pkg/ui"
	log "github.com/sirupsen/logrus"
)

// StatusbarView is structure for statusbar view
type StatusbarView = ui.View

// NewStatusbarView returns a new statusbar view
func NewStatusbarView() *StatusbarView {
	return ui.NewView("statusbar")
}

// UpdateStatusbar updates the statusbar view
func (ct *Cointop) UpdateStatusbar(s string) error {
	log.Debug("UpdateStatusbar()")
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
	log.Debug("RefreshRowLink()")
	var shortcut string
	if !open.CommandExists() {
		shortcut = "[O]Open "
	}

	url := ct.RowLinkShort()
	ct.UpdateStatusbar(fmt.Sprintf("%s%s", shortcut, url))

	return nil
}

// StatusbarMouseLeftClick is called on mouse left click event
func (ct *Cointop) StatusbarMouseLeftClick() error {
	_, x, _, err := ct.g.GetViewRelativeMousePosition(ct.g.CurrentEvent)
	if err != nil {
		return err
	}

	// Parse the statusbar text to identify hotspots and actions
	b := make([]byte, 1000)
	ct.Views.Statusbar.Rewind()
	if n, err := ct.Views.Statusbar.Read(b); err == nil {

		// Find all the "[X]word" substrings, then look for the one that was clicked
		matches := regexp.MustCompile(`\[.*?\]\w+`).FindAllIndex(b[:n], -1)
		for _, match := range matches {
			if x >= match[0] && x <= match[1] {
				s := string(b[match[0]:match[1]])
				bits := strings.Split(s, "]")
				word := bits[len(bits)-1]

				// Quit/Return Help Chart Range Search Convert Favorites Portfolio Edit(portfolio) Unfavorite
				switch word {
				case "Help":
					ct.ToggleHelp()
				case "Range":
					// left hand edge of "Range" is Prev, the rest is Next
					if x-match[0] < 3 {
						ct.PrevChartRange()
					} else {
						ct.NextChartRange()
					}
				case "Search":
					ct.OpenSearch()
				case "Convert":
					ct.ToggleConvertMenu()
				case "Favorites":
					ct.ToggleSelectedView(FavoritesView)
				case "Portfolio":
					ct.ToggleSelectedView(PortfolioView)
				}

			}
		}
	}

	return nil
}
