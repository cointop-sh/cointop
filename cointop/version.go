package cointop

import (
	"fmt"
	"os"
	"runtime/debug"
)

// version is the cointop version which will be populated by ldflags
var version string = "(unknown)"

// Version returns the cointop version
func (ct *Cointop) Version() string {
	return Version()
}

// Version returns cointop version
func Version() string {
	if version != "(unknown)" {
		return version
	}

	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		return buildInfo.Main.Version
	}

	return version
}

// PrintVersion prints the version
func PrintVersion() {
	fmt.Fprint(os.Stdout, Version())
}
