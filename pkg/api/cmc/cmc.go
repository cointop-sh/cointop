package api

import (
	types "github.com/miguelmota/cointop/pkg/api/types"
	cmc "github.com/miguelmota/go-coinmarketcap"
)

// Service service
type Service struct {
}

// New new service
func New() *Service {
	return &Service{}
}

// GetAllCoinData gets all coin data
func (s *Service) GetAllCoinData() (map[string]types.Coin, error) {
	ret := make(map[string]types.Coin)
	coins, err := cmc.GetAllCoinData(0)
	if err != nil {
		return ret, err
	}
	for _, v := range coins {
		ret[v.ID] = types.Coin{
			ID:               v.ID,
			Name:             v.Name,
			Symbol:           v.Symbol,
			Rank:             v.Rank,
			PriceUSD:         v.PriceUSD,
			PriceBTC:         v.PriceBTC,
			USD24HVolume:     v.USD24HVolume,
			MarketCapUSD:     v.MarketCapUSD,
			AvailableSupply:  v.AvailableSupply,
			TotalSupply:      v.TotalSupply,
			PercentChange1H:  v.PercentChange1H,
			PercentChange24H: v.PercentChange24H,
			PercentChange7D:  v.PercentChange7D,
			LastUpdated:      v.LastUpdated,
		}
	}
	return ret, nil
}

// GetCoinGraphData gets coin graph data
func (s *Service) GetCoinGraphData(coin string, start int64, end int64) (types.CoinGraph, error) {
	ret := types.CoinGraph{}
	graphData, err := cmc.GetCoinGraphData(coin, start, end)
	if err != nil {
		return ret, err
	}

	ret.MarketCapByAvailableSupply = graphData.MarketCapByAvailableSupply
	ret.PriceBTC = graphData.PriceBTC
	ret.PriceUSD = graphData.PriceUSD
	ret.VolumeUSD = graphData.VolumeUSD
	return ret, nil
}

// GetGlobalMarketGraphData gets global market graph data
func (s *Service) GetGlobalMarketGraphData(start int64, end int64) (types.MarketGraph, error) {
	ret := types.MarketGraph{}
	graphData, err := cmc.GetGlobalMarketGraphData(start, end)
	if err != nil {
		return ret, err
	}

	ret.MarketCapByAvailableSupply = graphData.MarketCapByAvailableSupply
	ret.VolumeUSD = graphData.VolumeUSD
	return ret, nil
}

// GetGlobalMarketData gets global market data
func (s *Service) GetGlobalMarketData() (types.GlobalMarketData, error) {
	ret := types.GlobalMarketData{}
	market, err := cmc.GetGlobalMarketData()
	if err != nil {
		return ret, err
	}
	ret = types.GlobalMarketData{
		TotalMarketCapUSD:            market.TotalMarketCapUSD,
		Total24HVolumeUSD:            market.Total24HVolumeUSD,
		BitcoinPercentageOfMarketCap: market.BitcoinPercentageOfMarketCap,
		ActiveCurrencies:             market.ActiveCurrencies,
		ActiveAssets:                 market.ActiveCurrencies,
		ActiveMarkets:                market.ActiveAssets,
	}
	return ret, nil
}
