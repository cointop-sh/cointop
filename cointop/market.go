package cointop

import (
	"fmt"

	humanize "github.com/dustin/go-humanize"
	"github.com/miguelmota/cointop/pkg/pad"
)

func (ct *Cointop) updateMarket() error {
	maxX, _ := ct.g.Size()
	market, err := ct.api.GetGlobalMarketData()
	if err != nil {
		return err
	}
	fmt.Fprintln(ct.marketview, pad.Right(fmt.Sprintf("%10.stotal market cap: %s", "", humanize.Commaf(market.TotalMarketCapUSD)), maxX, " "))
	return nil
}
