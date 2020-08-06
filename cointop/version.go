package cointop

import (
	"fmt"
	"runtime/debug"
)

// version is the cointop version which will be populated by ldflags
var version string

// Version returns the cointop version
func (ct *Cointop) Version() string {
	return Version()
}

// Version returns cointop version
func Version() string {
	if version != "" {
		return version
	}

	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		return buildInfo.Main.Version
	}

	return "(devel)"
}

// PrintVersion prints the version
func PrintVersion() {
	fmt.Println(Version())
}
