package coinmarketcap

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	apitypes "github.com/miguelmota/cointop/cointop/common/api/types"
	util "github.com/miguelmota/cointop/cointop/common/api/util"
	cmc "github.com/miguelmota/go-coinmarketcap/pro/v1"
	cmcv2 "github.com/miguelmota/go-coinmarketcap/v2"
)

// ErrQuoteNotFound is the error for when a quote is not found
var ErrQuoteNotFound = errors.New("Quote not found")

// ErrPingFailed is the error for when pinging the API fails
var ErrPingFailed = errors.New("Failed to ping")

// Service service
type Service struct {
	client *cmc.Client
}

// NewCMC new service
func NewCMC(apiKey string) *Service {
	if apiKey == "" {
		apiKey = os.Getenv("CMC_PRO_API_KEY")
	}
	client := cmc.NewClient(&cmc.Config{
		ProAPIKey: apiKey,
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
		return ErrPingFailed
	}
	if info == nil {
		return ErrPingFailed
	}
	return nil
}

func (s *Service) getLimitedCoinData(convert string, offset int) ([]apitypes.Coin, error) {
	var ret []apitypes.Coin
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
		quote, ok := v.Quote[convert]
		if !ok {
			return nil, ErrQuoteNotFound
		}

		ret = append(ret, apitypes.Coin{
			ID:               util.FormatID(v.Name),
			Name:             util.FormatName(v.Name),
			Symbol:           util.FormatSymbol(v.Symbol),
			Rank:             util.FormatRank(v.CMCRank),
			AvailableSupply:  util.FormatSupply(v.CirculatingSupply),
			TotalSupply:      util.FormatSupply(v.TotalSupply),
			MarketCap:        util.FormatMarketCap(quote.MarketCap),
			Price:            util.FormatPrice(v.Quote[convert].Price, convert),
			PercentChange1H:  util.FormatPercentChange(quote.PercentChange1H),
			PercentChange24H: util.FormatPercentChange(quote.PercentChange24H),
			PercentChange7D:  util.FormatPercentChange(quote.PercentChange7D),
			Volume24H:        util.FormatVolume(v.Quote[convert].Volume24H),
			LastUpdated:      util.FormatLastUpdated(v.LastUpdated),
		})
	}
	return ret, nil
}

// GetAllCoinData gets all coin data. Need to paginate through all pages
func (s *Service) GetAllCoinData(convert string, ch chan []apitypes.Coin) error {
	go func() {
		maxPages := 10
		defer close(ch)
		for i := 0; i < maxPages; i++ {
			if i > 0 {
				time.Sleep(1 * time.Second)
			}

			coins, err := s.getLimitedCoinData(convert, i)
			if err != nil {
				return
			}

			ch <- coins
		}
	}()
	return nil
}

// GetCoinGraphData gets coin graph data
func (s *Service) GetCoinGraphData(symbol string, name string, start int64, end int64) (apitypes.CoinGraph, error) {
	ret := apitypes.CoinGraph{}
	graphData, err := cmcv2.TickerGraph(&cmcv2.TickerGraphOptions{
		Symbol: symbol,
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

// Price returns the current price of the coin
func (s *Service) Price(name string, convert string) (float64, error) {
	convert = strings.ToUpper(convert)
	symbol, err := cmcv2.CoinSymbol(name)
	if err != nil {
		return 0, err
	}

	price, err := cmcv2.Price(&cmcv2.PriceOptions{
		Symbol:  symbol,
		Convert: convert,
	})
	if err != nil {
		return 0, err
	}

	return util.FormatPrice(price, convert), nil
}

// CoinLink returns the URL link for the coin
func (s *Service) CoinLink(name string) string {
	slug := util.NameToSlug(name)
	return fmt.Sprintf("https://coinmarketcap.com/currencies/%s", slug)
}

// SupportedCurrencies returns a list of supported currencies
func (s *Service) SupportedCurrencies() []string {
	return []string{
		"BTC",
		"ETH",
		"AUD",
		"BRL",
		"CAD",
		"CFH",
		"CLP",
		"CNY",
		"CZK",
		"DKK",
		"EUR",
		"GBP",
		"HKD",
		"HUF",
		"IDR",
		"ILS",
		"INR",
		"JPY",
		"KRW",
		"MXN",
		"MYR",
		"NOK",
		"NZD",
		"PLN",
		"PHP",
		"PKR",
		"RUB",
		"SEK",
		"SGD",
		"THB",
		"TRY",
		"TWD",
		"USD",
		"ZAR",
	}
}
