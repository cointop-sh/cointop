package cointop

import (
	"fmt"
	"runtime/debug"
	"strings"
)

// version is the cointop version which will be populated by ldflags
var version string

// Version returns the cointop version
func (ct *Cointop) Version() string {
	return Version()
}

// Version returns cointop version
func Version() string {
	ver := "(devel)"
	if version != "" {
		ver = version
	} else if buildInfo, ok := debug.ReadBuildInfo(); ok {
		ver = buildInfo.Main.Version
	}

	if !strings.HasPrefix(ver, "v") {
		ver = fmt.Sprintf("v%s", ver)
	}

	return ver
}

// PrintVersion prints the version
func PrintVersion() {
	fmt.Println(Version())
}
