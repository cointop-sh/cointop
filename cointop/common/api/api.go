package api

import (
	cmc "github.com/miguelmota/cointop/cointop/common/api/impl/coinmarketcap"
)

// NewCMC new CoinMarketCap API
func NewCMC() Interface {
	return cmc.New()
}

// NewCC new CryptoCompare API
func NewCC() {
	// TODO
}

// NewCG new CoinGecko API
func NewCG() {
	// TODO
}
