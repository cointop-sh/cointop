package cmd

import (
	"github.com/cointop-sh/cointop/cointop"
	"github.com/cointop-sh/cointop/pkg/filecache"
	"github.com/spf13/cobra"
)

// ResetCmd ...
func ResetCmd() *cobra.Command {
	cacheDir := filecache.DefaultCacheDir

	resetCmd := &cobra.Command{
		Use:   "reset",
		Short: "Resets the config and clear the cache",
		Long:  `The reset command resets the config and clears the cache`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// NOTE: if reset command, reset but don't run cointop
			return cointop.Reset(&cointop.ResetConfig{
				Log:      true,
				CacheDir: cacheDir,
			})
		},
	}

	resetCmd.Flags().StringVarP(&cacheDir, "cache-dir", "", cacheDir, "Cache directory")

	return resetCmd
}
