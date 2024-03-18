// Package coingecko is forked from https://github.com/superoo7/go-gecko
package coingecko

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"os"

	"github.com/cointop-sh/cointop/pkg/api/vendors/coingecko/format"
	"github.com/cointop-sh/cointop/pkg/api/vendors/coingecko/v3/types"
	log "github.com/sirupsen/logrus"
)

// Client struct
type Client struct {
	httpClient *http.Client
	apiKey     string
}

// NewClient create new client object
func NewClient(httpClient *http.Client, apiKey string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{httpClient: httpClient, apiKey: apiKey}
}

// helper
// doReq HTTP client
func doReq(req *http.Request, client *http.Client) ([]byte, error) {
	debugHttp := os.Getenv("DEBUG_HTTP") != ""
	if debugHttp {
		log.Debugf("doReq %s %s", req.Method, req.URL)
	}

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
		if debugHttp {
			log.Warnf("doReq Got Status '%s' from %s %s", resp.Status, req.Method, req.URL)
			log.Debugf("doReq Got Body: %s", body)
		}
		return nil, fmt.Errorf("%s", body)
	}
	return body, nil
}

func (c *Client) getApiUrl(path string, params *url.Values) string {
	urlParams := url.Values{}
	subdomain := "api"
	if params != nil {
		urlParams = *params
	}
	if c.apiKey != "" {
		subdomain = "pro-api"
		urlParams.Add("x_cg_pro_api_key", c.apiKey)
	}
	url := fmt.Sprintf("https://%s.coingecko.com/api/v3%s?%s", subdomain, path, urlParams.Encode())
	return url
}

// MakeReq HTTP request helper
func (c *Client) MakeReq(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}
	resp, err := doReq(req, c.httpClient)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// API

// Ping /ping endpoint
func (c *Client) Ping() (*types.Ping, error) {
	url := c.getApiUrl("/ping", nil)
	resp, err := c.MakeReq(url)
	if err != nil {
		return nil, err
	}
	var data *types.Ping
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SimpleSinglePrice /simple/price  Single ID and Currency (ids, vsCurrency)
func (c *Client) SimpleSinglePrice(id string, vsCurrency string) (*types.SimpleSinglePrice, error) {
	idParam := []string{strings.ToLower(id)}
	vcParam := []string{strings.ToLower(vsCurrency)}

	t, err := c.SimplePrice(idParam, vcParam)
	if err != nil {
		return nil, err
	}
	curr := (*t)[id]
	data := &types.SimpleSinglePrice{ID: id, Currency: vsCurrency, MarketPrice: curr[vsCurrency]}
	return data, nil
}

// SimplePrice /simple/price Multiple ID and Currency (ids, vs_currencies)
func (c *Client) SimplePrice(ids []string, vsCurrencies []string) (*map[string]map[string]float32, error) {
	params := &url.Values{}
	idsParam := strings.Join(ids[:], ",")
	vsCurrenciesParam := strings.Join(vsCurrencies[:], ",")

	params.Add("ids", idsParam)
	params.Add("vs_currencies", vsCurrenciesParam)

	url := c.getApiUrl("/simple/price", params)
	resp, err := c.MakeReq(url)
	if err != nil {
		return nil, err
	}

	t := make(map[string]map[string]float32)
	err = json.Unmarshal(resp, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// SimpleSupportedVSCurrencies /simple/supported_vs_currencies
func (c *Client) SimpleSupportedVSCurrencies() (*types.SimpleSupportedVSCurrencies, error) {
	url := c.getApiUrl("/simple/supported_vs_currencies", nil)
	resp, err := c.MakeReq(url)
	if err != nil {
		return nil, err
	}
	var data *types.SimpleSupportedVSCurrencies
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// CoinsList /coins/list
func (c *Client) CoinsList() (*types.CoinList, error) {
	url := c.getApiUrl("/coins/list", nil)
	resp, err := c.MakeReq(url)
	if err != nil {
		return nil, err
	}

	var data *types.CoinList
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// CoinsMarket /coins/market
func (c *Client) CoinsMarket(vsCurrency string, ids []string, order string, perPage int, page int, sparkline bool, priceChangePercentage []string) (*types.CoinsMarket, error) {
	if len(vsCurrency) == 0 {
		return nil, fmt.Errorf("vsCurrency is required")
	}
	params := &url.Values{}
	// vsCurrency
	params.Add("vs_currency", vsCurrency)
	// order
	if len(order) == 0 {
		order = types.OrderTypeObject.MarketCapDesc
	}
	params.Add("order", order)
	// ids
	if len(ids) != 0 {
		idsParam := strings.Join(ids[:], ",")
		params.Add("ids", idsParam)
	}
	// per_page
	if perPage <= 0 || perPage > 250 {
		perPage = 100
	}
	params.Add("per_page", format.Int2String(perPage))
	params.Add("page", format.Int2String(page))
	// sparkline
	params.Add("sparkline", format.Bool2String(sparkline))
	// price_change_percentage
	if len(priceChangePercentage) != 0 {
		priceChangePercentageParam := strings.Join(priceChangePercentage[:], ",")
		params.Add("price_change_percentage", priceChangePercentageParam)
	}
	url := c.getApiUrl("/coins/markets", params)
	resp, err := c.MakeReq(url)
	if err != nil {
		return nil, err
	}
	var data *types.CoinsMarket
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// CoinsID /coins/{id}
func (c *Client) CoinsID(id string, localization bool, tickers bool, marketData bool, communityData bool, developerData bool, sparkline bool) (*types.CoinsID, error) {

	if len(id) == 0 {
		return nil, fmt.Errorf("id is required")
	}
	params := &url.Values{}
	params.Add("localization", format.Bool2String(localization))
	params.Add("tickers", format.Bool2String(tickers))
	params.Add("market_data", format.Bool2String(marketData))
	params.Add("community_data", format.Bool2String(communityData))
	params.Add("developer_data", format.Bool2String(developerData))
	params.Add("sparkline", format.Bool2String(sparkline))
	url := c.getApiUrl(fmt.Sprintf("/coins/%s", id), params)
	resp, err := c.MakeReq(url)
	if err != nil {
		return nil, err
	}

	var data *types.CoinsID
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// CoinsIDTickers /coins/{id}/tickers
func (c *Client) CoinsIDTickers(id string, page int) (*types.CoinsIDTickers, error) {
	if len(id) == 0 {
		return nil, fmt.Errorf("id is required")
	}
	params := &url.Values{}
	if page > 0 {
		params.Add("page", format.Int2String(page))
	}
	url := c.getApiUrl(fmt.Sprintf("/coins/%s/tickers", id), params)
	resp, err := c.MakeReq(url)
	if err != nil {
		return nil, err
	}
	var data *types.CoinsIDTickers
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// CoinsIDHistory /coins/{id}/history?date={date}&localization=false
func (c *Client) CoinsIDHistory(id string, date string, localization bool) (*types.CoinsIDHistory, error) {
	if len(id) == 0 || len(date) == 0 {
		return nil, fmt.Errorf("id and date is required")
	}
	params := &url.Values{}
	params.Add("date", date)
	params.Add("localization", format.Bool2String(localization))

	url := c.getApiUrl(fmt.Sprintf("/coins/%s/history", id), params)
	resp, err := c.MakeReq(url)
	if err != nil {
		return nil, err
	}
	var data *types.CoinsIDHistory
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// CoinsIDMarketChart /coins/{id}/market_chart?vsCurrency={usd, eur, jpy, etc.}&days={1,14,30,max}
func (c *Client) CoinsIDMarketChart(id string, vsCurrency string, days string) (*types.CoinsIDMarketChart, error) {
	if len(id) == 0 || len(vsCurrency) == 0 || len(days) == 0 {
		return nil, fmt.Errorf("id, vsCurrency, and days is required")
	}

	params := &url.Values{}
	params.Add("vs_currency", vsCurrency)
	params.Add("days", days)

	url := c.getApiUrl(fmt.Sprintf("/coins/%s/market_chart", id), params)
	resp, err := c.MakeReq(url)
	if err != nil {
		return nil, err
	}

	m := types.CoinsIDMarketChart{}
	err = json.Unmarshal(resp, &m)
	if err != nil {
		return &m, err
	}

	return &m, nil
}

// CoinsIDStatusUpdates

// CoinsIDContractAddress https://api.coingecko.com/api/v3/coins/{id}/contract/{contract_address}
// func CoinsIDContractAddress(id string, address string) (nil, error) {
// 	url := c.getApiUrl(fmt.Sprintf("/coins/%s/contract/%s", id, address), nil)
// 	resp, err := request.MakeReq(url)
// 	if err != nil {
// 		return nil, err
// 	}
// }

// EventsCountries https://api.coingecko.com/api/v3/events/countries
func (c *Client) EventsCountries() ([]types.EventCountryItem, error) {
	url := c.getApiUrl("/events/countries", nil)
	resp, err := c.MakeReq(url)
	if err != nil {
		return nil, err
	}
	var data *types.EventsCountries
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return data.Data, nil

}

// EventsTypes https://api.coingecko.com/api/v3/events/types
func (c *Client) EventsTypes() (*types.EventsTypes, error) {
	url := c.getApiUrl("/events/types", nil)
	resp, err := c.MakeReq(url)
	if err != nil {
		return nil, err
	}
	var data *types.EventsTypes
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return data, nil

}

// ExchangeRates https://api.coingecko.com/api/v3/exchange_rates
func (c *Client) ExchangeRates() (*types.ExchangeRatesItem, error) {
	url := c.getApiUrl("/exchange_rates", nil)
	resp, err := c.MakeReq(url)
	if err != nil {
		return nil, err
	}
	var data *types.ExchangeRatesResponse
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return &data.Rates, nil
}

// Global https://api.coingecko.com/api/v3/global
func (c *Client) Global() (*types.Global, error) {
	url := c.getApiUrl("/global", nil)
	resp, err := c.MakeReq(url)
	if err != nil {
		return nil, err
	}
	var data *types.GlobalResponse
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return &data.Data, nil
}

// GlobalCharts https://www.coingecko.com/market_cap/total_charts_data?duration=7&locale=en&vs_currency=usd
func (c *Client) GlobalCharts(vsCurrency string, days string) (*types.GlobalCharts, error) {
	if len(vsCurrency) == 0 || len(days) == 0 {
		return nil, fmt.Errorf("vsCurrency, and days is required")
	}

	params := url.Values{}
	params.Add("locale", "en")
	params.Add("vs_currency", vsCurrency)
	params.Add("duration", days)

	url := fmt.Sprintf("https://www.coingecko.com/market_cap/total_charts_data?%s", params.Encode())
	resp, err := c.MakeReq(url)
	if err != nil {
		return nil, err
	}

	m := types.GlobalCharts{}
	err = json.Unmarshal(resp, &m)
	if err != nil {
		return &m, err
	}

	return &m, nil
}
