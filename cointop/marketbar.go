package cointop

import (
	"fmt"
	"time"

	types "github.com/miguelmota/cointop/pkg/api/types"
	"github.com/miguelmota/cointop/pkg/color"
	"github.com/miguelmota/cointop/pkg/fcache"
	"github.com/miguelmota/cointop/pkg/humanize"
	"github.com/miguelmota/cointop/pkg/pad"
)

func (ct *Cointop) updateMarketbar() error {
	maxX := ct.width()

	var market types.GlobalMarketData
	var err error
	cachekey := "market"
	cached, found := ct.cache.Get(cachekey)
	if found {
		// cache hit
		var ok bool
		market, ok = cached.(types.GlobalMarketData)
		if ok {
			ct.debuglog("soft cache hit")
		}
	} else {
		market, err = ct.api.GetGlobalMarketData()
		if err != nil {
			return err
		}

		ct.cache.Set(cachekey, market, 10*time.Second)
		go func() {
			_ = fcache.Set(cachekey, market, 24*time.Hour)
		}()
	}

	timeframe := "7 Day"
	chartname := ct.selectedCoinName()
	if chartname == "" {
		chartname = "Global"
	}
	ct.update(func() {
		ct.marketbarview.Clear()
		fmt.Fprintln(
			ct.marketbarview,
			pad.Right(
				fmt.Sprintf(
					"%s [ Chart: %s %s ] Global ▶ Market Cap: %s • 24H Volume: %s • BTC Dominance: %.2f%% • Active Currencies: %s • Active Markets: %s",
					fmt.Sprintf("%s%s%s%s", color.Green("❯"), color.Cyan("❯"), color.Green("❯"), color.Cyan("cointop")),
					color.Cyan(chartname),
					timeframe,
					humanize.Commaf(market.TotalMarketCapUSD),
					humanize.Commaf(market.Total24HVolumeUSD),
					market.BitcoinPercentageOfMarketCap,
					humanize.Commaf(float64(market.ActiveCurrencies)),
					humanize.Commaf(float64(market.ActiveMarkets)),
				),
				maxX,
				" ",
			),
		)
	})
	return nil
}
