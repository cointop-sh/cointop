package cointop

import (
	"sync"
	"time"

	types "github.com/miguelmota/cointop/cointop/common/api/types"
	"github.com/miguelmota/cointop/cointop/common/filecache"
)

var coinslock sync.Mutex
var updatecoinsmux sync.Mutex

func (ct *Cointop) updateCoins() error {
	coinslock.Lock()
	defer coinslock.Unlock()
	cachekey := ct.cacheKey("allCoinsSlugMap")

	var err error
	var allCoinsSlugMap map[string]types.Coin
	cached, found := ct.cache.Get(cachekey)
	_ = cached
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

func (ct *Cointop) processCoinsMap(coinsMap map[string]types.Coin) {
	var coins []types.Coin
	for _, v := range coinsMap {
		coins = append(coins, v)
	}

	ct.processCoins(coins)
}

func (ct *Cointop) processCoins(coins []types.Coin) {
	updatecoinsmux.Lock()
	defer updatecoinsmux.Unlock()

	cachekey := ct.cacheKey("allCoinsSlugMap")
	ct.cache.Set(cachekey, ct.State.allCoinsSlugMap, 10*time.Second)
	filecache.Set(cachekey, ct.State.allCoinsSlugMap, 24*time.Hour)

	for _, v := range coins {
		k := v.Name
		last := ct.State.allCoinsSlugMap[k]
		ct.State.allCoinsSlugMap[k] = &Coin{
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
			LastUpdated:      v.LastUpdated,
		}
		if last != nil {
			ct.State.allCoinsSlugMap[k].Favorite = last.Favorite
		}
	}
	if len(ct.State.allCoins) < len(ct.State.allCoinsSlugMap) {
		list := []*Coin{}
		for _, v := range coins {
			k := v.Name
			coin := ct.State.allCoinsSlugMap[k]
			list = append(list, coin)
		}
		ct.State.allCoins = append(ct.State.allCoins, list...)
	} else {
		// update list in place without changing order
		for i := range ct.State.allCoinsSlugMap {
			cm := ct.State.allCoinsSlugMap[i]
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
					c.LastUpdated = cm.LastUpdated
					c.Favorite = cm.Favorite
				}
			}
		}
	}

	time.AfterFunc(10*time.Millisecond, func() {
		ct.sort(ct.State.sortBy, ct.State.sortDesc, ct.State.coins, true)
		ct.updateTable()
	})
}

func (ct *Cointop) getListCount() int {
	return len(ct.allCoins())
}
