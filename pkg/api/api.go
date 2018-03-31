package api

import (
	cmc "github.com/miguelmota/cointop/pkg/api/cmc"
)

// NewCMC new CoinMarketCap api
func NewCMC() Interface {
	return cmc.New()
}
