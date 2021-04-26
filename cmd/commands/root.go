package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/miguelmota/cointop/cointop"
	"github.com/spf13/cobra"
)

// RootCmd ...
func RootCmd() *cobra.Command {
	var version bool
	var refreshRate uint
	test := getEnvBool("COINTOP_TEST")
	clean := getEnvBool("COINTOP_CLEAN")
	reset := getEnvBool("COINTOP_RESET")
	hideMarketbar := getEnvBool("COINTOP_HIDE_MARKETBAR")
	hideChart := getEnvBool("COINTOP_HIDE_CHART")
	hideTable := getEnvBool("COINTOP_HIDE_TABLE")
	hideStatusbar := getEnvBool("COINTOP_HIDE_STATUSBAR")
	onlyTable := getEnvBool("COINTOP_ONLY_TABLE")
	onlyChart := getEnvBool("COINTOP_ONLY_CHART")
	silent := getEnvBool("COINTOP_SILENT")
	noCache := getEnvBool("COINTOP_NO_CACHE")
	colorscheme := os.Getenv("COINTOP_COLORSCHEME")
	cacheDir := os.Getenv("COINTOP_CACHE_DIR")
	colorsDir := os.Getenv("COINTOP_COLORS_DIR")
	config := os.Getenv("COINTOP_CONFIG")
	apiChoice := os.Getenv("COINTOP_API")
	cmcAPIKey := os.Getenv("CMC_PRO_API_KEY")
	perPage := cointop.DefaultPerPage
	maxPages := cointop.DefaultMaxPages

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
			if refreshRateP == nil {
				value, ok := getEnvInt("COINTOP_REFRESH_RATE")
				if ok {
					uv := uint(value)
					refreshRateP = &uv
				}
			}

			ct, err := cointop.NewCointop(&cointop.Config{
				CacheDir:            cacheDir,
				ColorsDir:           colorsDir,
				NoCache:             noCache,
				ConfigFilepath:      config,
				CoinMarketCapAPIKey: cmcAPIKey,
				APIChoice:           apiChoice,
				Colorscheme:         colorscheme,
				HideMarketbar:       hideMarketbar,
				HideChart:           hideChart,
				HideTable:           hideTable,
				HideStatusbar:       hideStatusbar,
				OnlyTable:           onlyTable,
				OnlyChart:           onlyChart,
				RefreshRate:         refreshRateP,
				PerPage:             perPage,
				MaxPages:            maxPages,
			})
			if err != nil {
				return err
			}

			return ct.Run()
		},
	}

	rootCmd.Flags().BoolVarP(&version, "version", "v", version, "Display current version")
	rootCmd.Flags().BoolVarP(&test, "test", "", test, "Run test (for Homebrew)")
	rootCmd.Flags().BoolVarP(&clean, "clean", "", clean, "Wipe clean the cache")
	rootCmd.Flags().BoolVarP(&reset, "reset", "", reset, "Reset the config. Make sure to backup any relevant changes first!")
	rootCmd.Flags().BoolVarP(&hideMarketbar, "hide-marketbar", "", hideMarketbar, "Hide the top marketbar")
	rootCmd.Flags().BoolVarP(&hideChart, "hide-chart", "", hideChart, "Hide the chart view")
	rootCmd.Flags().BoolVarP(&hideTable, "hide-table", "", hideTable, "Hide the table view")
	rootCmd.Flags().BoolVarP(&hideStatusbar, "hide-statusbar", "", hideStatusbar, "Hide the bottom statusbar")
	rootCmd.Flags().BoolVarP(&onlyTable, "only-table", "", onlyTable, "Show only the table. Hides the chart and top and bottom bars")
	rootCmd.Flags().BoolVarP(&onlyChart, "only-chart", "", onlyChart, "Show only the chart. Hides the table and top and bottom bars")
	rootCmd.Flags().BoolVarP(&silent, "silent", "s", silent, "Silence log ouput")
	rootCmd.Flags().BoolVarP(&noCache, "no-cache", "", noCache, "No cache")
	rootCmd.Flags().UintVarP(&refreshRate, "refresh-rate", "r", 60, "Refresh rate in seconds. Set to 0 to not auto-refresh")
	rootCmd.Flags().UintVarP(&perPage, "per-page", "", perPage, "Per page")
	rootCmd.Flags().UintVarP(&maxPages, "max-pages", "", maxPages, "Max number of pages")
	rootCmd.Flags().StringVarP(&config, "config", "c", config, fmt.Sprintf("Config filepath. (default %s)", cointop.DefaultConfigFilepath))
	rootCmd.Flags().StringVarP(&cmcAPIKey, "coinmarketcap-api-key", "", cmcAPIKey, "Set the CoinMarketCap API key")
	rootCmd.Flags().StringVarP(&apiChoice, "api", "", apiChoice, "API choice. Available choices are \"coinmarketcap\" and \"coingecko\"")
	rootCmd.Flags().StringVarP(&colorscheme, "colorscheme", "", colorscheme, fmt.Sprintf("Colorscheme to use (default \"cointop\").\n%s", cointop.ColorschemeHelpString()))
	rootCmd.Flags().StringVarP(&cacheDir, "cache-dir", "", cacheDir, fmt.Sprintf("Cache directory (default %s)", cointop.DefaultCacheDir))
	rootCmd.Flags().StringVarP(&colorsDir, "colors-dir", "", colorsDir, "Colorschemes directory")

	return rootCmd
}

func getEnvBool(key string) bool {
	val := os.Getenv(key)
	if val == "" {
		return false
	}

	v, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}

	return v
}

func getEnvInt(key string) (int, bool) {
	val := os.Getenv(key)
	if val == "" {
		return 0, false
	}

	v, err := strconv.Atoi(val)
	if err != nil {
		return 0, false
	}

	return v, true
}
