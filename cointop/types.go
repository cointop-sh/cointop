package cointop

type coin struct {
	ID               string
	Name             string
	Slug             string
	Symbol           string
	Rank             int
	Price            float64
	Volume24H        float64
	MarketCap        float64
	AvailableSupply  float64
	TotalSupply      float64
	PercentChange1H  float64
	PercentChange24H float64
	PercentChange7D  float64
	LastUpdated      string
	Favorite         bool
}
