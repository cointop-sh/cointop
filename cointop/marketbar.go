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
	ct.Update(func() {
		ct.marketview.Clear()
		fmt.Fprintln(ct.marketview, color.White(pad.Right(fmt.Sprintf("%s Chart: %s • Total Market Cap: %s • 24H Volume: %s • BTC Dominance: %.2f%% • Active Currencies: %s • Active Markets: %s", color.Cyan("cointop"), timeframe, humanize.Commaf(market.TotalMarketCapUSD), humanize.Commaf(market.Total24HVolumeUSD), market.BitcoinPercentageOfMarketCap, humanize.Commaf(float64(market.ActiveCurrencies)), humanize.Commaf(float64(market.ActiveMarkets))), maxX, " ")))
	})
	return nil
}
