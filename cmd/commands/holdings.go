package cmd

import (
	"fmt"

	"github.com/cointop-sh/cointop/cointop"
	"github.com/spf13/cobra"
)

// HoldingsCmd ...
func HoldingsCmd() *cobra.Command {
	var help bool
	var total bool
	var percentChange24H bool
	var noCache bool
	var noHeader bool
	var config string
	var sortBy string
	var sortDesc bool
	var format = "table"
	var humanReadable bool
	var filter []string
	var cols []string
	var convert string
	var hideBalances bool

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

			if total || percentChange24H {
				if percentChange24H {
					return ct.PrintHoldings24HChange(&cointop.TablePrintOptions{
						HumanReadable: humanReadable,
						Format:        format,
						Filter:        filter,
						Convert:       convert,
					})
				}
				if total {
					return ct.PrintHoldingsTotal(&cointop.TablePrintOptions{
						HumanReadable: humanReadable,
						Format:        format,
						Filter:        filter,
						Convert:       convert,
					})
				}
			}

			return ct.PrintHoldingsTable(&cointop.TablePrintOptions{
				SortBy:        sortBy,
				SortDesc:      sortDesc,
				HumanReadable: humanReadable,
				Format:        format,
				Filter:        filter,
				Cols:          cols,
				Convert:       convert,
				NoHeader:      noHeader,
				HideBalances:  hideBalances,
			})
		},
	}

	holdingsCmd.Flags().BoolVarP(&help, "help", "", help, "Help for holdings")
	holdingsCmd.Flags().BoolVarP(&total, "total", "t", total, "Show portfolio total only")
	holdingsCmd.Flags().BoolVarP(&percentChange24H, "24h", "", percentChange24H, "Show portfolio 24H change only")
	holdingsCmd.Flags().BoolVarP(&noCache, "no-cache", "", noCache, "No cache")
	holdingsCmd.Flags().BoolVarP(&humanReadable, "human", "h", humanReadable, "Human readable output")
	holdingsCmd.Flags().BoolVarP(&noHeader, "no-header", "", noHeader, "Don't display header columns")
	holdingsCmd.Flags().BoolVarP(&hideBalances, "hide-balances", "", hideBalances, "Hide portfolio balances. Useful for when sharing screen or taking screenshotss")
	holdingsCmd.Flags().StringVarP(&config, "config", "c", "", fmt.Sprintf("Config filepath. (default %s)", cointop.DefaultConfigFilepath))
	holdingsCmd.Flags().StringVarP(&sortBy, "sort-by", "s", sortBy, `Sort by column. Options are "name", "symbol", "price", "holdings", "balance", "24h"`)
	holdingsCmd.Flags().BoolVarP(&sortDesc, "sort-desc", "d", sortDesc, "Sort in descending order")
	holdingsCmd.Flags().StringVarP(&format, "format", "", format, `Ouput format. Options are "table", "csv", "json"`)
	holdingsCmd.Flags().StringSliceVarP(&filter, "filter", "", filter, `Filter portfolio entries by coin name or symbol, comma separated without spaces. Example: "btc,eth,doge"`)
	holdingsCmd.Flags().StringSliceVarP(&cols, "cols", "", cols, `Filter portfolio columns, comma separated without spaces. Example: "symbol,holdings,balance"`)
	holdingsCmd.Flags().StringVarP(&convert, "convert", "f", convert, "The currency to convert to")

	return holdingsCmd
}
