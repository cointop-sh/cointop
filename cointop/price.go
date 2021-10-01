package cointop

import (
	"fmt"
	"math"
	"strings"

	"github.com/cointop-sh/cointop/pkg/api"
	"github.com/cointop-sh/cointop/pkg/humanize"
)

// PriceConfig is the config options for the coin price method
type PriceConfig struct {
	Coin      string
	Currency  string
	APIChoice string
}

// PricesConfig is the config options for the coin prices method
type PricesConfig struct {
	Coins     []string
	Currency  string
	APIChoice string
}

// PrintPrices outputs the current price of the coins
func PrintPrices(config *PricesConfig) error {
	prices, err := GetCoinPrices(config)
	if err != nil {
		return err
	}
	fmt.Println(strings.Join(prices, "\n"))
	return nil
}

// PrintPrice outputs the current price of the coin
func PrintPrice(config *PriceConfig) error {
	prices, err := GetCoinPrices(&PricesConfig{
		Coins:     []string{config.Coin},
		Currency:  config.Currency,
		APIChoice: config.APIChoice,
	})
	if err != nil {
		return err
	}

	fmt.Println(prices[0])
	return nil
}

// GetCoinPrices returns the current price of the specified coins
func GetCoinPrices(config *PricesConfig) ([]string, error) {
	if len(config.Coins) == 0 {
		return nil, ErrCoinNameOrSymbolRequired
	}
	var priceAPI api.Interface
	if config.APIChoice == CoinMarketCap {
		priceAPI = api.NewCMC("")
	} else if config.APIChoice == CoinGecko {
		priceAPI = api.NewCG(0, 0)
	} else {
		return nil, ErrInvalidAPIChoice
	}

	var prices []string
	for _, coin := range config.Coins {
		price, err := priceAPI.Price(coin, config.Currency)
		if err != nil {
			return nil, err
		}

		symbol := CurrencySymbol(config.Currency)
		value := fmt.Sprintf("%s%s", symbol, humanize.Monetaryf(price, 2))
		prices = append(prices, value)
	}

	return prices, nil
}

// FormatPrice formats the coin price number of decimals and currency format
func (ct *Cointop) FormatPrice(price float64) string {
	decimals := 2
	if price < 1 {
		decimals = 8
	}
	if price == math.Trunc(price) {
		decimals = 2
	}
	return humanize.Monetaryf(price, decimals)
}
