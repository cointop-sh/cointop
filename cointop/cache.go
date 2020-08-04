package cointop

import (
	"fmt"
	"strings"
	"time"
)

// CacheKey returns cached value given key
func (ct *Cointop) CacheKey(key string) string {
	ct.debuglog("CacheKey()")
	return strings.ToLower(fmt.Sprintf("%s_%s", ct.apiChoice, key))
}

// CacheAllCoinsSlugMap writes the coins map to the memory and disk cache
func (ct *Cointop) CacheAllCoinsSlugMap() {
	ct.debuglog("CacheAllCoinsSlugMap()")
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
