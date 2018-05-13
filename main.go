package main

import (
	"flag"
	"fmt"

	"github.com/miguelmota/cointop/cointop"
)

var version = "1.0.1"

func main() {
	var ver, test bool
	flag.BoolVar(&ver, "v", false, "Version")
	flag.BoolVar(&test, "test", false, "Run test")
	flag.Parse()
	if ver {
		fmt.Println(version)
	} else if test {
		doTest()
	} else {
		cointop.New().Run()
	}
}

func doTest() {
	cointop.New().Exit()
}
