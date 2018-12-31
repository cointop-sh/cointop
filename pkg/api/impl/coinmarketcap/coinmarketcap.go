package coinmarketcap

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	apitypes "github.com/miguelmota/cointop/pkg/api/types"
	cmc "github.com/miguelmota/cointop/pkg/cmc"
)

// Service service
type Service struct {
}

// New new service
func New() *Service {
	return &Service{}
}

// Ping ping API
func (s *Service) Ping() error {
	ticker, err := cmc.Ticker(&cmc.TickerOptions{
		Symbol: "ETH",
	})
	if err != nil {
		return errors.New("failed to ping")
	}
	if ticker == nil {
		return errors.New("failed to ping")
	}
	return nil
}

func getLimitedCoinDataV2(convert string, offset int) (map[string]apitypes.Coin, error) {
	ret := make(map[string]apitypes.Coin)
	max := 100
	coins, err := cmc.Tickers(&cmc.TickersOptions{
		Convert: convert,
		Start:   max * offset,
		Limit:   max,
	})
	if err != nil {
		return ret, err
	}
	for _, v := range coins {
		ret[v.Name] = apitypes.Coin{
			ID:               v.Slug,
			Name:             v.Name,
			Symbol:           v.Symbol,
			Rank:             v.Rank,
			AvailableSupply:  v.CirculatingSupply,
			TotalSupply:      v.TotalSupply,
			MarketCap:        v.Quotes[convert].MarketCap,
			Price:            v.Quotes[convert].Price,
			PercentChange1H:  v.Quotes[convert].PercentChange1H,
			PercentChange24H: v.Quotes[convert].PercentChange24H,
			PercentChange7D:  v.Quotes[convert].PercentChange7D,
			Volume24H:        formatVolume(v.Quotes[convert].Volume24H),
			LastUpdated:      strconv.Itoa(v.LastUpdated),
		}
	}
	return ret, nil
}

func getLimitedCoinData(convert string, offset int) (map[string]apitypes.Coin, error) {
	ret := make(map[string]apitypes.Coin)
	max := 100
	coins, err := cmc.Tickers(&cmc.TickersOptions{
		Convert: convert,
		Start:   max * offset,
		Limit:   max,
	})
	if err != nil {
		return ret, err
	}
	for _, v := range coins {
		price := formatPrice(v.Quotes[convert].Price, convert)
		ret[v.Name] = apitypes.Coin{
			ID:               strings.ToLower(v.Name),
			Name:             v.Name,
			Symbol:           v.Symbol,
			Rank:             v.Rank,
			AvailableSupply:  v.CirculatingSupply,
			TotalSupply:      v.TotalSupply,
			MarketCap:        v.Quotes[convert].MarketCap,
			Price:            price,
			PercentChange1H:  v.Quotes[convert].PercentChange1H,
			PercentChange24H: v.Quotes[convert].PercentChange24H,
			PercentChange7D:  v.Quotes[convert].PercentChange7D,
			Volume24H:        formatVolume(v.Quotes[convert].Volume24H),
			LastUpdated:      strconv.Itoa(v.LastUpdated),
		}
	}
	return ret, nil
}

// GetAllCoinData1 get all coin data
func (s *Service) GetAllCoinData1(convert string) (map[string]apitypes.Coin, error) {
	ret := make(map[string]apitypes.Coin)
	coins, err := cmc.V1Tickers(0, convert)
	if err != nil {
		return ret, err
	}
	for _, v := range coins {
		price := formatPrice(v.Quotes[convert].Price, convert)
		ret[v.Name] = apitypes.Coin{
			ID:               strings.ToLower(v.Name),
			Name:             v.Name,
			Symbol:           v.Symbol,
			Rank:             v.Rank,
			AvailableSupply:  v.AvailableSupply,
			TotalSupply:      v.TotalSupply,
			MarketCap:        v.Quotes[convert].MarketCap,
			Price:            price,
			PercentChange1H:  v.PercentChange1H,
			PercentChange24H: v.PercentChange24H,
			PercentChange7D:  v.PercentChange7D,
			Volume24H:        formatVolume(v.Quotes[convert].Volume24H),
			LastUpdated:      strconv.Itoa(v.LastUpdated),
		}
	}
	return ret, nil
}

// GetAllCoinDataV2 gets all coin data
// V1 is currently better for fetching all the coins at once
func (s *Service) GetAllCoinDataV2(convert string) (map[string]apitypes.Coin, error) {
	var wg sync.WaitGroup
	ret := make(map[string]apitypes.Coin)
	var mutex sync.Mutex
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			coins, err := getLimitedCoinDataV2(convert, j)
			if err != nil {
				return
			}
			mutex.Lock()
			for k, v := range coins {
				ret[k] = v
			}
			mutex.Unlock()
		}(i)
	}
	wg.Wait()
	return ret, nil
}

// GetAllCoinData gets all coin data. Need to paginate through all pages
func (s *Service) GetAllCoinData(convert string) (map[string]apitypes.Coin, error) {
	var wg sync.WaitGroup
	ret := make(map[string]apitypes.Coin)
	var mutex sync.Mutex
	for i := 0; i < 17; i++ {
		time.Sleep(500 * time.Millisecond)
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			coins, err := getLimitedCoinData(convert, j)
			if err != nil {
				return
			}
			mutex.Lock()
			for k, v := range coins {
				ret[k] = v
			}
			mutex.Unlock()
		}(i)
	}
	wg.Wait()
	return ret, nil
}

// GetCoinGraphData gets coin graph data
func (s *Service) GetCoinGraphData(coin string, start int64, end int64) (apitypes.CoinGraph, error) {
	ret := apitypes.CoinGraph{}
	graphData, err := cmc.TickerGraph(&cmc.TickerGraphOptions{
		Symbol: coin,
		Start:  start,
		End:    end,
	})
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
func (s *Service) GetGlobalMarketGraphData(start int64, end int64) (apitypes.MarketGraph, error) {
	ret := apitypes.MarketGraph{}
	graphData, err := cmc.GlobalMarketGraph(&cmc.GlobalMarketGraphOptions{
		Start: start,
		End:   end,
	})
	if err != nil {
		return ret, err
	}

	ret.MarketCapByAvailableSupply = graphData.MarketCapByAvailableSupply
	ret.VolumeUSD = graphData.VolumeUSD
	return ret, nil
}

// GetGlobalMarketData gets global market data
func (s *Service) GetGlobalMarketData(convert string) (apitypes.GlobalMarketData, error) {
	ret := apitypes.GlobalMarketData{}
	market, err := cmc.GlobalMarket(&cmc.GlobalMarketOptions{
		Convert: convert,
	})
	if err != nil {
		return ret, err
	}
	ret = apitypes.GlobalMarketData{
		TotalMarketCapUSD:            market.Quotes[convert].TotalMarketCap,
		Total24HVolumeUSD:            market.Quotes[convert].TotalVolume24H,
		BitcoinPercentageOfMarketCap: market.BitcoinPercentageOfMarketCap,
		ActiveCurrencies:             market.ActiveCurrencies,
		ActiveAssets:                 0,
		ActiveMarkets:                market.ActiveMarkets,
	}
	return ret, nil
}

func formatPrice(price float64, convert string) float64 {
	pricestr := fmt.Sprintf("%.2f", price)
	if convert == "ETH" || convert == "BTC" || price < 1 {
		pricestr = fmt.Sprintf("%.5f", price)
	}
	price, _ = strconv.ParseFloat(pricestr, 64)
	return price
}

func formatVolume(volume float64) float64 {
	return float64(int64(volume))
}
