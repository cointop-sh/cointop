package cmd

import (
	"github.com/miguelmota/cointop/cointop"
	"github.com/miguelmota/cointop/pkg/filecache"
	"github.com/spf13/cobra"
)

// CleanCmd ...
func CleanCmd() *cobra.Command {
	cacheDir := filecache.DefaultCacheDir

	cleanCmd := &cobra.Command{
		Use:   "clean",
		Short: "Clear the cache",
		Long:  `The clean command clears the cache`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// NOTE: if clean command, clean but don't run cointop
			return cointop.Clean(&cointop.CleanConfig{
				Log:      true,
				CacheDir: cacheDir,
			})
		},
	}

	cleanCmd.Flags().StringVarP(&cacheDir, "cache-dir", "", cacheDir, "Cache directory")

	return cleanCmd
}
