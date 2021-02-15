package cmd

import (
	"github.com/miguelmota/cointop/cointop"
	"github.com/spf13/cobra"
)

// PriceCmd ...
func PriceCmd() *cobra.Command {
	var apiChoice string
	var coins []string
	var currency string

	priceCmd := &cobra.Command{
		Use:   "price",
		Short: "Displays the current price of a coin",
		Long:  `The price command display the current price of a coin`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cointop.PrintPrices(&cointop.PricesConfig{
				Coins:     coins,
				Currency:  currency,
				APIChoice: apiChoice,
			})
		},
	}

	priceCmd.Flags().StringSliceVarP(&coins, "coins", "c", nil, "Name or symbol of coin(s), comma separated. E.g. \"Bitcoin\" Eg. \"btc,eth,doge\"")
	priceCmd.Flags().StringVarP(&currency, "currency", "f", "USD", "The currency to convert to")
	priceCmd.Flags().StringVarP(&apiChoice, "api", "a", cointop.CoinGecko, "API choice. Available choices are \"coinmarketcap\" and \"coingecko\"")

	return priceCmd
}
