package api

import (
	types "github.com/miguelmota/cointop/pkg/api/types"
)

// Interface interface
type Interface interface {
	Ping() error
	GetAllCoinData(convert string) (map[string]types.Coin, error)
	GetCoinGraphData(coin string, start int64, end int64) (types.CoinGraph, error)
	GetGlobalMarketGraphData(start int64, end int64) (types.MarketGraph, error)
	GetGlobalMarketData(convert string) (types.GlobalMarketData, error)
	//GetCoinData(coin string) (types.Coin, error)
	//GetAltcoinMarketGraphData(start int64, end int64) (types.MarketGraph, error)
	//GetCoinPriceUSD(coin string) (float64, error)
	//GetCoinMarkets(coin string) ([]types.Market, error)
}
