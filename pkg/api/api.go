package api

import (
	cg "github.com/miguelmota/cointop/pkg/api/impl/coingecko"
	cmc "github.com/miguelmota/cointop/pkg/api/impl/coinmarketcap"
)

// NewCMC new CoinMarketCap API
func NewCMC(apiKey string) Interface {
	return cmc.NewCMC(apiKey)
}

// NewCC new CryptoCompare API
func NewCC() {
	// TODO
}

// NewCG new CoinGecko API
func NewCG() Interface {
	return cg.NewCoinGecko()
}
