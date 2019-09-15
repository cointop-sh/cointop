package cointop

import (
	"fmt"
	"strings"
	"time"

	"github.com/miguelmota/cointop/cointop/common/filecache"
)

func (ct *Cointop) cacheKey(key string) string {
	ct.debuglog("cacheKey()")
	return strings.ToLower(fmt.Sprintf("%s_%s", ct.apiChoice, key))
}

func (ct *Cointop) cacheAllCoinsSlugMap() {
	ct.debuglog("cacheAllCoinsSlugMap()")
	allCoinsSlugMap := make(map[string]*Coin)
	ct.State.allCoinsSlugMap.Range(func(key, value interface{}) bool {
		allCoinsSlugMap[key.(string)] = value.(*Coin)
		return true
	})

	// NOTE: do not override with empty data on startup
	if len(allCoinsSlugMap) != 0 {
		cachekey := ct.cacheKey("allCoinsSlugMap")
		ct.cache.Set(cachekey, allCoinsSlugMap, 10*time.Second)
		filecache.Set(cachekey, allCoinsSlugMap, 24*time.Hour)
	}
}
