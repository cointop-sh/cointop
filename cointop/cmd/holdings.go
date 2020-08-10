package cmd

import (
	"fmt"

	"github.com/miguelmota/cointop/cointop"
	"github.com/spf13/cobra"
)

// HoldingsCmd ...
func HoldingsCmd() *cobra.Command {
	var total, noCache bool
	var config string

	holdingsCmd := &cobra.Command{
		Use:   "holdings",
		Short: "Displays current holdings",
		Long:  `The holdings command shows your current holdings`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ct, err := cointop.NewCointop(&cointop.Config{
				ConfigFilepath: config,
				CacheDir:       cointop.DefaultCacheDir,
			})
			if err != nil {
				return err
			}

			if total {
				return ct.PrintTotalHoldings()
			}

			return ct.PrintHoldingsTable()
		},
	}

	holdingsCmd.Flags().BoolVarP(&total, "total", "t", false, "Show total only")
	holdingsCmd.Flags().BoolVarP(&noCache, "no-cache", "", false, "No cache")
	holdingsCmd.Flags().StringVarP(&config, "config", "c", "", fmt.Sprintf("Config filepath. (default %s)", cointop.DefaultConfigFilepath))

	return holdingsCmd
}
