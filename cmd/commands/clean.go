package cmd

import (
	"fmt"

	"github.com/cointop-sh/cointop/cointop"
	"github.com/spf13/cobra"
)

// CleanCmd ...
func CleanCmd() *cobra.Command {
	var config string

	cleanCmd := &cobra.Command{
		Use:   "clean",
		Short: "Clear the cache",
		Long:  `The clean command clears the cache`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load config etc
			ct, err := cointop.NewCointop(&cointop.Config{
				ConfigFilepath: config,
			})
			if err != nil {
				return err
			}

			// NOTE: if clean command, clean but don't run cointop
			return cointop.Clean(&cointop.CleanConfig{
				Log:      true,
				CacheDir: ct.State.GetCacheDir(),
			})
		},
	}

	cleanCmd.Flags().StringVarP(&config, "config", "c", "", fmt.Sprintf("Config filepath. (default %s)", cointop.DefaultConfigFilepath))

	return cleanCmd
}
