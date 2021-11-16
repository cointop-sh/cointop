package coingecko

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	apitypes "github.com/cointop-sh/cointop/pkg/api/types"
	"github.com/cointop-sh/cointop/pkg/api/util"
	gecko "github.com/cointop-sh/cointop/pkg/api/vendors/coingecko/v3"
	"github.com/cointop-sh/cointop/pkg/api/vendors/coingecko/v3/types"
	geckoTypes "github.com/cointop-sh/cointop/pkg/api/vendors/coingecko/v3/types"
)

// ErrPingFailed is the error for when pinging the API fails
var ErrPingFailed = errors.New("failed to ping")

// ErrNotFound is the error when the target is not found
var ErrNotFound = errors.New("not found")

// Config config
type Config struct {
	PerPage  uint
	MaxPages uint
}

// Service service
type Service struct {
	client            *gecko.Client
	maxResultsPerPage uint
	maxPages          uint
	cacheMap          sync.Map
	cachedRates       *types.ExchangeRatesItem
}

// NewCoinGecko new service
func NewCoinGecko(config *Config) *Service {
	maxResultsPerPage := 250 // absolute max
	maxResults := uint(0)
	maxPages := uint(10)
	perPage := uint(100)
	if config.PerPage > 0 {
		perPage = config.PerPage
	}
	if config.MaxPages > 0 {
		maxPages = config.MaxPages
		maxResults = perPage * maxPages
		maxPages = uint(math.Ceil(math.Max(float64(maxResults)/float64(maxResultsPerPage), 1)))
	}

	client := gecko.NewClient(nil)
	svc := &Service{
		client:            client,
		maxResultsPerPage: uint(math.Min(float64(maxResults), float64(maxResultsPerPage))),
		maxPages:          maxPages,
		cacheMap:          sync.Map{},
	}
	svc.cacheCoinsIDList()
	return svc
}

// Ping ping API
func (s *Service) Ping() error {
	if _, err := s.client.Ping(); err != nil {
		return err
	}

	return nil
}

// GetAllCoinData gets all coin data. Need to paginate through all pages
func (s *Service) GetAllCoinData(convert string, ch chan []apitypes.Coin) error {
	go func() {
		defer close(ch)

		for i := 0; i < int(s.maxPages); i++ {
			if i > 0 {
				time.Sleep(1 * time.Second)
			}

			coins, err := s.getPaginatedCoinData(convert, i, []string{})
			if err != nil {
				return
			}

			ch <- coins
		}
	}()
	return nil
}

// GetCoinData gets all data of a coin.
func (s *Service) GetCoinData(name string, convert string) (apitypes.Coin, error) {
	ret := apitypes.Coin{}
	ids := []string{name}
	coins, err := s.getPaginatedCoinData(convert, 0, ids)
	if err != nil {
		return ret, err
	}

	if len(coins) > 0 {
		ret = coins[0]
	}

	return ret, nil
}

// GetCoinDataBatch gets all data of specified coins.
func (s *Service) GetCoinDataBatch(names []string, convert string) ([]apitypes.Coin, error) {
	return s.getPaginatedCoinData(convert, 0, names)
}

// GetCoinGraphData gets coin graph data
func (s *Service) GetCoinGraphData(convert, symbol, name string, start, end int64) (apitypes.CoinGraph, error) {
	ret := apitypes.CoinGraph{}
	days := strconv.Itoa(util.CalcDays(start, end))
	chart, err := s.client.CoinsIDMarketChart(s.coinNameToID(name), convert, days)
	if err != nil {
		return ret, err
	}

	var marketCap [][]float64
	var priceCoin [][]float64
	var priceBTC [][]float64
	var volumeCoin [][]float64

	if chart.Prices != nil {
		for _, item := range *chart.Prices {
			timestamp := float64(item[0])
			price := float64(item[1])

			priceCoin = append(priceCoin, []float64{
				timestamp,
				price,
			})
		}
	}

	ret.MarketCapByAvailableSupply = marketCap
	ret.PriceBTC = priceBTC
	ret.Price = priceCoin
	ret.Volume = volumeCoin

	return ret, nil
}

// GetExchangeRates returns the exchange rates from the backend, or a cached copy if requested and available
func (s *Service) GetExchangeRates(cached bool) (*types.ExchangeRatesItem, error) {
	if s.cachedRates == nil || !cached {
		rates, err := s.client.ExchangeRates()
		if err != nil {
			return nil, err
		}
		s.cachedRates = rates
	}
	return s.cachedRates, nil
}

// GetExchangeRate gets the current excange rate between two currencies
func (s *Service) GetExchangeRate(convertFrom, convertTo string, cached bool) (float64, error) {
	convertFrom = strings.ToLower(convertFrom)
	convertTo = strings.ToLower(convertTo)
	if convertFrom == convertTo {
		return 1.0, nil
	}
	rates, err := s.GetExchangeRates(cached)
	if err != nil {
		return 0, err
	}
	if rates == nil {
		return 0, fmt.Errorf("expected rates, received nil")
	}
	// Combined rate is convertFrom->BTC->convertTo
	fromRate, found := (*rates)[convertFrom]
	if !found {
		return 0, fmt.Errorf("unsupported currency conversion: %s", convertFrom)
	}
	toRate, found := (*rates)[convertTo]
	if !found {
		return 0, fmt.Errorf("unsupported currency conversion: %s", convertTo)
	}
	rate := toRate.Value / fromRate.Value
	return rate, nil
}

// GetGlobalMarketGraphData gets global market graph data
func (s *Service) GetGlobalMarketGraphData(convert string, start int64, end int64) (apitypes.MarketGraph, error) {
	days := strconv.Itoa(util.CalcDays(start, end))
	ret := apitypes.MarketGraph{}
	convertTo := strings.ToLower(convert)
	if convertTo == "" {
		convertTo = "usd"
	}
	graphData, err := s.client.GlobalCharts("usd", days)
	if err != nil {
		return ret, err
	}

	// This API does not appear to support vs_currency and only returns USD, so use ExchangeRates to convert
	// TODO: watch out - this is not cached, so we hit the backend every time!
	rate, err := s.GetExchangeRate("usd", convertTo, true)
	if err != nil {
		return ret, err
	}

	var marketCapUSD [][]float64
	var marketVolumeUSD [][]float64
	if graphData.Stats != nil {
		for _, item := range *graphData.Stats {
			marketCapUSD = append(marketCapUSD, []float64{
				float64(item[0]),
				float64(item[1]) * rate,
			})
		}
	}

	ret.MarketCapByAvailableSupply = marketCapUSD
	ret.VolumeUSD = marketVolumeUSD
	return ret, nil
}

// GetGlobalMarketData gets global market data
func (s *Service) GetGlobalMarketData(convert string) (apitypes.GlobalMarketData, error) {
	convert = strings.ToLower(convert)
	ret := apitypes.GlobalMarketData{}
	market, err := s.client.Global()
	if err != nil {
		return ret, err
	}

	totalMarketCap := market.TotalMarketCap[convert]
	totalVolume := market.TotalVolume[convert]
	btcDominance := market.MarketCapPercentage["btc"]

	ret = apitypes.GlobalMarketData{
		TotalMarketCapUSD:            totalMarketCap,
		Total24HVolumeUSD:            totalVolume,
		BitcoinPercentageOfMarketCap: btcDominance,
		ActiveCurrencies:             int(market.ActiveCryptocurrencies),
		ActiveAssets:                 0,
		ActiveMarkets:                int(market.Markets),
	}

	return ret, nil
}

// Price returns the current price of the coin
func (s *Service) Price(name string, convert string) (float64, error) {
	ids := []string{s.coinNameToID(name)}
	convert = strings.ToLower(convert)
	currencies := []string{convert}
	priceList, err := s.client.SimplePrice(ids, currencies)
	if err != nil {
		return 0, err
	}

	for _, item := range *priceList {
		if p, ok := item[convert]; ok {
			return util.FormatPrice(float64(p), convert), nil
		}
	}

	return 0, ErrNotFound
}

// CoinLink returns the URL link for the coin
func (s *Service) CoinLink(name string) string {
	ID := s.coinNameToID(name)
	return fmt.Sprintf("https://www.coingecko.com/en/coins/%s", ID)
}

// SupportedCurrencies returns a list of supported currencies
func (s *Service) SupportedCurrencies() []string {

	// keep these in alphabetical order
	return []string{
		"AED",
		"ARS",
		"AUD",
		"BDT",
		"BHD",
		"BMD",
		"BNB",
		"BRL",
		"BTC",
		"CAD",
		"CHF",
		"CLP",
		"CNY",
		"CZK",
		"DKK",
		"EOS",
		"ETH",
		"EUR",
		"GBP",
		"HKD",
		"HUF",
		"IDR",
		"ILS",
		"INR",
		"JPY",
		"KRW",
		"KWD",
		"LKR",
		"MMK",
		"MXN",
		"MYR",
		"NOK",
		"NZD",
		"PHP",
		"PKR",
		"PLN",
		"RUB",
		"SAR",
		"SATS",
		"SEK",
		"SGD",
		"THB",
		"TRY",
		"TWD",
		"UAH",
		"USD",
		"VEF",
		"VND",
		"XAG",
		"XDR",
		"ZAR",
	}
}

// cacheCoinsIDList fetches list of all coin IDS by name and symbols and caches it in a map for fast lookups
func (s *Service) cacheCoinsIDList() error {
	list, err := s.client.CoinsList()
	if err != nil {
		return err
	}
	if list == nil {
		return nil
	}
	var firstWords [][]string
	for _, item := range *list {
		keys := []string{
			strings.ToLower(item.Name),
			strings.ToLower(item.Symbol),
			util.NameToSlug(item.Name),
		}
		parts := strings.Split(strings.ToLower(item.Name), " ")
		if len(parts) > 1 {
			if parts[1] == "coin" {
				keys = append(keys, parts[0])
			} else {
				firstWords = append(firstWords, []string{parts[0], item.ID})
			}
		}
		for _, key := range keys {
			_, exists := s.cacheMap.Load(key)
			if !exists {
				s.cacheMap.Store(key, item.ID)
			}
		}
	}
	for _, parts := range firstWords {
		_, exists := s.cacheMap.Load(parts[0])
		if !exists {
			s.cacheMap.Store(parts[0], parts[1])
		}
	}
	return nil
}

// coinNameToID attempts to get coin ID based on coin name or coin symbol
func (s *Service) coinNameToID(name string) string {
	id, ok := s.cacheMap.Load(strings.ToLower(strings.TrimSpace(name)))
	if ok {
		return id.(string)
	}
	return util.NameToSlug(name)
}

// getPaginatedCoinData fetches coin data from page offset
func (s *Service) getPaginatedCoinData(convert string, offset int, names []string) ([]apitypes.Coin, error) {
	var ret []apitypes.Coin
	page := offset + 1 // page starts at 1
	sparkline := false
	pcp := geckoTypes.PriceChangePercentageObject
	priceChangePercentage := []string{
		pcp.PCP1h,
		pcp.PCP24h,
		pcp.PCP7d,
		pcp.PCP30d,
		pcp.PCP1y,
	}
	order := geckoTypes.OrderTypeObject.MarketCapDesc
	convertTo := strings.ToLower(convert)
	if convertTo == "" {
		convertTo = "usd"
	}

	ids := make([]string, len(names))
	for i, name := range names {
		ids[i] = s.coinNameToID(name)
	}
	list, err := s.client.CoinsMarket(convertTo, ids, order, int(s.maxResultsPerPage), page, sparkline, priceChangePercentage)
	if err != nil {
		return nil, err
	}

	if list != nil {
		// for fetching "simple prices"
		currencies := make([]string, len(*list))
		for i, item := range *list {
			currencies[i] = item.Name
		}

		for _, item := range *list {
			price := item.CurrentPrice
			var percentChange1H float64
			var percentChange24H float64
			var percentChange7D float64
			var percentChange30D float64
			var percentChange1Y float64

			if item.PriceChangePercentage1hInCurrency != nil {
				percentChange1H = *item.PriceChangePercentage1hInCurrency
			}
			if item.PriceChangePercentage24hInCurrency != nil {
				percentChange24H = *item.PriceChangePercentage24hInCurrency
			}
			if item.PriceChangePercentage7dInCurrency != nil {
				percentChange7D = *item.PriceChangePercentage7dInCurrency
			}
			if item.PriceChangePercentage30dInCurrency != nil {
				percentChange30D = *item.PriceChangePercentage30dInCurrency
			}
			if item.PriceChangePercentage1yInCurrency != nil {
				percentChange1Y = *item.PriceChangePercentage1yInCurrency
			}

			availableSupply := item.CirculatingSupply
			totalSupply := item.TotalSupply
			if totalSupply == 0 {
				totalSupply = availableSupply
			}

			ret = append(ret, apitypes.Coin{
				ID:               util.FormatID(item.ID),
				Name:             util.FormatName(item.Name),
				Symbol:           util.FormatSymbol(item.Symbol),
				Rank:             util.FormatRank(item.MarketCapRank),
				AvailableSupply:  util.FormatSupply(availableSupply),
				TotalSupply:      util.FormatSupply(totalSupply),
				MarketCap:        util.FormatMarketCap(item.MarketCap),
				Price:            util.FormatPrice(price, convert),
				PercentChange1H:  util.FormatPercentChange(percentChange1H),
				PercentChange24H: util.FormatPercentChange(percentChange24H),
				PercentChange7D:  util.FormatPercentChange(percentChange7D),
				PercentChange30D: util.FormatPercentChange(percentChange30D),
				PercentChange1Y:  util.FormatPercentChange(percentChange1Y),
				Volume24H:        util.FormatVolume(item.TotalVolume),
				LastUpdated:      util.FormatLastUpdated(item.LastUpdated),
			})
		}
	}

	return ret, nil
}
