package api

import (
	"github.com/cointop-sh/cointop/pkg/api/types"
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
	CoinLink(slug string) string
	AltCoinLink(coin types.Coin) string
	SupportedCurrencies() []string
	Price(name string, convert string) (float64, error)
	GetExchangeRate(convertFrom, convertTo string, cached bool) (float64, error) // I don't love this caching
}
