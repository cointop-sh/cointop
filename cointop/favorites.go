package cointop

func (ct *Cointop) toggleFavorite() error {
	ct.State.portfolioVisible = false
	coin := ct.HighlightedRowCoin()
	if coin == nil {
		return nil
	}
	_, ok := ct.State.favorites[coin.Name]
	if ok {
		delete(ct.State.favorites, coin.Name)
		coin.Favorite = false
	} else {
		ct.State.favorites[coin.Name] = true
		coin.Favorite = true
	}
	go ct.updateTable()
	return nil
}

func (ct *Cointop) toggleShowFavorites() error {
	ct.State.portfolioVisible = false
	ct.State.filterByFavorites = !ct.State.filterByFavorites
	go ct.updateTable()
	return nil
}
