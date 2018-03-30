package apis

import (
	types "github.com/miguelmota/cointop/apis/types"
)

// Interface interface
type Interface interface {
	GetAllCoinData(limit int) (map[string]types.Coin, error)
	GetCoinGraphData(coin string, start int64, end int64) (types.CoinGraph, error)
	GetGlobalMarketGraphData(start int64, end int64) (types.MarketGraph, error)
	//GetCoinData(coin string) (types.Coin, error)
	//GetGlobalMarketData() (types.GlobalMarketData, error)
	//GetAltcoinMarketGraphData(start int64, end int64) (types.MarketGraph, error)
	//GetCoinPriceUSD(coin string) (float64, error)
	//GetCoinMarkets(coin string) ([]types.Market, error)
}
