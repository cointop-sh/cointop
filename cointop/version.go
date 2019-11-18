package cointop

import (
	"fmt"
	"os"
)

// version is the cointop version which will be populated by ldflags
var version = "dev"

// Version returns the cointop version
func (ct *Cointop) Version() string {
	return version
}

// Version returns cointop version
func Version() string {
	return version
}

// PrintVersion prints the version
func PrintVersion() {
	fmt.Fprint(os.Stdout, Version())
}
