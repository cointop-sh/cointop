package cointop

import (
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// Refresh triggers a force refresh of coin data
func (ct *Cointop) Refresh() error {
	log.Debug("Refresh()")
	go func() {
		<-ct.limiter
		ct.forceRefresh <- true
	}()
	return nil
}

// RefreshAll triggers a force refresh of all data
func (ct *Cointop) RefreshAll() error {
	log.Debug("RefreshAll()")
	ct.refreshMux.Lock()
	defer ct.refreshMux.Unlock()
	ct.setRefreshStatus()
	ct.cache.Delete("allCoinsSlugMap")
	ct.cache.Delete("market")
	go func() {
		ct.UpdateCoins()
		ct.UpdateTable()
		ct.UpdateChart()
	}()
	return nil
}

// SetRefreshStatus sets the refresh ticker
func (ct *Cointop) setRefreshStatus() {
	log.Debug("setRefreshStatus()")
	go func() {
		ct.loadingTicks("refreshing", 900)
		ct.RowChanged()
	}()
}

// LoadingTicks sets the loading ticking dots
func (ct *Cointop) loadingTicks(s string, t int) {
	log.Debug("loadingTicks()")
	interval := 150
	k := 0
	for i := 0; i < (t / interval); i++ {
		ct.UpdateStatusbar(s + strings.Repeat(".", k))
		time.Sleep(time.Duration(i*interval) * time.Millisecond)
		k = k + 1
		if k > 3 {
			k = 0
		}
	}
}

// intervalFetchData does a force refresh at every interval
func (ct *Cointop) intervalFetchData() {
	log.Debug("intervalFetchData()")
	go func() {
		for {
			select {
			case <-ct.forceRefresh:
				ct.RefreshAll()
			case <-ct.refreshTicker.C:
				ct.RefreshAll()
			}
		}
	}()
}
