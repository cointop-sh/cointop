package cointop

import (
	"sync"
	"time"

	"github.com/cointop-sh/cointop/pkg/api/types"
	log "github.com/sirupsen/logrus"
)

var (
	coinslock      sync.Mutex
	updatecoinsmux sync.Mutex
)

// UpdateCoins updates coins view
func (ct *Cointop) UpdateCoins() error {
	log.Debug("UpdateCoins()")
	coinslock.Lock()
	defer coinslock.Unlock()
	cachekey := ct.CacheKey("allCoinsSlugMap")

	var err error
	var allCoinsSlugMap map[string]types.Coin
	cached, found := ct.cache.Get(cachekey)
	if found {
		// cache hit
		allCoinsSlugMap, _ = cached.(map[string]types.Coin)
		log.Debug("UpdateCoins() soft cache hit")
	}

	// cache miss or coin struct has been changed from the last time
	isCacheMissed := allCoinsSlugMap == nil
	currentCoinHash, _ := getStructHash(Coin{})
	isCoinStructHashChanged := currentCoinHash != ct.config.CoinStructHash
	if isCacheMissed || isCoinStructHashChanged {
		log.Debug("UpdateCoins() cache miss or coin struct has changed")
		ch := make(chan []types.Coin)
		err = ct.api.GetAllCoinData(ct.State.currencyConversion, ch)
		if err != nil {
			return err
		}

		for coins := range ch {
			go ct.processCoins(coins)
		}
	} else {
		ct.processCoinsMap(allCoinsSlugMap)
	}

	return nil
}

// UpdateCurrentPageCoins updates all the coins in the current page
func (ct *Cointop) UpdateCurrentPageCoins() error {
	log.Debugf("UpdateCurrentPageCoins(%d)", len(ct.State.coins))
	currentPageCoins := make([]string, len(ct.State.coins))
	for i, entry := range ct.State.coins {
		currentPageCoins[i] = entry.Name
	}

	coins, err := ct.api.GetCoinDataBatch(currentPageCoins, ct.State.currencyConversion)
	if err != nil {
		return err
	}
	go ct.processCoins(coins)
	return nil
}

// ProcessCoinsMap processes coins map
func (ct *Cointop) processCoinsMap(coinsMap map[string]types.Coin) {
	log.Debug("ProcessCoinsMap()")

	var coins []types.Coin
	for _, v := range coinsMap {
		coins = append(coins, v)
	}

	ct.processCoins(coins)
}

// ProcessCoins processes coins list
func (ct *Cointop) processCoins(coins []types.Coin) {
	log.Debug("ProcessCoins()")
	updatecoinsmux.Lock()
	defer updatecoinsmux.Unlock()

	ct.CacheAllCoinsSlugMap()

	for _, v := range coins {
		k := v.Name

		// Fix for https://github.com/cointop-sh/cointop/issues/59
		// some APIs returns rank 0 for new coins
		// or coins with low market cap data so we need to put them
		// at the end of the list
		if v.Rank == 0 {
			v.Rank = 10000
		}

		ilast, _ := ct.State.allCoinsSlugMap.Load(k)
		ct.State.allCoinsSlugMap.Store(k, &Coin{
			ID:               v.ID,
			Name:             v.Name,
			Symbol:           v.Symbol,
			Rank:             v.Rank,
			Price:            v.Price,
			Volume24H:        v.Volume24H,
			MarketCap:        v.MarketCap,
			AvailableSupply:  v.AvailableSupply,
			TotalSupply:      v.TotalSupply,
			PercentChange1H:  v.PercentChange1H,
			PercentChange24H: v.PercentChange24H,
			PercentChange7D:  v.PercentChange7D,
			PercentChange30D: v.PercentChange30D,
			PercentChange1Y:  v.PercentChange1Y,
			LastUpdated:      v.LastUpdated,
			Slug:             v.Slug,
		})
		if ilast != nil {
			last, _ := ilast.(*Coin)
			if last != nil {
				ivalue, _ := ct.State.allCoinsSlugMap.Load(k)
				l, _ := ivalue.(*Coin)
				l.Favorite = last.Favorite
				ct.State.allCoinsSlugMap.Store(k, l)
			}
		}
	}

	size := 0
	// NOTE: there's no Len method on sync.Map so need to manually count
	ct.State.allCoinsSlugMap.Range(func(key, value interface{}) bool {
		size++
		return true
	})

	if len(ct.State.allCoins) < size {
		var list []*Coin
		for _, v := range coins {
			k := v.Name
			icoin, _ := ct.State.allCoinsSlugMap.Load(k)
			coin, _ := icoin.(*Coin)
			list = append(list, coin)
		}
		ct.State.allCoins = append(ct.State.allCoins, list...)
	} else {
		// update list in place without changing order
		ct.State.allCoinsSlugMap.Range(func(key, value interface{}) bool {
			cm, _ := value.(*Coin)
			for k := range ct.State.allCoins {
				c := ct.State.allCoins[k]
				if c.ID == cm.ID {
					// TODO: improve this
					c.ID = cm.ID
					c.Name = cm.Name
					c.Symbol = cm.Symbol
					c.Rank = cm.Rank
					c.Price = cm.Price
					c.Volume24H = cm.Volume24H
					c.MarketCap = cm.MarketCap
					c.AvailableSupply = cm.AvailableSupply
					c.TotalSupply = cm.TotalSupply
					c.PercentChange1H = cm.PercentChange1H
					c.PercentChange24H = cm.PercentChange24H
					c.PercentChange7D = cm.PercentChange7D
					c.PercentChange30D = cm.PercentChange30D
					c.PercentChange1Y = cm.PercentChange1Y
					c.LastUpdated = cm.LastUpdated
					c.Favorite = cm.Favorite
					c.Slug = cm.Slug
				}
			}

			return true
		})
	}

	time.AfterFunc(10*time.Millisecond, func() {
		ct.Sort(ct.State.viewSorts[ct.State.selectedView], ct.State.coins, true)
		ct.UpdateTable()
	})
}

// GetListCount returns count of coins list
func (ct *Cointop) GetListCount() int {
	log.Debug("GetListCount()")
	if ct.IsFavoritesVisible() {
		return len(ct.State.favorites)
	} else if ct.IsPortfolioVisible() {
		return len(ct.State.portfolio.Entries)
	} else {
		return len(ct.State.allCoins)
	}
}
