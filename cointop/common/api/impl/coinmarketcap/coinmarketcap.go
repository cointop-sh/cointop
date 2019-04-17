package coinmarketcap

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	apitypes "github.com/miguelmota/cointop/cointop/common/api/types"
	cmc "github.com/miguelmota/go-coinmarketcap/pro/v1"
	cmcv2 "github.com/miguelmota/go-coinmarketcap/v2"
)

// Service service
type Service struct {
	client *cmc.Client
}

// New new service
func New() *Service {
	client := cmc.NewClient(&cmc.Config{
		ProAPIKey: os.Getenv("CMC_PRO_API_KEY"),
	})
	return &Service{
		client: client,
	}
}

// Ping ping API
func (s *Service) Ping() error {
	info, err := s.client.Cryptocurrency.Info(&cmc.InfoOptions{
		Symbol: "BTC",
	})
	if err != nil {
		return errors.New("failed to ping")
	}
	if info == nil {
		return errors.New("failed to ping")
	}
	return nil
}

func (s *Service) getLimitedCoinData(convert string, offset int) (map[string]apitypes.Coin, error) {
	ret := make(map[string]apitypes.Coin)
	max := 100

	listings, err := s.client.Cryptocurrency.LatestListings(&cmc.ListingOptions{
		Limit:   max,
		Convert: convert,
		Start:   max * offset,
	})
	if err != nil {
		return nil, err
	}
	for _, v := range listings {
		price := formatPrice(v.Quote[convert].Price, convert)
		lastUpdated, err := time.Parse(time.RFC3339, v.LastUpdated)
		if err != nil {
			return nil, err
		}
		ret[v.Name] = apitypes.Coin{
			ID:               strings.ToLower(v.Name),
			Name:             v.Name,
			Symbol:           v.Symbol,
			Rank:             int(v.CMCRank),
			AvailableSupply:  v.CirculatingSupply,
			TotalSupply:      v.TotalSupply,
			MarketCap:        float64(int(v.Quote[convert].MarketCap)),
			Price:            price,
			PercentChange1H:  v.Quote[convert].PercentChange1H,
			PercentChange24H: v.Quote[convert].PercentChange24H,
			PercentChange7D:  v.Quote[convert].PercentChange7D,
			Volume24H:        formatVolume(v.Quote[convert].Volume24H),
			LastUpdated:      strconv.Itoa(int(lastUpdated.Unix())),
		}
	}
	return ret, nil
}

// GetAllCoinData gets all coin data. Need to paginate through all pages
func (s *Service) GetAllCoinData(convert string) (chan map[string]apitypes.Coin, error) {
	var wg sync.WaitGroup
	ch := make(chan map[string]apitypes.Coin)
	go func() {
		var mutex sync.Mutex
		maxPages := 15
		for i := 0; i < maxPages; i++ {
			time.Sleep(time.Duration(i) * time.Second)
			wg.Add(1)
			go func(j int) {
				defer wg.Done()
				coins, err := s.getLimitedCoinData(convert, j)
				if err != nil {
					return
				}
				mutex.Lock()
				defer mutex.Unlock()
				ret := make(map[string]apitypes.Coin)
				for k, v := range coins {
					ret[k] = v
				}
				ch <- ret
			}(i)
		}
		wg.Wait()
	}()
	return ch, nil
}

// GetCoinGraphData gets coin graph data
func (s *Service) GetCoinGraphData(coin string, start int64, end int64) (apitypes.CoinGraph, error) {
	ret := apitypes.CoinGraph{}
	graphData, err := cmcv2.TickerGraph(&cmcv2.TickerGraphOptions{
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
	graphData, err := cmcv2.GlobalMarketGraph(&cmcv2.GlobalMarketGraphOptions{
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
	market, err := s.client.GlobalMetrics.LatestQuotes(&cmc.QuoteOptions{
		Convert: convert,
	})

	if err != nil {
		return ret, err
	}
	ret = apitypes.GlobalMarketData{
		TotalMarketCapUSD:            market.Quote[convert].TotalMarketCap,
		Total24HVolumeUSD:            market.Quote[convert].TotalVolume24H,
		BitcoinPercentageOfMarketCap: market.BTCDominance,
		ActiveCurrencies:             int(market.ActiveCryptocurrencies),
		ActiveAssets:                 0,
		ActiveMarkets:                int(market.ActiveMarketPairs),
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
