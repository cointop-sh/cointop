package coinmarketcap

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/miguelmota/cointop/pkg/cmc/types"
)

// V1Tickers get information about all coins listed in Coin Market Cap
func V1Tickers(limit int, convert string) (map[string]*types.V1Ticker, error) {
	var params []string
	if limit >= 0 {
		params = append(params, fmt.Sprintf("limit=%v", limit))
	}
	if convert != "" {
		params = append(params, fmt.Sprintf("convert=%v", convert))
	}
	baseURL := "https://api.coinmarketcap.com/v1"
	url := fmt.Sprintf("%s/ticker?%s", baseURL, strings.Join(params, "&"))

	resp, err := makeReq(url)
	var data []*types.V1Ticker
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	var mapstring []map[string]interface{}
	err = json.Unmarshal(resp, &mapstring)
	if err != nil {
		return nil, err
	}

	// creating map from the array
	allCoins := make(map[string]*types.V1Ticker)
	for i := 0; i < len(data); i++ {
		allCoins[data[i].ID] = data[i]
	}

	for _, item := range mapstring {
		id, _ := item["id"].(string)
		priceifc := item[fmt.Sprintf("price_%s", strings.ToLower(convert))]
		pricestr, _ := priceifc.(string)
		price, _ := strconv.ParseFloat(pricestr, 64)
		marketcapifc := item[fmt.Sprintf("market_cap_%s", strings.ToLower(convert))]
		marketcapstr, _ := marketcapifc.(string)
		marketcap, _ := strconv.ParseFloat(marketcapstr, 64)
		volumeifc := item[fmt.Sprintf("24h_volume_%s", strings.ToLower(convert))]
		volumestr, _ := volumeifc.(string)
		volume, _ := strconv.ParseFloat(volumestr, 64)
		quotes := &types.TickerQuote{
			Price:     price,
			Volume24H: volume,
			MarketCap: marketcap,
		}
		allCoins[id].Quotes = map[string]*types.TickerQuote{}
		allCoins[id].Quotes[strings.ToUpper(convert)] = quotes
	}

	return allCoins, nil
}
