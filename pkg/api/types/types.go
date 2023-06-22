package types

// Coin struct
type Coin struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Symbol           string  `json:"symbol"`
	Rank             int     `json:"rank"`
	Price            float64 `json:"price"`
	Volume24H        float64 `json:"volume24H"`
	MarketCap        float64 `json:"marketCap"`
	AvailableSupply  float64 `json:"availableSupply"`
	TotalSupply      float64 `json:"totalSupply"`
	PercentChange1H  float64 `json:"percentChange1H"`
	PercentChange24H float64 `json:"percentChange24H"`
	PercentChange7D  float64 `json:"percentChange7D"`
	PercentChange30D float64 `json:"percentChange30D"`
	PercentChange1Y  float64 `json:"percentChange1Y"`
	LastUpdated      string  `json:"lastUpdated"`
	// Slug uses to access the coin's info web page
	Slug string `json:"slug"`
}

// GlobalMarketData struct
type GlobalMarketData struct {
	TotalMarketCapUSD            float64 `json:"totalMarketCapUSD"`
	Total24HVolumeUSD            float64 `json:"total24HVolumeUSD"`
	BitcoinPercentageOfMarketCap float64 `json:"bitcoinPercentageOfMarketCap"`
	ActiveCurrencies             int     `json:"activeCurrencies"`
	ActiveAssets                 int     `json:"activeAssets"`
	ActiveMarkets                int     `json:"activeMarkets"`
}

// CoinGraph struct
type CoinGraph struct {
	MarketCapByAvailableSupply [][]float64 `json:"marketCapByAvailableSupply"`
	PriceBTC                   [][]float64 `json:"priceBTC"`
	Price                      [][]float64 `json:"price"`
	Volume                     [][]float64 `json:"volume"`
}

// Market struct
type Market struct {
	Rank          int     `json:"rank"`
	Exchange      string  `json:"exchange"`
	Pair          string  `json:"pair"`
	VolumeUSD     float64 `json:"volumeUSD"`
	Price         float64 `json:"price"`
	VolumePercent float64 `json:"volumePercent"`
	Updated       string  `json:"updated"`
}

// MarketGraph struct
type MarketGraph struct {
	MarketCapByAvailableSupply [][]float64 `json:"marketCapByAvailableSupply"`
	VolumeUSD                  [][]float64 `json:"volumeUSD"`
}
