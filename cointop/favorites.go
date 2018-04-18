package cointop

func (ct *Cointop) toggleFavorite() error {
	coin := ct.highlightedRowCoin()
	if coin == nil {
		return nil
	}
	_, ok := ct.favorites[coin.Symbol]
	if ok {
		delete(ct.favorites, coin.Symbol)
		coin.Favorite = false
	} else {
		ct.favorites[coin.Symbol] = true
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
