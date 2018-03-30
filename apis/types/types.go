package apis

// Coin struct
type Coin struct {
	ID               string
	Name             string
	Symbol           string
	Rank             int
	PriceUSD         float64
	PriceBTC         float64
	USD24HVolume     float64
	MarketCapUSD     float64
	AvailableSupply  float64
	TotalSupply      float64
	PercentChange1H  float64
	PercentChange24H float64
	PercentChange7D  float64
	LastUpdated      string
}

// GlobalMarketData struct
type GlobalMarketData struct {
	TotalMarketCapUSD            float64
	Total24HVolumeUSD            float64
	BitcoinPercentageOfMarketCap float64
	ActiveCurrencies             int
	ActiveAssets                 int
	ActiveMarkets                int
}

// CoinGraph struct
type CoinGraph struct {
	MarketCapByAvailableSupply [][]float64
	PriceBTC                   [][]float64
	PriceUSD                   [][]float64
	VolumeUSD                  [][]float64
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
	MarketCapByAvailableSupply [][]float64
	VolumeUSD                  [][]float64
}
