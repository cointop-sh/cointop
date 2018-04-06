package cointop

import (
	"fmt"

	"github.com/miguelmota/cointop/pkg/color"
	"github.com/miguelmota/cointop/pkg/humanize"
	"github.com/miguelmota/cointop/pkg/pad"
)

func (ct *Cointop) updateMarketbar() error {
	maxX := ct.Width()
	market, err := ct.api.GetGlobalMarketData()
	if err != nil {
		return err
	}
	timeframe := "7 Day"
	chartname := ct.selectedCoinName()
	if chartname == "" {
		chartname = "Global"
	}
	ct.Update(func() {
		ct.marketview.Clear()
		fmt.Fprintln(
			ct.marketview,
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
