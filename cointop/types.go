package cointop

type coin struct {
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
	Favorite         bool
}
