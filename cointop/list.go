package cointop

import (
	"sync"

	types "github.com/miguelmota/cointop/cointop/common/api/types"
)

var coinslock sync.Mutex

func (ct *Cointop) updateCoins() error {
	coinslock.Lock()
	defer coinslock.Unlock()
	cachekey := "allcoinsslugmap"

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
		ch, err := ct.api.GetAllCoinData(ct.currencyconversion)
		if err != nil {
			return err
		}

		for {
			coins, ok := <-ch
			if !ok {
				break
			}
			ct.updateCoinsMap(coins, true)
			ct.updateTable()
		}

		/*
			ct.cache.Set(cachekey, allcoinsslugmap, 10*time.Second)
			go func() {
				filecache.Set(cachekey, allcoinsslugmap, 24*time.Hour)
			}()
		*/
	} else {
		ct.updateCoinsMap(allcoinsslugmap, false)
	}

	return nil
}

func (ct *Cointop) updateCoinsMap(allcoinsslugmap map[string]types.Coin, b bool) {
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

		if b {
			ct.allcoins = append(ct.allcoins, ct.allcoinsslugmap[k])
		}
	}

	//if len(ct.allcoins) == 0 {
	if b {
		//ct.sort(ct.sortby, ct.sortdesc, ct.allcoins)
	}
	if !b {
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
}
