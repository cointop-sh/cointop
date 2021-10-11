package cointop

import (
	"regexp"
	"strings"

	"github.com/cointop-sh/cointop/pkg/ui"
	log "github.com/sirupsen/logrus"
)

// SearchFieldView is structure for search field view
type SearchFieldView = ui.View

// NewSearchFieldView returns a new search field view
func NewSearchFieldView() *SearchFieldView {
	return ui.NewView("searchfield")
}

// InputView is structure for help view
type InputView = ui.View

// NewInputView returns a new help view
func NewInputView() *InputView {
	return ui.NewView("input")
}

// OpenSearch opens the search field
func (ct *Cointop) OpenSearch() error {
	log.Debug("OpenSearch()")
	if ct.ui.ActiveViewName() != ct.Views.Table.Name() {
		return nil
	}
	ct.State.searchFieldVisible = true
	ct.ui.SetCursor(true)
	ct.SetActiveView(ct.Views.SearchField.Name())
	return nil
}

// CancelSearch closes the search field
func (ct *Cointop) CancelSearch() error {
	log.Debug("CancelSearch()")
	ct.State.searchFieldVisible = false
	ct.ui.SetCursor(false)
	ct.SetActiveView(ct.Views.Table.Name())
	return nil
}

// DoSearch triggers the search and sets views
func (ct *Cointop) DoSearch() error {
	log.Debug("DoSearch()")
	ct.Views.SearchField.Rewind()
	b := make([]byte, 100)
	n, err := ct.Views.SearchField.Read(b)
	if err != nil {
		return err
	}

	// TODO: do this a better way (SoC)
	ct.SetSelectedView(CoinsView)

	defer ct.SetActiveView(ct.Views.Table.Name())
	if err != nil {
		return nil
	}
	if n == 0 {
		return nil
	}
	q := string(b)
	// remove slash
	regex := regexp.MustCompile(`/(.*)`)
	matches := regex.FindStringSubmatch(q)
	if len(matches) > 0 {
		q = matches[1]
	}
	return ct.Search(q)
}

// Search performs the search and filtering
func (ct *Cointop) Search(q string) error {
	log.Debugf("Search(%s)", q)

	// If there are no coins, return no result
	if len(ct.State.coins) == 0 {
		return nil
	}

	// If search term is empty, use the previous search term.
	q = strings.TrimSpace(strings.ToLower(q))
	if q == "" {
		q = ct.State.lastSearchQuery
	} else {
		ct.State.lastSearchQuery = q
	}

	// Start the search from the current position (+1), looking names that start with the search term, or symbols that match completely
	currentIndex := ct.GetGlobalCoinIndex(ct.HighlightedRowCoin()) + 1
	if ct.IsLastPage() && ct.IsLastRow() {
		currentIndex = 0
	}
	for i, coin := range ct.State.allCoins[currentIndex:] {
		name := strings.ToLower(coin.Name)
		symbol := strings.ToLower(coin.Symbol)
		if strings.HasPrefix(name, q) || symbol == q {
			ct.GoToGlobalIndex(currentIndex + i)
			return nil
		}
	}
	return nil
}
