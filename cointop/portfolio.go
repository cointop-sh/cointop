package cointop

func (ct *Cointop) togglePortfolio() error {
	ct.filterByFavorites = false
	ct.portfoliovisible = !ct.portfoliovisible
	ct.updateTable()
	return nil
}

func (ct *Cointop) toggleShowPortfolio() error {
	ct.filterByFavorites = false
	ct.portfoliovisible = true
	ct.updateTable()
	return nil
}
