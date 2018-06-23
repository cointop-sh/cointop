package cointop

func (ct *Cointop) toggleFavorite() error {
	coin := ct.highlightedRowCoin()
	if coin == nil {
		return nil
	}
	_, ok := ct.favorites[coin.Name]
	if ok {
		delete(ct.favorites, coin.Name)
		coin.Favorite = false
	} else {
		ct.favorites[coin.Name] = true
		coin.Favorite = true
	}
	ct.updateTable()
	return nil
}

func (ct *Cointop) toggleShowFavorites() error {
	ct.filterByFavorites = !ct.filterByFavorites
	ct.updateTable()
	return nil
}
