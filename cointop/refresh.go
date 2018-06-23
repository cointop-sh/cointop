package cointop

import (
	"strings"
	"time"
)

func (ct *Cointop) refresh() error {
	ct.forcerefresh <- true
	return nil
}

func (ct *Cointop) refreshAll() error {
	ct.refreshmux.Lock()
	defer ct.refreshmux.Unlock()
	ct.setRefreshStatus()
	ct.cache.Delete("allcoinsslugmap")
	ct.cache.Delete("allcoinssymbolmap")
	ct.cache.Delete("market")
	ct.updateCoins()
	ct.updateTable()
	ct.updateMarketbar()
	ct.updateChart()
	return nil
}

func (ct *Cointop) setRefreshStatus() {
	go func() {
		ct.loadingTicks("refreshing", 900)
		ct.updateStatusbar("")
		ct.rowChanged()
	}()
}

func (ct *Cointop) loadingTicks(s string, t int) {
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
