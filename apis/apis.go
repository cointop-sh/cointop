package apis

import (
	cmc "github.com/miguelmota/cointop/apis/cmc"
)

// NewCMC new CoinMarketCap api
func NewCMC() Interface {
	return cmc.New()
}
