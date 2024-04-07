package api

import (
	cg "github.com/cointop-sh/cointop/pkg/api/impl/coingecko"
	cmc "github.com/cointop-sh/cointop/pkg/api/impl/coinmarketcap"
)

type CoinGeckoConfig struct {
	PerPage   uint
	MaxPages  uint
	ApiKey    string
	ProApiKey string
}

// NewCMC new CoinMarketCap API
func NewCMC(apiKey string) Interface {
	return cmc.NewCMC(apiKey)
}

// NewCC new CryptoCompare API
func NewCC() {
	// TODO
}

// NewCG new CoinGecko API
func NewCG(config *CoinGeckoConfig) Interface {
	return cg.NewCoinGecko(&cg.Config{
		PerPage:   config.PerPage,
		MaxPages:  config.MaxPages,
		ApiKey:    config.ApiKey,
		ProApiKey: config.ProApiKey,
	})
}
