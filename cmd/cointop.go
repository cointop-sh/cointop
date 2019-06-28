package cmd

import (
	"flag"
	"fmt"

	"github.com/miguelmota/cointop/cointop"
)

// Run ...
func Run() {
	var v, ver, test, clean, reset, hideMarketbar, hideChart, hideStatusbar, onlyTable bool
	var config, cmcAPIKey, apiChoice, colorscheme string
	flag.BoolVar(&v, "v", false, "Version")
	flag.BoolVar(&ver, "version", false, "Display current version")
	flag.BoolVar(&test, "test", false, "Run test (for Homebrew)")
	flag.BoolVar(&clean, "clean", false, "Wipe clean the cache")
	flag.BoolVar(&reset, "reset", false, "Reset the config. Make sure to backup any relevant changes first!")
	flag.BoolVar(&hideMarketbar, "hide-marketbar", false, "Hide the top marketbar")
	flag.BoolVar(&hideChart, "hide-chart", false, "Hide the chart view")
	flag.BoolVar(&hideStatusbar, "hide-statusbar", false, "Hide the bottom statusbar")
	flag.BoolVar(&onlyTable, "only-table", false, "Show only the table. Hides the chart and top and bottom bars")
	refreshRateFlag := flag.Int("refresh-rate", -1, "Refresh rate in seconds. Set to 0 to not auto-refresh. Default is 60")
	flag.StringVar(&config, "config", "", "Config filepath. Default is ~/.cointop/config.toml")
	flag.StringVar(&cmcAPIKey, "coinmarketcap-api-key", "", "Set the CoinMarketCap API key")
	flag.StringVar(&apiChoice, "api", cointop.CoinGecko, "API choice")
	flag.StringVar(&colorscheme, "colorscheme", "", "Colorscheme to use. Default is \"cointop\". To install standard themes, do:\n\ngit clone git@github.com:cointop-sh/colors.git ~/.cointop/colors\n\nFor additional instructions, visit: https://github.com/cointop-sh/colors")
	flag.Parse()

	var refreshRate *uint
	if refreshRateFlag != nil {
		if *refreshRateFlag > -1 {
			t := uint(*refreshRateFlag)
			refreshRate = &t
		}
	}

	if v || ver {
		fmt.Printf("cointop v%s", cointop.Version())
	} else if test {
		doTest()
	} else if clean {
		cointop.Clean()
	} else if reset {
		cointop.Reset()
	} else {
		cointop.NewCointop(&cointop.Config{
			ConfigFilepath:      config,
			CoinMarketCapAPIKey: cmcAPIKey,
			APIChoice:           apiChoice,
			Colorscheme:         colorscheme,
			HideMarketbar:       hideMarketbar,
			HideChart:           hideChart,
			HideStatusbar:       hideStatusbar,
			OnlyTable:           onlyTable,
			RefreshRate:         refreshRate,
		}).Run()
	}
}

func doTest() {
	cointop.NewCointop(&cointop.Config{
		NoPrompts: true,
	}).Exit()
}
