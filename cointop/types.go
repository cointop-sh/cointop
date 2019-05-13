package cointop

// Coin is the row structure
type Coin struct {
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
	// for favorites
	Favorite bool
	// for portfolio
	Holdings float64
	Balance  float64
}
