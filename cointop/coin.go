package cointop

import log "github.com/sirupsen/logrus"

// Coin is the row structure
type Coin struct {
	ID               string
	Name             string
	Slug             string
	Symbol           string
	Rank             int
	Price            float64
	Volume24H        float64
	MarketCap        float64
	AvailableSupply  float64
	TotalSupply      float64
	PercentChange1H  float64
	PercentChange24H float64
	PercentChange7D  float64
	PercentChange30D float64
	PercentChange1Y  float64
	LastUpdated      string
	// for favorites
	Favorite bool
	// for portfolio
	Holdings float64
	Balance  float64
}

// AllCoins returns a slice of all the coins
func (ct *Cointop) AllCoins() []*Coin {
	log.Debug("AllCoins()")
	if ct.IsFavoritesVisible() {
		var list []*Coin
		for i := range ct.State.allCoins {
			coin := ct.State.allCoins[i]
			if coin.Favorite {
				list = append(list, coin)
			}
		}
		return list
	}

	if ct.IsPortfolioVisible() {
		var list []*Coin
		for i := range ct.State.allCoins {
			coin := ct.State.allCoins[i]
			if ct.PortfolioEntryExists(coin) {
				list = append(list, coin)
			}
		}
		return list
	}

	return ct.State.allCoins
}

// CoinBySymbol returns the coin struct given the symbol
func (ct *Cointop) CoinBySymbol(symbol string) *Coin {
	log.Debug("CoinBySymbol()")
	for i := range ct.State.allCoins {
		coin := ct.State.allCoins[i]
		if coin.Symbol == symbol {
			return coin
		}
	}
	return nil
}

// CoinByName returns the coin struct given the name
func (ct *Cointop) CoinByName(name string) *Coin {
	log.Debug("CoinByName()")
	for i := range ct.State.allCoins {
		coin := ct.State.allCoins[i]
		if coin.Name == name {
			return coin
		}
	}
	return nil
}

// CoinByID returns the coin struct given the ID
func (ct *Cointop) CoinByID(id string) *Coin {
	log.Debug("CoinByID()")
	for i := range ct.State.allCoins {
		coin := ct.State.allCoins[i]
		if coin.ID == id {
			return coin
		}
	}
	return nil
}
