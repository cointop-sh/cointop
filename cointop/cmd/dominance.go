package cmd

import (
	"github.com/miguelmota/cointop/cointop"
	"github.com/spf13/cobra"
)

// DominanceCmd ...
func DominanceCmd() *cobra.Command {
	var apiChoice, currency string

	dominanceCmd := &cobra.Command{
		Use:   "dominance",
		Short: "Displays bitcoin dominance",
		Long:  `The dominance command display the dominance percentage of bitcoin`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cointop.PrintBitcoinDominance(&cointop.DominanceConfig{
				Currency:  currency,
				APIChoice: apiChoice,
			})
		},
	}

	dominanceCmd.Flags().StringVarP(&currency, "currency", "f", "USD", "The currency to convert to")
	dominanceCmd.Flags().StringVarP(&apiChoice, "api", "a", cointop.CoinGecko, "API choice. Available choices are \"coinmarketcap\" and \"coingecko\"")

	return dominanceCmd
}
