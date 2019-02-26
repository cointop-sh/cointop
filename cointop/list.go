package cointop

import (
	"sync"
	"time"

	types "github.com/miguelmota/cointop/cointop/common/api/types"
	"github.com/miguelmota/cointop/cointop/common/filecache"
)

var coinslock sync.Mutex

func (ct *Cointop) updateCoins() error {
	coinslock.Lock()
	defer coinslock.Unlock()
	list := []*coin{}
	cachekey := "allcoinsslugmap"

	var err error
	var allcoinsslugmap map[string]types.Coin
	cached, found := ct.cache.Get(cachekey)
	if found {
		// cache hit
		allcoinsslugmap, _ = cached.(map[string]types.Coin)
		ct.debuglog("soft cache hit")
	}

	// cache miss
	if allcoinsslugmap == nil {
		ct.debuglog("cache miss")
		allcoinsslugmap, err = ct.api.GetAllCoinData(ct.currencyconversion)
		if err != nil {
			return err
		}
		ct.cache.Set(cachekey, allcoinsslugmap, 10*time.Second)
		go func() {
			filecache.Set(cachekey, allcoinsslugmap, 24*time.Hour)
		}()
	}

	if len(ct.allcoinsslugmap) == 0 {
		ct.allcoinsslugmap = map[string]*coin{}
	}
	for k, v := range allcoinsslugmap {
		last := ct.allcoinsslugmap[k]
		ct.allcoinsslugmap[k] = &coin{
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
	if len(ct.allcoins) == 0 {
		for i := range ct.allcoinsslugmap {
			coin := ct.allcoinsslugmap[i]
			list = append(list, coin)
		}
		ct.allcoins = list
		ct.sort(ct.sortby, ct.sortdesc, ct.allcoins)
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
	return nil
}
