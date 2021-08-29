package cointop

import (
	"regexp"
	"strings"

	"github.com/miguelmota/cointop/pkg/levenshtein"
	"github.com/miguelmota/cointop/pkg/ui"
)

// SearchFieldView is structure for search field view
type SearchFieldView = ui.View

// NewSearchFieldView returns a new search field view
func NewSearchFieldView() *SearchFieldView {
	var view *SearchFieldView = ui.NewView("searchfield")
	return view
}

// InputView is structure for help view
type InputView = ui.View

// NewInputView returns a new help view
func NewInputView() *InputView {
	var view *InputView = ui.NewView("input")
	return view
}

// OpenSearch opens the search field
func (ct *Cointop) openSearch() error {
	ct.debuglog("openSearch()")
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
	ct.debuglog("cancelSearch()")
	ct.State.searchFieldVisible = false
	ct.ui.SetCursor(false)
	ct.SetActiveView(ct.Views.Table.Name())
	return nil
}

// DoSearch triggers the search and sets views
func (ct *Cointop) DoSearch() error {
	ct.debuglog("doSearch()")
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
	ct.debuglog("search()")
	q = strings.TrimSpace(strings.ToLower(q))
	idx := -1
	min := -1
	var hasprefixidx []int
	var hasprefixdist []int
	for i := range ct.State.allCoins {
		coin := ct.State.allCoins[i]
		name := strings.ToLower(coin.Name)
		symbol := strings.ToLower(coin.Symbol)
		// if query matches symbol, return immediately
		if symbol == q {
			ct.GoToGlobalIndex(i)
			return nil
		}
		// if query matches name, return immediately
		if name == q {
			ct.GoToGlobalIndex(i)
			return nil
		}
		// store index with the smallest levenshtein
		dist := levenshtein.DamerauLevenshteinDistance(name, q)
		if min == -1 || dist <= min {
			idx = i
			min = dist
		}
		// store index where query is substring to name
		if strings.HasPrefix(name, q) {
			if len(hasprefixdist) == 0 || dist < hasprefixdist[0] {
				hasprefixidx = append(hasprefixidx, i)
				hasprefixdist = append(hasprefixdist, dist)
			}
		}
	}
	// go to row if prefix match
	if len(hasprefixidx) > 0 && hasprefixidx[0] != -1 && min > 0 {
		ct.GoToGlobalIndex(hasprefixidx[0])
		return nil
	}
	// go to row if levenshtein distance is small enough
	if idx > -1 && min <= 6 {
		ct.GoToGlobalIndex(idx)
		return nil
	}
	return nil
}
