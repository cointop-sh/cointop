package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// NameToSlug converts a coin name to slug for URLs
func NameToSlug(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	return strings.ToLower(re.ReplaceAllString(name, "-"))
}

// FormatID formats the ID value
func FormatID(id string) string {
	return strings.ToLower(id)
}

// FormatSymbol formats the symbol value
func FormatSymbol(id string) string {
	return strings.ToUpper(id)
}

// FormatName formats the name value
func FormatName(name string) string {
	return name
}

// FormatRank formats the rank value
func FormatRank(rank interface{}) int {
	switch v := rank.(type) {
	case int:
		return v
	case uint:
		return int(v)
	case int16:
		return int(v)
	case uint16:
		return int(v)
	case int32:
		return int(v)
	case uint32:
		return int(v)
	case int64:
		return int(v)
	case uint64:
		return int(v)
	case float32:
		return int(v)
	case float64:
		return int(v)
	}

	return 0
}

// FormatPrice formats the price value
func FormatPrice(price float64, convert string) float64 {
	convert = strings.ToUpper(convert)
	pricestr := fmt.Sprintf("%.2f", price)
	if convert == "ETH" || convert == "BTC" || price < 1 {
		pricestr = fmt.Sprintf("%.5f", price)
	}
	price, _ = strconv.ParseFloat(pricestr, 64)
	return price
}

// FormatVolume formats the volume value
func FormatVolume(volume float64) float64 {
	return float64(int64(volume))
}

// FormatMarketCap formats the market cap value
func FormatMarketCap(marketCap float64) float64 {
	return float64(int64(marketCap))
}

// FormatSupply formats the supply value
func FormatSupply(supply float64) float64 {
	return float64(int64(supply))
}

// FormatPercentChange formats the percent change value
func FormatPercentChange(percentChange float64) float64 {
	return percentChange
}

// FormatLastUpdated formats the last updated value
func FormatLastUpdated(lastUpdated string) string {
	lastUpdatedTime, err := time.Parse(time.RFC3339, lastUpdated)
	if err != nil {
		return ""
	}

	return strconv.Itoa(int(lastUpdatedTime.Unix()))
}

// CalcDays calculates the number of days between two timestamps
func CalcDays(start, end int64) int {
	return int(time.Unix(end, 0).Sub(time.Unix(start, 0)).Hours() / 24)
}
