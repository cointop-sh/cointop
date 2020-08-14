package cmd

import (
	"fmt"

	"github.com/miguelmota/cointop/cointop"
	"github.com/spf13/cobra"
)

// RootCmd ...
func RootCmd() *cobra.Command {
	var version, test, clean, reset, hideMarketbar, hideChart, hideStatusbar, onlyTable, silent, noCache bool
	var refreshRate uint
	var config, cmcAPIKey, apiChoice, colorscheme string
	perPage := cointop.DefaultPerPage
	cacheDir := cointop.DefaultCacheDir

	rootCmd := &cobra.Command{
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

See git.io/cointop for more info.`,
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
				if err := cointop.Reset(&cointop.ResetConfig{
					Log: !silent,
				}); err != nil {
					return err
				}
			}

			// NOTE: if clean flag enabled, clean and run cointop
			if clean {
				if err := cointop.Clean(&cointop.CleanConfig{
					Log: !silent,
				}); err != nil {
					return err
				}
			}

			var refreshRateP *uint
			if cmd.Flags().Changed("refresh-rate") {
				refreshRateP = &refreshRate
			}

			ct, err := cointop.NewCointop(&cointop.Config{
				CacheDir:            cacheDir,
				NoCache:             noCache,
				ConfigFilepath:      config,
				CoinMarketCapAPIKey: cmcAPIKey,
				APIChoice:           apiChoice,
				Colorscheme:         colorscheme,
				HideMarketbar:       hideMarketbar,
				HideChart:           hideChart,
				HideStatusbar:       hideStatusbar,
				OnlyTable:           onlyTable,
				RefreshRate:         refreshRateP,
				PerPage:             perPage,
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
	rootCmd.Flags().BoolVarP(&silent, "silent", "s", false, "Silence log ouput")
	rootCmd.Flags().BoolVarP(&noCache, "no-cache", "", false, "No cache")
	rootCmd.Flags().UintVarP(&refreshRate, "refresh-rate", "r", 60, "Refresh rate in seconds. Set to 0 to not auto-refresh")
	rootCmd.Flags().UintVarP(&perPage, "per-page", "", perPage, "Per page")
	rootCmd.Flags().StringVarP(&config, "config", "c", "", fmt.Sprintf("Config filepath. (default %s)", cointop.DefaultConfigFilepath))
	rootCmd.Flags().StringVarP(&cmcAPIKey, "coinmarketcap-api-key", "", "", "Set the CoinMarketCap API key")
	rootCmd.Flags().StringVarP(&apiChoice, "api", "", "", "API choice. Available choices are \"coinmarketcap\" and \"coingecko\"")
	rootCmd.Flags().StringVarP(&colorscheme, "colorscheme", "", "", fmt.Sprintf("Colorscheme to use (default \"cointop\").\n%s", cointop.ColorschemeHelpString()))
	rootCmd.Flags().StringVarP(&cacheDir, "cache-dir", "", cacheDir, "Cache directory")

	return rootCmd
}
