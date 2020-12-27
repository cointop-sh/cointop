package cointop

import (
	"fmt"
	"log"
	"time"

	"github.com/miguelmota/cointop/pkg/humanize"
	"github.com/miguelmota/cointop/pkg/notifier"
	"github.com/miguelmota/cointop/pkg/table"
)

// GetAlertsTableHeaders returns the alerts table headers
func (ct *Cointop) GetAlertsTableHeaders() []string {
	return []string{
		"name",
		"symbol",
		"targetprice", //>600
		"price",
		"frequency",
	}
}

var gt = ">"
var gte = "≥"
var lte = "≤"
var lt = "<"
var eq = "="

// PriceAlertDirectionsMap is map of valid price alert direction symbols
var PriceAlertDirectionsMap = map[string]bool{
	">":  true,
	"<":  true,
	">=": true,
	"<=": true,
	"=":  true,
}

// PriceAlertFrequencyMap is map of valid price alert frequency values
var PriceAlertFrequencyMap = map[string]bool{
	"once":        true,
	"reoccurring": true,
}

// GetAlertsTable returns the table for displaying alerts
func (ct *Cointop) GetAlertsTable() *table.Table {
	maxX := ct.width()
	t := table.NewTable().SetWidth(maxX)

	for _, entry := range ct.State.priceAlerts.Entries {
		ifc, ok := ct.State.allCoinsSlugMap.Load(entry.CoinName)
		if !ok {
			continue
		}
		coin, ok := ifc.(*Coin)
		if !ok {
			continue
		}
		name := TruncateString(entry.CoinName, 20)
		symbol := TruncateString(coin.Symbol, 6)
		namecolor := ct.colorscheme.TableRow
		frequency := entry.Frequency
		targetPrice := fmt.Sprintf("%s%v", gte, entry.TargetPrice)

		t.AddRowCells(
			&table.RowCell{
				LeftMargin: 1,
				Width:      22,
				LeftAlign:  true,
				Color:      namecolor,
				Text:       name,
			},
			&table.RowCell{
				LeftMargin: 1,
				Width:      10,
				LeftAlign:  true,
				Color:      ct.colorscheme.TableRow,
				Text:       symbol,
			},
			&table.RowCell{
				LeftMargin: 1,
				Width:      16,
				LeftAlign:  false,
				Color:      ct.colorscheme.TableRow,
				Text:       targetPrice,
			},
			&table.RowCell{
				LeftMargin: 1,
				Width:      11,
				LeftAlign:  false,
				Color:      ct.colorscheme.TableRow,
				Text:       humanize.Commaf(coin.Price),
			},
			&table.RowCell{
				LeftMargin: 2,
				Width:      11,
				LeftAlign:  true,
				Color:      ct.colorscheme.TableRow,
				Text:       frequency,
			},
		)
	}

	return t
}

// ToggleAlerts toggles the alerts view
func (ct *Cointop) ToggleAlerts() error {
	ct.debuglog("toggleAlerts()")
	ct.ToggleSelectedView(AlertsView)
	go ct.UpdateTable()
	return nil
}

// IsAlertsVisible returns true if alerts view is visible
func (ct *Cointop) IsAlertsVisible() bool {
	return ct.State.selectedView == AlertsView
}

// PriceAlertWatcher starts the price alert watcher
func (ct *Cointop) PriceAlertWatcher() {
	ct.debuglog("priceAlertWatcher()")
	alerts := ct.State.priceAlerts.Entries
	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-ticker.C:
			for _, alert := range alerts {
				err := ct.CheckPriceAlert(alert)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

// CheckPriceAlert checks the price alert
func (ct *Cointop) CheckPriceAlert(alert *PriceAlert) error {
	ct.debuglog("checkPriceAlert()")
	if alert.Expired {
		return nil
	}

	cacheKey := ct.CacheKey("priceAlerts")
	var cachedEntries []*PriceAlert
	ct.filecache.Get(cacheKey, &cachedEntries)
	for _, cachedEntry := range cachedEntries {
		if cachedEntry.ID == alert.ID {
			alert.Expired = cachedEntry.Expired
			if alert.Expired {
				return nil
			}
		}
	}

	coinIfc, _ := ct.State.allCoinsSlugMap.Load(alert.CoinName)
	coin, ok := coinIfc.(*Coin)
	if !ok {
		return nil
	}
	var msg string
	title := "Cointop Alert"
	priceStr := fmt.Sprintf("$%s", humanize.Commaf(alert.TargetPrice))
	if alert.Direction == ">" {
		if coin.Price > alert.TargetPrice {
			msg = fmt.Sprintf("%s price is greater than %v", alert.CoinName, priceStr)
		}
	} else if alert.Direction == ">=" {
		if coin.Price >= alert.TargetPrice {
			msg = fmt.Sprintf("%s price is greater than or equal to %v", alert.CoinName, priceStr)
		}
	} else if alert.Direction == "<" {
		if coin.Price < alert.TargetPrice {
			msg = fmt.Sprintf("%s price is less than %v", alert.CoinName, priceStr)
		}
	} else if alert.Direction == "<=" {
		if coin.Price <= alert.TargetPrice {
			msg = fmt.Sprintf("%s price is less than or equal to %v", alert.CoinName, priceStr)
		}
	} else if alert.Direction == "=" {
		if coin.Price == alert.TargetPrice {
			msg = fmt.Sprintf("%s price is equal to %v", alert.CoinName, priceStr)
		}
	}

	if msg != "" {
		if ct.State.priceAlerts.SoundEnabled {
			notifier.NotifyWithSound(title, msg)
		} else {
			notifier.Notify(title, msg)
		}

		alert.Expired = true
		ct.filecache.Set(cacheKey, ct.State.priceAlerts.Entries, 87600*time.Hour)
	}
	return nil
}
