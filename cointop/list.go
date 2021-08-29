package cointop

import (
	"sync"
	"time"

	types "github.com/miguelmota/cointop/pkg/api/types"
)

var coinslock sync.Mutex
var updatecoinsmux sync.Mutex

// UpdateCoins updates coins view
func (ct *Cointop) UpdateCoins() error {
	ct.debuglog("updateCoins()")
	coinslock.Lock()
	defer coinslock.Unlock()
	cachekey := ct.CacheKey("allCoinsSlugMap")

	var err error
	var allCoinsSlugMap map[string]types.Coin
	cached, found := ct.cache.Get(cachekey)
	if found {
		// cache hit
		allCoinsSlugMap, _ = cached.(map[string]types.Coin)
		ct.debuglog("soft cache hit")
	}

	// cache miss
	if allCoinsSlugMap == nil {
		ct.debuglog("cache miss")
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

// ProcessCoinsMap processes coins map
func (ct *Cointop) processCoinsMap(coinsMap map[string]types.Coin) {
	ct.debuglog("processCoinsMap()")

	var coins []types.Coin
	for _, v := range coinsMap {
		coins = append(coins, v)
	}

	ct.processCoins(coins)
}

// ProcessCoins processes coins list
func (ct *Cointop) processCoins(coins []types.Coin) {
	ct.debuglog("processCoins()")
	updatecoinsmux.Lock()
	defer updatecoinsmux.Unlock()

	ct.CacheAllCoinsSlugMap()

	for _, v := range coins {
		k := v.Name

		// Fix for https://github.com/miguelmota/cointop/issues/59
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
		list := []*Coin{}
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
				}
			}

			return true
		})
	}

	time.AfterFunc(10*time.Millisecond, func() {
		ct.Sort(ct.State.sortBy, ct.State.sortDesc, ct.State.coins, true)
		ct.UpdateTable()
	})
}

// GetListCount returns count of coins list
func (ct *Cointop) GetListCount() int {
	ct.debuglog("getListCount()")
	if ct.IsFavoritesVisible() {
		return len(ct.State.favorites)
	} else if ct.IsPortfolioVisible() {
		return len(ct.State.portfolio.Entries)
	} else {
		return len(ct.State.allCoins)
	}
}
