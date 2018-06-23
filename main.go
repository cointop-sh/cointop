package main

import (
	"flag"
	"fmt"

	"github.com/miguelmota/cointop/cointop"
)

func main() {
	var v, ver, test bool
	flag.BoolVar(&v, "v", false, "Version")
	flag.BoolVar(&ver, "version", false, "Version")
	flag.BoolVar(&test, "test", false, "Run test")
	flag.Parse()
	if v || ver {
		fmt.Printf("cointop v%s", cointop.Version())
	} else if test {
		doTest()
	} else {
		cointop.New().Run()
	}
}

func doTest() {
	cointop.New().Exit()
}
