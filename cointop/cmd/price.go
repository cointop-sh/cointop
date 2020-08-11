package cmd

import (
	"github.com/miguelmota/cointop/cointop"
	"github.com/spf13/cobra"
)

// PriceCmd ...
func PriceCmd() *cobra.Command {
	var apiChoice, coin, currency string

	priceCmd := &cobra.Command{
		Use:   "price",
		Short: "Displays the current price of a coin",
		Long:  `The price command display the current price of a coin`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cointop.PrintPrice(&cointop.PriceConfig{
				Coin:      coin,
				Currency:  currency,
				APIChoice: apiChoice,
			})
		},
	}

	priceCmd.Flags().StringVarP(&coin, "coin", "c", "bitcoin", "Full name of the coin")
	priceCmd.Flags().StringVarP(&currency, "currency", "f", "USD", "The currency to convert to")
	priceCmd.Flags().StringVarP(&apiChoice, "api", "a", cointop.CoinGecko, "API choice. Available choices are \"coinmarketcap\" and \"coingecko\"")

	return priceCmd
}
