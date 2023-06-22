package cointop

import (
	"sort"

	log "github.com/sirupsen/logrus"
)

// GetFavoritesTableHeaders returns the favorites table headers
func (ct *Cointop) GetFavoritesTableHeaders() []string {
	return ct.State.favoritesTableColumns
}

// ToggleFavorite toggles coin as favorite
func (ct *Cointop) ToggleFavorite() error {
	log.Debug("ToggleFavorite()")
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

	if err := ct.Save(); err != nil {
		return err
	}

	go ct.UpdateTable()

	return nil
}

// ToggleFavorites toggles the favorites view
func (ct *Cointop) ToggleFavorites() error {
	log.Debug("ToggleFavorites()")
	ct.ToggleSelectedView(FavoritesView)
	go ct.UpdateTable()
	return nil
}

// ToggleShowFavorites shows the favorites view
func (ct *Cointop) ToggleShowFavorites() error {
	log.Debug("ToggleShowFavorites()")
	ct.ToggleSelectedView(FavoritesView)
	go ct.UpdateTable()
	return nil
}

// GetFavoritesSlice returns coin favorites as slice
func (ct *Cointop) GetFavoritesSlice() []*Coin {
	log.Debug("GetFavoritesSlice()")
	var sliced []*Coin
	for i := range ct.State.allCoins {
		coin := ct.State.allCoins[i]
		if coin.Favorite {
			sliced = append(sliced, coin)
		}
	}

	sort.SliceStable(sliced, func(i, j int) bool {
		return sliced[i].MarketCap > sliced[j].MarketCap
	})

	for i, coin := range sliced {
		coin.Rank = i + 1
	}

	return sliced
}

// IsFavoritesVisible returns true if favorites view is visible
func (ct *Cointop) IsFavoritesVisible() bool {
	return ct.State.selectedView == FavoritesView
}

// FavoritesLen returns the number of favorite coins
func (ct *Cointop) FavoritesLen() int {
	return len(ct.GetFavoritesSlice())
}
