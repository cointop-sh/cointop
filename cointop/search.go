package cointop

import (
	"regexp"
	"strings"

	"github.com/miguelmota/cointop/pkg/levenshtein"
)

func (ct *Cointop) openSearch() error {
	ct.setActiveView("searchfield")
	return nil
}

func (ct *Cointop) cancelSearch() error {
	ct.setActiveView("table")
	return nil
}

func (ct *Cointop) doSearch() error {
	ct.searchfield.Rewind()
	b := make([]byte, 100)
	n, err := ct.searchfield.Read(b)
	defer ct.setActiveView("table")
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
	for i := range ct.allcoins {
		coin := ct.allcoins[i]
		name := strings.ToLower(coin.Name)
		symbol := strings.ToLower(coin.Symbol)
		if symbol == q {
			idx = i
			min = 0
			break
		}
		dist := levenshtein.Distance(name, q)
		if min == -1 {
			min = dist
		}
		if dist <= min {
			idx = i
			min = dist
		}
	}

	if idx > -1 && min <= 6 {
		ct.goToGlobalIndex(idx)
	}
	return nil
}

func (ct *Cointop) goToGlobalIndex(idx int) error {
	perpage := ct.totalPerPage()
	atpage := idx / perpage
	ct.setPage(atpage)
	rowIndex := (idx % perpage)
	ct.highlightRow(rowIndex)
	ct.updateTable()
	ct.rowChanged()
	return nil
}
