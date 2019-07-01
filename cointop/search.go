package cointop

import (
	"regexp"
	"strings"

	"github.com/miguelmota/cointop/cointop/common/levenshtein"
)

func (ct *Cointop) openSearch() error {
	ct.State.searchFieldVisible = true
	ct.setActiveView(ct.Views.SearchField.Name)
	return nil
}

func (ct *Cointop) cancelSearch() error {
	ct.State.searchFieldVisible = false
	ct.setActiveView(ct.Views.Table.Name)
	return nil
}

func (ct *Cointop) doSearch() error {
	ct.Views.SearchField.Backing.Rewind()
	b := make([]byte, 100)
	n, err := ct.Views.SearchField.Backing.Read(b)

	// TODO: do this a better way (SoC)
	ct.State.filterByFavorites = false
	ct.State.portfolioVisible = false

	defer ct.setActiveView(ct.Views.Table.Name)
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
	return ct.search(q)
}

func (ct *Cointop) search(q string) error {
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
			ct.goToGlobalIndex(i)
			return nil
		}
		// if query matches name, return immediately
		if name == q {
			ct.goToGlobalIndex(i)
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
		ct.goToGlobalIndex(hasprefixidx[0])
		return nil
	}
	// go to row if levenshtein distance is small enough
	if idx > -1 && min <= 6 {
		ct.goToGlobalIndex(idx)
		return nil
	}
	return nil
}
