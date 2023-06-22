package cmd

import (
	"github.com/cointop-sh/cointop/cointop"
	"github.com/spf13/cobra"
)

// VersionCmd ...
func VersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Displays the current version",
		Long:  `The version command displays the current version`,
		Run: func(cmd *cobra.Command, args []string) {
			cointop.PrintVersion()
		},
	}

	return versionCmd
}
