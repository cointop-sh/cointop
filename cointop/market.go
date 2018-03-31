package cointop

import (
	"fmt"

	humanize "github.com/dustin/go-humanize"
	"github.com/miguelmota/cointop/pkg/color"
	"github.com/miguelmota/cointop/pkg/pad"
)

func (ct *Cointop) updateMarket() error {
	maxX, _ := ct.g.Size()
	market, err := ct.api.GetGlobalMarketData()
	if err != nil {
		return err
	}
	fmt.Fprintln(ct.marketview, pad.Right(fmt.Sprintf("%s   Total Market Cap: %s • 24H Volume: %s • BTC Dominance: %.2f%% • Active Currencies: %s • Active Assets: %s • Active Markets: %s", color.Cyan("cointop"), color.WhiteBold(humanize.Commaf(market.TotalMarketCapUSD)), humanize.Commaf(market.Total24HVolumeUSD), market.BitcoinPercentageOfMarketCap, humanize.Comma(int64(market.ActiveCurrencies)), humanize.Comma(int64(market.ActiveAssets)), humanize.Comma(int64(market.ActiveMarkets))), maxX, " "))
	return nil
}
