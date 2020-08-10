package api

import (
	types "github.com/miguelmota/cointop/cointop/common/api/types"
)

// Interface interface
type Interface interface {
	Ping() error
	GetAllCoinData(convert string, ch chan []types.Coin) error
	GetCoinGraphData(convert string, symbol string, name string, start int64, end int64) (types.CoinGraph, error)
	GetGlobalMarketGraphData(convert string, start int64, end int64) (types.MarketGraph, error)
	GetGlobalMarketData(convert string) (types.GlobalMarketData, error)
	GetCoinData(name string, convert string) (types.Coin, error)
	GetCoinDataBatch(names []string, convert string) ([]types.Coin, error)
	//GetAltcoinMarketGraphData(start int64, end int64) (types.MarketGraph, error)
	//GetCoinPriceUSD(coin string) (float64, error)
	//GetCoinMarkets(coin string) ([]types.Market, error)
	CoinLink(name string) string
	SupportedCurrencies() []string
	Price(name string, convert string) (float64, error)
}
