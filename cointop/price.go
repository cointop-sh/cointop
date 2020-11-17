package cointop

import (
	"fmt"

	"github.com/miguelmota/cointop/pkg/api"
	"github.com/miguelmota/cointop/pkg/humanize"
)

// PriceConfig is the config options for the price command
type PriceConfig struct {
	Coin      string
	Currency  string
	APIChoice string
}

// PrintPrice outputs the current price of the coin
func PrintPrice(config *PriceConfig) error {
	var priceAPI api.Interface
	if config.APIChoice == CoinMarketCap {
		priceAPI = api.NewCMC("")
	} else if config.APIChoice == CoinGecko {
		priceAPI = api.NewCG()
	} else {
		return ErrInvalidAPIChoice
	}

	price, err := priceAPI.Price(config.Coin, config.Currency)
	if err != nil {
		return err
	}

	symbol := CurrencySymbol(config.Currency)
	value := fmt.Sprintf("%s%s", symbol, humanize.Commaf(price))
	fmt.Println(value)

	return nil
}
