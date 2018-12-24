package api

import (
	cmc "github.com/miguelmota/cointop/pkg/api/impl/coinmarketcap"
)

// NewCMC new CoinMarketCap API
func NewCMC() Interface {
	return cmc.New()
}

// NewCC new CryptoCompare API
func NewCC() {
	// TODO
}
