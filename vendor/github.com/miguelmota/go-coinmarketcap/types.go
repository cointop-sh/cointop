package coinmarketcap

// Interface interface
type Interface interface {
	GetGlobalMarketData() (GlobalMarketData, error)
	GetGlobalMarketGraphData(start int64, end int64) (MarketGraph, error)
	GetAltcoinMarketGraphData(start int64, end int64) (MarketGraph, error)
	GetCoinData(coin string) (Coin, error)
	GetAllCoinData(limit int) (map[string]Coin, error)
	GetCoinGraphData(coin string, start int64, end int64) (CoinGraph, error)
	GetCoinPriceUSD(coin string) (float64, error)
	GetCoinMarkets(coin string) ([]Market, error)
}

// Coin struct
type Coin struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Symbol           string  `json:"symbol"`
	Rank             int     `json:"rank,string"`
	PriceUSD         float64 `json:"price_usd,string"`
	PriceBTC         float64 `json:"price_btc,string"`
	USD24HVolume     float64 `json:"24h_volume_usd,string"`
	MarketCapUSD     float64 `json:"market_cap_usd,string"`
	AvailableSupply  float64 `json:"available_supply,string"`
	TotalSupply      float64 `json:"total_supply,string"`
	PercentChange1H  float64 `json:"percent_change_1h,string"`
	PercentChange24H float64 `json:"percent_change_24h,string"`
	PercentChange7D  float64 `json:"percent_change_7d,string"`
	LastUpdated      string  `json:"last_updated"`
}

// GlobalMarketData struct
type GlobalMarketData struct {
	TotalMarketCapUSD            float64 `json:"total_market_cap_usd"`
	Total24HVolumeUSD            float64 `json:"total_24h_volume_usd"`
	BitcoinPercentageOfMarketCap float64 `json:"bitcoin_percentage_of_market_cap"`
	ActiveCurrencies             int     `json:"active_currencies"`
	ActiveAssets                 int     `json:"active_assets"`
	ActiveMarkets                int     `json:"active_markets"`
}

// CoinGraph struct
type CoinGraph struct {
	MarketCapByAvailableSupply [][]float64 `json:"market_cap_by_available_supply"`
	PriceBTC                   [][]float64 `json:"price_btc"`
	PriceUSD                   [][]float64 `json:"price_usd"`
	VolumeUSD                  [][]float64 `json:"volume_usd"`
}

// Market struct
type Market struct {
	Rank          int
	Exchange      string
	Pair          string
	VolumeUSD     float64
	Price         float64
	VolumePercent float64
	Updated       string
}

// MarketGraph struct
type MarketGraph struct {
	MarketCapByAvailableSupply [][]float64 `json:"market_cap_by_available_supply"`
	VolumeUSD                  [][]float64 `json:"volume_usd"`
}
