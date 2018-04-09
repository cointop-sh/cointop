package cointop

import (
	"regexp"
	"strings"
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
	for i := range ct.allcoins {
		coin := ct.allcoins[i]
		if strings.ToLower(coin.Name) == q || strings.ToLower(coin.Symbol) == q {
			ct.goToGlobalIndex(i)
			return nil
		}
	}
	return nil
}

func (ct *Cointop) goToGlobalIndex(idx int) error {
	perpage := ct.getTotalPerPage()
	atpage := idx / perpage
	ct.setPage(atpage)
	rowIndex := (idx % perpage)
	ct.highlightRow(rowIndex)
	ct.updateTable()
	ct.rowChanged()
	return nil
}
