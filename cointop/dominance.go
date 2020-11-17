package cointop

import (
	"fmt"

	"github.com/miguelmota/cointop/pkg/api"
)

// DominanceConfig is the config options for the dominance command
type DominanceConfig struct {
	Currency  string
	APIChoice string
}

// PrintBitcoinDominance outputs the dominance percentage of bitcoin
func PrintBitcoinDominance(config *DominanceConfig) error {
	if config == nil {
		config = &DominanceConfig{}
	}

	var coinAPI api.Interface
	if config.APIChoice == CoinMarketCap {
		coinAPI = api.NewCMC("")
	} else if config.APIChoice == CoinGecko {
		coinAPI = api.NewCG()
	} else {
		return ErrInvalidAPIChoice
	}

	data, err := coinAPI.GetGlobalMarketData(config.Currency)
	if err != nil {
		return err
	}

	value := fmt.Sprintf("Bitcoin: %.2f%%", data.BitcoinPercentageOfMarketCap)
	fmt.Println(value)

	return nil
}
