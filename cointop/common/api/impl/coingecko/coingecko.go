package coingecko

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	gecko "github.com/miguelmota/cointop/cointop/api/coingecko/v3"
	geckoTypes "github.com/miguelmota/cointop/cointop/api/coingecko/v3/types"
	apitypes "github.com/miguelmota/cointop/cointop/common/api/types"
	util "github.com/miguelmota/cointop/cointop/common/api/util"
)

// ErrPingFailed is the error for when pinging the API fails
var ErrPingFailed = errors.New("Failed to ping")

// ErrNotFound is the error when the target is not found
var ErrNotFound = errors.New("Not found")

// Service service
type Service struct {
	client            *gecko.Client
	maxResultsPerPage int
	maxPages          int
}

// NewCoinGecko new service
func NewCoinGecko() *Service {
	client := gecko.NewClient(nil)
	return &Service{
		client:            client,
		maxResultsPerPage: 250,
		maxPages:          10,
	}
}

// Ping ping API
func (s *Service) Ping() error {
	if _, err := s.client.Ping(); err != nil {
		return err
	}

	return nil
}

func (s *Service) getPaginatedCoinData(convert string, offset int, names []string) ([]apitypes.Coin, error) {
	var ret []apitypes.Coin
	page := offset
	sparkline := false
	pcp := geckoTypes.PriceChangePercentageObject
	priceChangePercentage := []string{pcp.PCP1h, pcp.PCP24h, pcp.PCP7d}
	order := geckoTypes.OrderTypeObject.MarketCapDesc
	convertTo := strings.ToLower(convert)
	if convertTo == "" {
		convertTo = "usd"
	}

	ids := make([]string, len(names))
	for i, name := range names {
		slug := util.NameToSlug(name)
		ids[i] = slug
	}
	list, err := s.client.CoinsMarket(convertTo, ids, order, s.maxResultsPerPage, page, sparkline, priceChangePercentage)
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

			if item.PriceChangePercentage1hInCurrency != nil {
				percentChange1H = *item.PriceChangePercentage1hInCurrency
			}

			if item.PriceChangePercentage24hInCurrency != nil {
				percentChange24H = *item.PriceChangePercentage24hInCurrency
			}

			if item.PriceChangePercentage7dInCurrency != nil {
				percentChange7D = *item.PriceChangePercentage7dInCurrency
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
				Volume24H:        util.FormatVolume(item.TotalVolume),
				LastUpdated:      util.FormatLastUpdated(item.LastUpdated),
			})
		}
	}

	return ret, nil
}

// GetAllCoinData gets all coin data. Need to paginate through all pages
func (s *Service) GetAllCoinData(convert string, ch chan []apitypes.Coin) error {
	go func() {
		defer close(ch)

		for i := 0; i < s.maxPages; i++ {
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
	chart, err := s.client.CoinsIDMarketChart(util.NameToSlug(name), convert, days)
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

// GetGlobalMarketGraphData gets global market graph data
func (s *Service) GetGlobalMarketGraphData(convert string, start int64, end int64) (apitypes.MarketGraph, error) {
	days := strconv.Itoa(util.CalcDays(start, end))
	ret := apitypes.MarketGraph{}
	convertTo := strings.ToLower(convert)
	if convertTo == "" {
		convertTo = "usd"
	}
	graphData, err := s.client.GlobalCharts(convertTo, days)
	if err != nil {
		return ret, err
	}

	var marketCapUSD [][]float64
	var marketVolumeUSD [][]float64
	if graphData.Stats != nil {
		for _, item := range *graphData.Stats {
			marketCapUSD = append(marketCapUSD, []float64{
				float64(item[0]),
				float64(item[1]),
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
	list, err := s.client.CoinsList()
	if err != nil {
		return 0, err
	}

	for _, item := range *list {
		if item.Symbol == strings.ToLower(name) {
			name = item.Name
		}
	}

	ids := []string{util.NameToSlug(name)}
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
	slug := util.NameToSlug(name)
	return fmt.Sprintf("https://www.coingecko.com/en/coins/%s", slug)
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
		"SEK",
		"SGD",
		"THB",
		"TRY",
		"TWD",
		"USD",
		"VEF",
		"VND",
		"XAG",
		"XDR",
		"ZAR",
	}
}
