package cointop

func (ct *Cointop) toggleShowFavorites() error {
	ct.filterByFavorites = !ct.filterByFavorites
	ct.updateTable()
	return nil
}
