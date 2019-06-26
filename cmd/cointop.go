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
	flag.BoolVar(&ver, "version", false, "Version")
	flag.BoolVar(&test, "test", false, "Run test")
	flag.BoolVar(&clean, "clean", false, "Clean cache")
	flag.BoolVar(&reset, "reset", false, "Reset config")
	flag.BoolVar(&hideMarketbar, "hide-marketbar", false, "Hide marketbar")
	flag.BoolVar(&hideChart, "hide-chart", false, "Hide chart view")
	flag.BoolVar(&hideStatusbar, "hide-statusbar", false, "Hide statusbar")
	flag.BoolVar(&onlyTable, "only-table", false, "Show only the table")
	flag.StringVar(&config, "config", "", "Config filepath")
	flag.StringVar(&cmcAPIKey, "coinmarketcap-api-key", "", "CoinMarketCap API key")
	flag.StringVar(&apiChoice, "api", cointop.CoinGecko, "API choice")
	flag.StringVar(&colorscheme, "colorscheme", "", "Colorscheme name")
	flag.Parse()

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
		}).Run()
	}
}

func doTest() {
	cointop.NewCointop(&cointop.Config{
		NoPrompts: true,
	}).Exit()
}
