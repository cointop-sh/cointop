package cointop

import (
	"strings"
	"time"
)

func (ct *Cointop) refresh() error {
	ct.debuglog("refresh()")
	go func() {
		<-ct.limiter
		ct.forceRefresh <- true
	}()
	return nil
}

func (ct *Cointop) refreshAll() error {
	ct.debuglog("refreshAll()")
	ct.refreshMux.Lock()
	defer ct.refreshMux.Unlock()
	ct.setRefreshStatus()
	ct.cache.Delete("allCoinsSlugMap")
	ct.cache.Delete("market")
	go func() {
		ct.updateCoins()
		ct.updateTable()
		ct.UpdateChart()
	}()
	return nil
}

func (ct *Cointop) setRefreshStatus() {
	ct.debuglog("setRefreshStatus()")
	go func() {
		ct.loadingTicks("refreshing", 900)
		ct.rowChanged()
	}()
}

func (ct *Cointop) loadingTicks(s string, t int) {
	ct.debuglog("loadingTicks()")
	interval := 150
	k := 0
	for i := 0; i < (t / interval); i++ {
		ct.updateStatusbar(s + strings.Repeat(".", k))
		time.Sleep(time.Duration(i*interval) * time.Millisecond)
		k = k + 1
		if k > 3 {
			k = 0
		}
	}
}

func (ct *Cointop) intervalFetchData() {
	ct.debuglog("intervalFetchData()")
	go func() {
		for {
			select {
			case <-ct.forceRefresh:
				ct.refreshAll()
			case <-ct.refreshTicker.C:
				ct.refreshAll()
			}
		}
	}()
}
