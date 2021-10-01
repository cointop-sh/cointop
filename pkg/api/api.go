package api

import (
	cg "github.com/cointop-sh/cointop/pkg/api/impl/coingecko"
	cmc "github.com/cointop-sh/cointop/pkg/api/impl/coinmarketcap"
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
func NewCG(perPage, maxPages uint) Interface {
	return cg.NewCoinGecko(&cg.Config{
		PerPage:  perPage,
		MaxPages: maxPages,
	})
}
