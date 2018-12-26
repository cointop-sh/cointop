package cointop

import (
	"fmt"
	"math"
	"time"

	types "github.com/miguelmota/cointop/pkg/api/types"
	"github.com/miguelmota/cointop/pkg/color"
	"github.com/miguelmota/cointop/pkg/fcache"
	"github.com/miguelmota/cointop/pkg/humanize"
	"github.com/miguelmota/cointop/pkg/pad"
)

func (ct *Cointop) updateMarketbar() error {
	maxX := ct.width()
	logo := fmt.Sprintf("%s%s%s%s", color.Green("❯"), color.Cyan("❯"), color.Green("❯"), color.Cyan("cointop"))
	var content string

	if ct.portfoliovisible {
		total := ct.getPortfolioTotal()
		if !(ct.currencyconversion == "BTC" || ct.currencyconversion == "ETH" || total < 1) {
			total = math.Round(total*1e2) / 1e2
		}

		timeframe := ct.selectedchartrange
		chartname := ct.selectedCoinName()
		var charttitle string
		if chartname == "" {
			chartname = "Portfolio"
			charttitle = color.Cyan(chartname)
		} else {
			charttitle = fmt.Sprintf("Portfolio - %s", color.Cyan(chartname))
		}

		content = fmt.Sprintf(
			"[ Chart: %s %s ] Current Portfolio Value: %s%s",
			charttitle,
			timeframe,
			ct.currencySymbol(),
			humanize.Commaf(total),
		)
	} else {
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
			market, err = ct.api.GetGlobalMarketData(ct.currencyconversion)
			if err != nil {
				return err
			}

			ct.cache.Set(cachekey, market, 10*time.Second)
			go func() {
				_ = fcache.Set(cachekey, market, 24*time.Hour)
			}()
		}

		timeframe := ct.selectedchartrange
		chartname := ct.selectedCoinName()
		if chartname == "" {
			chartname = "Global"
		}

		content = fmt.Sprintf(
			"[ Chart: %s %s ] Global ▶ Market Cap: $%s • 24H Volume: $%s • BTC Dominance: %.2f%% • Active Currencies: %s • Active Markets: %s",
			color.Cyan(chartname),
			timeframe,
			humanize.Commaf(market.TotalMarketCapUSD),
			humanize.Commaf(market.Total24HVolumeUSD),
			market.BitcoinPercentageOfMarketCap,
			humanize.Commaf(float64(market.ActiveCurrencies)),
			humanize.Commaf(float64(market.ActiveMarkets)),
		)
	}

	content = fmt.Sprintf("%s %s", logo, content)
	content = pad.Right(content, maxX, " ")

	ct.update(func() {
		ct.marketbarview.Clear()
		fmt.Fprintln(ct.marketbarview, content)
	})

	return nil
}
