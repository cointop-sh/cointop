package cmd

import (
	"fmt"

	"github.com/miguelmota/cointop/cointop"
	"github.com/spf13/cobra"
)

// HoldingsCmd ...
func HoldingsCmd() *cobra.Command {
	var help bool
	var total bool
	var noCache bool
	var config string
	var sortBy string
	var sortDesc bool
	var format string = "table"
	var humanReadable bool
	var filter []string

	holdingsCmd := &cobra.Command{
		Use:   "holdings",
		Short: "Displays current holdings",
		Long:  `The holdings command shows your current holdings`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if help {
				return cmd.Help()
			}

			ct, err := cointop.NewCointop(&cointop.Config{
				ConfigFilepath: config,
				CacheDir:       cointop.DefaultCacheDir,
			})
			if err != nil {
				return err
			}

			if total {
				return ct.PrintTotalHoldings(&cointop.TablePrintOptions{
					HumanReadable: humanReadable,
					Format:        format,
					Filter:        filter,
				})
			}

			return ct.PrintHoldingsTable(&cointop.TablePrintOptions{
				SortBy:        sortBy,
				SortDesc:      sortDesc,
				HumanReadable: humanReadable,
				Format:        format,
				Filter:        filter,
			})
		},
	}

	holdingsCmd.Flags().BoolVarP(&help, "help", "", help, "Help for holdings")
	holdingsCmd.Flags().BoolVarP(&total, "total", "t", total, "Show total only")
	holdingsCmd.Flags().BoolVarP(&noCache, "no-cache", "", noCache, "No cache")
	holdingsCmd.Flags().BoolVarP(&humanReadable, "human", "h", humanReadable, "Human readable output")
	holdingsCmd.Flags().StringVarP(&config, "config", "c", "", fmt.Sprintf("Config filepath. (default %s)", cointop.DefaultConfigFilepath))
	holdingsCmd.Flags().StringVarP(&sortBy, "sort-by", "s", sortBy, `Sort by column. Options are "name", "symbol", "price", "holdings", "balance", "24h"`)
	holdingsCmd.Flags().BoolVarP(&sortDesc, "sort-desc", "d", sortDesc, "Sort in descending order")
	holdingsCmd.Flags().StringVarP(&format, "format", "", format, `Ouput format. Options are "table", "csv", "json"`)
	holdingsCmd.Flags().StringSliceVarP(&filter, "filter", "f", filter, `Filter portfolio entries by coin name or symbol, comma separated. Example: "btc,eth,doge"`)

	return holdingsCmd
}
