package cointop

import (
	"fmt"
	"math"
	"time"

	types "github.com/miguelmota/cointop/cointop/common/api/types"
	"github.com/miguelmota/cointop/cointop/common/color"
	"github.com/miguelmota/cointop/cointop/common/filecache"
	"github.com/miguelmota/cointop/cointop/common/humanize"
	"github.com/miguelmota/cointop/cointop/common/pad"
)

func (ct *Cointop) updateMarketbar() error {
	if ct.marketbarview == nil {
		return nil
	}

	maxX := ct.width()
	logo := "❯❯❯cointop"
	if ct.colorschemename == "cointop" {
		logo = fmt.Sprintf("%s%s%s%s", color.Green("❯"), color.Cyan("❯"), color.Green("❯"), color.Cyan("cointop"))
	}
	var content string

	if ct.portfoliovisible {
		total := ct.getPortfolioTotal()
		totalstr := humanize.Commaf(total)
		if !(ct.currencyconversion == "BTC" || ct.currencyconversion == "ETH" || total < 1) {
			total = math.Round(total*1e2) / 1e2
			totalstr = humanize.Commaf2(total)
		}

		timeframe := ct.selectedchartrange
		chartname := ct.selectedCoinName()
		var charttitle string
		if chartname == "" {
			chartname = "Portfolio"
			charttitle = ct.colorscheme.MarketBarLabelActive(chartname)
		} else {
			charttitle = fmt.Sprintf("Portfolio - %s", ct.colorscheme.MarketBarLabelActive(chartname))
		}

		var percentChange24H float64
		for _, p := range ct.getPortfolioSlice() {
			n := ((p.Balance / total) * p.PercentChange24H)
			if math.IsNaN(n) {
				continue
			}
			percentChange24H += n
		}

		color24h := ct.colorscheme.MarketbarSprintf()
		arrow := ""
		if percentChange24H > 0 {
			color24h = ct.colorscheme.MarketbarChangeUpSprintf()
			arrow = "▲"
		}
		if percentChange24H < 0 {
			color24h = ct.colorscheme.MarketbarChangeDownSprintf()
			arrow = "▼"
		}

		chartInfo := ""
		if !ct.hideChart {
			chartInfo = fmt.Sprintf(
				"[ Chart: %s %s ] ",
				charttitle,
				timeframe,
			)
		}

		content = fmt.Sprintf(
			"%sTotal Portfolio Value: %s • 24H: %s",
			chartInfo,
			ct.colorscheme.MarketBarLabelActive(fmt.Sprintf("%s%s", ct.currencySymbol(), totalstr)),
			color24h(fmt.Sprintf("%.2f%%%s", percentChange24H, arrow)),
		)
	} else {
		var market types.GlobalMarketData
		var err error
		cachekey := ct.cacheKey("market")
		cached, found := ct.cache.Get(cachekey)

		if found {
			// cache hit
			var ok bool
			market, ok = cached.(types.GlobalMarketData)
			if ok {
				ct.debuglog("soft cache hit")
			}
		}

		if market.TotalMarketCapUSD == 0 {
			market, err = ct.api.GetGlobalMarketData(ct.currencyconversion)
			if err != nil {
				filecache.Get(cachekey, &market)
			}

			ct.cache.Set(cachekey, market, 10*time.Second)
			go func() {
				filecache.Set(cachekey, market, 24*time.Hour)
			}()
		}

		timeframe := ct.selectedchartrange
		chartname := ct.selectedCoinName()
		if chartname == "" {
			chartname = "Global"
		}

		chartInfo := ""
		if !ct.hideChart {
			chartInfo = fmt.Sprintf(
				"[ Chart: %s %s ] ",
				ct.colorscheme.MarketBarLabelActive(chartname),
				timeframe,
			)
		}

		content = fmt.Sprintf(
			"%sGlobal ▶ Market Cap: %s • 24H Volume: %s • BTC Dominance: %.2f%%",
			chartInfo,
			fmt.Sprintf("%s%s", ct.currencySymbol(), humanize.Commaf(market.TotalMarketCapUSD)),
			fmt.Sprintf("%s%s", ct.currencySymbol(), humanize.Commaf(market.Total24HVolumeUSD)),
			market.BitcoinPercentageOfMarketCap,
		)
	}

	content = fmt.Sprintf("%s %s", logo, content)
	content = pad.Right(content, maxX, " ")
	content = ct.colorscheme.Marketbar(content)

	ct.update(func() {
		if ct.marketbarview == nil {
			return
		}

		ct.marketbarview.Clear()
		fmt.Fprintln(ct.marketbarview, content)
	})

	return nil
}
