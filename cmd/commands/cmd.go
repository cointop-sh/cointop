package cmd

import (
	"fmt"
	"os"
)

// Execute executes the program
func Execute() {
	rootCmd := RootCmd()
	rootCmd.AddCommand(
		VersionCmd(),
		CleanCmd(),
		ResetCmd(),
		HoldingsCmd(),
		PriceCmd(),
		DominanceCmd(),
		ServerCmd(),
		TestCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
