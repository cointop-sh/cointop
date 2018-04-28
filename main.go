package main

import (
	"flag"
	"fmt"

	"github.com/miguelmota/cointop/cointop"
)

var version = "1.0.0"

func main() {
	var ver bool
	flag.BoolVar(&ver, "v", false, "Version")
	var test bool
	flag.BoolVar(&test, "test", false, "Run test")
	flag.Parse()
	if ver {
		fmt.Println(version)
		return
	}
	if test {
		runTest()
		return
	}

	cointop.Run()
}

func runTest() {
	// TODO
	fmt.Println("ok")
}
