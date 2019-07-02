package cmd

import (
	"log"
	"os"

	"github.com/miguelmota/cointop/cointop"
	"github.com/spf13/cobra"
)

// Execute executes the program
func Execute() {
	var version, test, clean, reset, hideMarketbar, hideChart, hideStatusbar, onlyTable bool
	var refreshRate uint
	var config, cmcAPIKey, apiChoice, colorscheme, coin, currency string

	var rootCmd = &cobra.Command{
		Use:   "cointop",
		Short: "Cointop is an interactive terminal based app for tracking cryptocurrencies",
		Long: `
           _       _
  ___ ___ (_)_ __ | |_ ___  _ __
 / __/ _ \| | '_ \| __/ _ \| '_ \
| (_| (_) | | | | | || (_) | |_) |
 \___\___/|_|_| |_|\__\___/| .__/
                           |_|

Cointop is a fast and lightweight interactive terminal based UI application for tracking and monitoring cryptocurrency coin stats in real-time.

For more information, visit: https://github.com/miguelmota/cointop`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if version {
				cointop.PrintVersion()
				return nil
			}

			if test {
				// TODO: deprecate test flag, only have test command
				doTest()
				return nil
			}

			// NOTE: if reset flag enabled, reset and run cointop
			if reset {
				if err := cointop.Reset(); err != nil {
					return err
				}
			}

			// NOTE: if clean flag enabled, clean and run cointop
			if clean {
				if err := cointop.Clean(); err != nil {
					return err
				}
			}

			var refreshRateP *uint
			if cmd.Flags().Changed("refresh-rate") {
				refreshRateP = &refreshRate
			}

			ct, err := cointop.NewCointop(&cointop.Config{
				ConfigFilepath:      config,
				CoinMarketCapAPIKey: cmcAPIKey,
				APIChoice:           apiChoice,
				Colorscheme:         colorscheme,
				HideMarketbar:       hideMarketbar,
				HideChart:           hideChart,
				HideStatusbar:       hideStatusbar,
				OnlyTable:           onlyTable,
				RefreshRate:         refreshRateP,
			})
			if err != nil {
				return err
			}

			return ct.Run()
		},
	}

	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "Display current version")
	rootCmd.Flags().BoolVarP(&test, "test", "", false, "Run test (for Homebrew)")
	rootCmd.Flags().BoolVarP(&clean, "clean", "", false, "Wipe clean the cache")
	rootCmd.Flags().BoolVarP(&reset, "reset", "", false, "Reset the config. Make sure to backup any relevant changes first!")
	rootCmd.Flags().BoolVarP(&hideMarketbar, "hide-marketbar", "", false, "Hide the top marketbar")
	rootCmd.Flags().BoolVarP(&hideChart, "hide-chart", "", false, "Hide the chart view")
	rootCmd.Flags().BoolVarP(&hideStatusbar, "hide-statusbar", "", false, "Hide the bottom statusbar")
	rootCmd.Flags().BoolVarP(&onlyTable, "only-table", "", false, "Show only the table. Hides the chart and top and bottom bars")
	rootCmd.Flags().UintVarP(&refreshRate, "refresh-rate", "r", 60, "Refresh rate in seconds. Set to 0 to not auto-refresh")
	rootCmd.Flags().StringVarP(&config, "config", "c", "", "Config filepath. (default ~/.cointop/config.toml)")
	rootCmd.Flags().StringVarP(&cmcAPIKey, "coinmarketcap-api-key", "", "", "Set the CoinMarketCap API key")
	rootCmd.Flags().StringVarP(&apiChoice, "api", "", cointop.CoinGecko, "API choice. Available choices are \"coinmarketcap\" and \"coingecko\"")
	rootCmd.Flags().StringVarP(&colorscheme, "colorscheme", "", "", "Colorscheme to use (default \"cointop\"). To install standard themes, do:\n\ngit clone git@github.com:cointop-sh/colors.git ~/.cointop/colors\n\nFor additional instructions, visit: https://github.com/cointop-sh/colors")

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Displays the current version",
		Long:  `The version command displays the current version`,
		Run: func(cmd *cobra.Command, args []string) {
			cointop.PrintVersion()
		},
	}

	var cleanCmd = &cobra.Command{
		Use:   "clean",
		Short: "Clear the cache",
		Long:  `The clean command clears the cache`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// NOTE: if clean command, clean but don't run cointop
			return cointop.Clean()
		},
	}

	var resetCmd = &cobra.Command{
		Use:   "reset",
		Short: "Resets the config and clear the cache",
		Long:  `The reset command resets the config and clears the cache`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// NOTE: if reset command, reset but don't run cointop
			return cointop.Reset()
		},
	}

	var priceCmd = &cobra.Command{
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

	var testCmd = &cobra.Command{
		Use:   "test",
		Short: "Runs tests",
		Long:  `The test command runs tests for Homebrew`,
		Run: func(cmd *cobra.Command, args []string) {
			doTest()
		},
	}

	priceCmd.Flags().StringVarP(&coin, "coin", "c", "bitcoin", "Full name of the coin (default \"bitcoin\")")
	priceCmd.Flags().StringVarP(&currency, "currency", "f", "USD", "The currency to convert to (default \"USD\")")
	priceCmd.Flags().StringVarP(&apiChoice, "api", "a", cointop.CoinGecko, "API choice. Available choices are \"coinmarketcap\" and \"coingecko\"")

	rootCmd.AddCommand(versionCmd, cleanCmd, resetCmd, priceCmd, testCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func doTest() {
	ct, err := cointop.NewCointop(&cointop.Config{
		NoPrompts: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	ct.Exit()
}
