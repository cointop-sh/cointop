package coinmarketcap

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	apitypes "github.com/cointop-sh/cointop/pkg/api/types"
	"github.com/cointop-sh/cointop/pkg/api/util"
	cmc "github.com/miguelmota/go-coinmarketcap/pro/v1"
	cmcv2 "github.com/miguelmota/go-coinmarketcap/v2"
)

// ErrQuoteNotFound is the error for when a quote is not found
var ErrQuoteNotFound = errors.New("quote not found")

// ErrPingFailed is the error for when pinging the API fails
var ErrPingFailed = errors.New("failed to ping")

// ErrFetchGraphData is the error for when fetching graph data fails
var ErrFetchGraphData = errors.New("graph data fetch error")

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

func (s *Service) getPaginatedCoinData(convert string, offset int) ([]apitypes.Coin, error) {
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
			// TODO: Fix ID
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
			Slug:             util.FormatSlug(v.Slug),
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

			coins, err := s.getPaginatedCoinData(convert, i)
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
	coins, err := s.getPaginatedCoinData(convert, 0)
	if err != nil {
		return ret, err
	}

	for _, coin := range coins {
		if coin.Name == name {
			return coin, nil
		}
	}

	return ret, nil
}

// GetCoinDataBatch gets all data of specified coins.
func (s *Service) GetCoinDataBatch(names []string, convert string) ([]apitypes.Coin, error) {
	var ret []apitypes.Coin
	coins, err := s.getPaginatedCoinData(convert, 0)
	if err != nil {
		return ret, err
	}

	for _, coin := range coins {
		for _, name := range names {
			if coin.Name == name {
				ret = append(ret, coin)
				break
			}
		}
	}

	return ret, nil
}

// GetCoinGraphData gets coin graph data
func (s *Service) GetCoinGraphData(convert, symbol string, name string, start int64, end int64) (apitypes.CoinGraph, error) {
	ret := apitypes.CoinGraph{}
	symbol = strings.ToUpper(symbol)
	info, err := s.client.Cryptocurrency.Info(&cmc.InfoOptions{
		Slug: name,
	})
	if err != nil {
		return ret, err
	}
	var coinID string
	if len(info) == 0 {
		return ret, ErrFetchGraphData
	}
	for k := range info {
		coinID = fmt.Sprintf("%v", info[k].ID)
	}
	if convert == "" {
		convert = "usd"
	}
	convert = strings.ToUpper(convert)
	interval := getChartInterval(start, end)
	params := []string{
		fmt.Sprintf("convert=%s,%s", convert, symbol),
		"format=chart_crypto_details",
		fmt.Sprintf("id=%s", coinID),
		fmt.Sprintf("interval=%s", interval),
		fmt.Sprintf("time_start=%v", start),
		fmt.Sprintf("time_end=%v", end),
	}
	baseURL := "https://web-api.coinmarketcap.com/v1.1"
	url := fmt.Sprintf("%s/cryptocurrency/quotes/historical?%s", baseURL, strings.Join(params, "&"))
	resp, err := makeReq(url)
	if err != nil {
		return ret, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return ret, err
	}
	data, ok := result["data"]
	if !ok {
		return ret, ErrFetchGraphData
	}
	ifcs, ok := data.(map[string]interface{})
	if !ok {
		return ret, ErrFetchGraphData
	}
	var prices [][]float64
	for datetime, item := range ifcs {
		ifc, ok := item.(map[string]interface{})
		if !ok {
			return ret, ErrFetchGraphData
		}
		for key, obj := range ifc {
			if key != convert {
				continue
			}
			arrIfc, ok := obj.([]interface{})
			if !ok {
				return ret, ErrFetchGraphData
			}
			if len(arrIfc) == 0 {
				return ret, ErrFetchGraphData
			}
			val := arrIfc[0].(float64)
			t, err := time.Parse(time.RFC3339, datetime)
			if err != nil {
				return ret, err
			}
			prices = append(prices, []float64{float64(t.Unix()), val})
		}
	}
	sort.Slice(prices[:], func(i, j int) bool {
		return prices[i][0] < prices[j][0]
	})
	ret.Price = prices
	return ret, nil
}

// GetGlobalMarketGraphData gets global market graph data
func (s *Service) GetGlobalMarketGraphData(convert string, start int64, end int64) (apitypes.MarketGraph, error) {
	ret := apitypes.MarketGraph{}
	if convert == "" {
		convert = "usd"
	}
	convert = strings.ToUpper(convert)
	interval := getChartInterval(start, end)
	params := []string{
		fmt.Sprintf("convert=%s", convert),
		"format=chart",
		fmt.Sprintf("interval=%s", interval),
		fmt.Sprintf("time_start=%v", start),
		fmt.Sprintf("time_end=%v", end),
	}
	baseURL := "https://web-api.coinmarketcap.com/v1.1"
	url := fmt.Sprintf("%s/global-metrics/quotes/historical?%s", baseURL, strings.Join(params, "&"))
	resp, err := makeReq(url)
	if err != nil {
		return ret, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return ret, err
	}
	data, ok := result["data"]
	if !ok {
		return ret, ErrFetchGraphData
	}
	mapIfc, ok := data.(map[string]interface{})
	if !ok {
		return ret, ErrFetchGraphData
	}
	var marketCap [][]float64
	for datetime, item := range mapIfc {
		arrIfc, ok := item.([]interface{})
		if !ok {
			return ret, ErrFetchGraphData
		}
		if len(arrIfc) == 0 {
			return ret, ErrFetchGraphData
		}
		val := arrIfc[0].(float64)
		t, err := time.Parse(time.RFC3339, datetime)
		if err != nil {
			return ret, err
		}
		marketCap = append(marketCap, []float64{float64(t.Unix()), val})
	}
	sort.Slice(marketCap[:], func(i, j int) bool {
		return marketCap[i][0] < marketCap[j][0]
	})
	ret.MarketCapByAvailableSupply = marketCap
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
	symbol, err := cmcv2.CoinSymbol(util.NameToSlug(name))
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
func (s *Service) CoinLink(slug string) string {
	return fmt.Sprintf("https://coinmarketcap.com/currencies/%s/", slug)
}

// SupportedCurrencies returns a list of supported currencies
func (s *Service) SupportedCurrencies() []string {
	// keep these in alphabetical order
	return []string{
		"AUD",
		"BRL",
		"BTC",
		"CAD",
		"CFH",
		"CLP",
		"CNY",
		"CZK",
		"DKK",
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
		"MXN",
		"MYR",
		"NOK",
		"NZD",
		"PHP",
		"PKR",
		"PLN",
		"RUB",
		"SEK",
		"SGD",
		"THB",
		"TRY",
		"TWD",
		"USD",
		"VND",
		"ZAR",
	}
}

// doReq does HTTP request with client
func doReq(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s", body)
	}

	return body, nil
}

// makeReq is an HTTP GET request helper
func makeReq(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := doReq(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

// getChartInterval returns the interval to use for given time range
func getChartInterval(start, end int64) string {
	interval := "15m"
	delta := end - start
	if delta >= 604800 {
		interval = "1h"
	}
	if delta >= 2629746 {
		interval = "1d"
	}
	if delta >= 604800 {
		interval = "1h"
	}
	if delta >= 2592000 {
		interval = "1d"
	}
	return interval
}

// GetExchangeRate gets the current excange rate between two currencies
func (s *Service) GetExchangeRate(convertFrom, convertTo string, cached bool) (float64, error) {
	if convertFrom == convertTo {
		return 1.0, nil
	}
	return 0, fmt.Errorf("unsupported currency conversion: %s => %s", convertFrom, convertTo)
}
