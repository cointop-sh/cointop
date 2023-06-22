package cmd

import (
	"fmt"
	"os"

	"github.com/cointop-sh/cointop/cointop"
	"github.com/spf13/cobra"
)

// CleanCmd will wipe the cache only
func CleanCmd() *cobra.Command {
	config := os.Getenv("COINTOP_CONFIG")
	cacheDir := os.Getenv("COINTOP_CACHE_DIR")

	cleanCmd := &cobra.Command{
		Use:   "clean",
		Short: "Clear the cache",
		Long:  `The clean command clears the cache`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ct, err := cointop.NewCointop(&cointop.Config{
				ConfigFilepath: config,
			})
			if err != nil {
				return err
			}
			return ct.Clean(&cointop.CleanConfig{
				Log:      true,
				CacheDir: cacheDir,
			})
		},
	}

	cleanCmd.Flags().StringVarP(&cacheDir, "cache-dir", "", cacheDir, "Cache directory")
	cleanCmd.Flags().StringVarP(&config, "config", "c", config, fmt.Sprintf("Config filepath. (default %s)", cointop.DefaultConfigFilepath))

	return cleanCmd
}
