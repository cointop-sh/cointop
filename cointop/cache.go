package cointop

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// CacheKey returns cached value given key
func (ct *Cointop) CacheKey(key string) string {
	return strings.ToLower(fmt.Sprintf("%s_%s", ct.apiChoice, key))
}

// CompositeCacheKey returns a CacheKey for a coin (or globaldata)
func (ct *Cointop) CompositeCacheKey(symbol string, name string, convert string, chartRange string) string {
	keyname := symbol
	if name != "" {
		keyname += "-" + name
	}
	if convert != "" {
		keyname += "_" + convert
	}
	if chartRange != "" {
		keyname += "_" + strings.Replace(chartRange, " ", "", -1) // "All Time" contains space
	}
	return ct.CacheKey(keyname)
}

// CacheAllCoinsSlugMap writes the coins map to the memory and disk cache
func (ct *Cointop) CacheAllCoinsSlugMap() {
	log.Debug("CacheAllCoinsSlugMap()")
	allCoinsSlugMap := make(map[string]*Coin)
	ct.State.allCoinsSlugMap.Range(func(key, value interface{}) bool {
		allCoinsSlugMap[key.(string)] = value.(*Coin)
		return true
	})

	// NOTE: do not override with empty data on startup
	if len(allCoinsSlugMap) != 0 {
		cachekey := ct.CacheKey("allCoinsSlugMap")
		ct.cache.Set(cachekey, allCoinsSlugMap, 10*time.Second)
		if ct.filecache != nil {
			ct.filecache.Set(cachekey, allCoinsSlugMap, 24*time.Hour)
		}
	}
}
