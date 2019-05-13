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
	cachekey := ct.cacheKey("allcoinsslugmap")

	var err error
	var allcoinsslugmap map[string]types.Coin
	cached, found := ct.cache.Get(cachekey)
	_ = cached
	if found {
		// cache hit
		allcoinsslugmap, _ = cached.(map[string]types.Coin)
		ct.debuglog("soft cache hit")
	}

	// cache miss
	if allcoinsslugmap == nil {
		ct.debuglog("cache miss")
		ch := make(chan []types.Coin)
		err = ct.api.GetAllCoinData(ct.currencyconversion, ch)
		if err != nil {
			return err
		}

		for coins := range ch {
			go ct.processCoins(coins)
			ct.cache.Set(cachekey, ct.allcoinsslugmap, 10*time.Second)
			filecache.Set(cachekey, ct.allcoinsslugmap, 24*time.Hour)
		}
	} else {
		ct.processCoinsMap(allcoinsslugmap)
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
	for _, v := range coins {
		k := v.Name
		last := ct.allcoinsslugmap[k]
		ct.allcoinsslugmap[k] = &Coin{
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
			ct.allcoinsslugmap[k].Favorite = last.Favorite
		}
	}
	if len(ct.allcoins) < len(ct.allcoinsslugmap) {
		list := []*Coin{}
		for _, v := range coins {
			k := v.Name
			coin := ct.allcoinsslugmap[k]
			list = append(list, coin)
		}
		ct.allcoins = append(ct.allcoins, list...)
	} else {
		// update list in place without changing order
		for i := range ct.allcoinsslugmap {
			cm := ct.allcoinsslugmap[i]
			for k := range ct.allcoins {
				c := ct.allcoins[k]
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
		ct.sort(ct.sortby, ct.sortdesc, ct.coins, true)
		ct.updateTable()
	})
}

func (ct *Cointop) getListCount() int {
	return len(ct.allCoins())
}
