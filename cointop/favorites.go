package cointop

func (ct *Cointop) toggleShowFavorites() error {
	for _, v := range ct.allcoinsmap {
		if ct.favorites[v.Symbol] {
			v.Favorite = true
		}
	}

	ct.filterByFavorites = !ct.filterByFavorites
	ct.sort(ct.sortby, ct.sortdesc, ct.coins)
	//ct.updateCoins()
	/*
		ct.update(func() {
			ct.tableview.Clear()
			ct.updateTable()
		})
		ct.rowChanged()
	*/
	return nil
}
