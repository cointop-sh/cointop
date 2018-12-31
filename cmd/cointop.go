package cmd

import (
	"flag"
	"fmt"

	"github.com/miguelmota/cointop/cointop"
)

// Run ...
func Run() {
	var v, ver, test, clean, reset bool
	var config string
	flag.BoolVar(&v, "v", false, "Version")
	flag.BoolVar(&ver, "version", false, "Version")
	flag.BoolVar(&test, "test", false, "Run test")
	flag.BoolVar(&clean, "clean", false, "Clean cache")
	flag.BoolVar(&reset, "reset", false, "Reset config")
	flag.StringVar(&config, "config", "", "Config filepath")
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
			ConfigFilepath: config,
		}).Run()
	}
}

func doTest() {
	cointop.NewCointop(nil).Exit()
}
