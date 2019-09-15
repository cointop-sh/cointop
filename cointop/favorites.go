package cointop

import "sort"

func (ct *Cointop) toggleFavorite() error {
	ct.debuglog("toggleFavorite()")
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

	if err := ct.save(); err != nil {
		return err
	}

	go ct.updateTable()

	return nil
}

func (ct *Cointop) toggleShowFavorites() error {
	ct.debuglog("toggleShowFavorites()")
	ct.State.portfolioVisible = false
	ct.State.filterByFavorites = !ct.State.filterByFavorites
	go ct.updateTable()
	return nil
}

func (ct *Cointop) getFavoritesSlice() []*Coin {
	ct.debuglog("getFavoritesSlice()")
	sliced := []*Coin{}
	for i := range ct.State.allCoins {
		coin := ct.State.allCoins[i]
		if coin.Favorite {
			sliced = append(sliced, coin)
		}
	}

	sort.Slice(sliced, func(i, j int) bool {
		return sliced[i].MarketCap > sliced[j].MarketCap
	})

	for i, coin := range sliced {
		coin.Rank = i + 1
	}

	return sliced
}
