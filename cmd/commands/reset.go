package cmd

import (
	"fmt"
	"os"

	"github.com/cointop-sh/cointop/cointop"
	"github.com/spf13/cobra"
)

// ResetCmd will wipe cache and config file
func ResetCmd() *cobra.Command {
	config := os.Getenv("COINTOP_CONFIG")
	cacheDir := os.Getenv("COINTOP_CACHE_DIR")

	resetCmd := &cobra.Command{
		Use:   "reset",
		Short: "Resets the config and clear the cache",
		Long:  `The reset command resets the config and clears the cache`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ct, err := cointop.NewCointop(&cointop.Config{
				ConfigFilepath: config,
			})
			if err != nil {
				return err
			}
			return ct.Reset(&cointop.ResetConfig{
				Log:      true,
				CacheDir: cacheDir,
			})
		},
	}

	resetCmd.Flags().StringVarP(&cacheDir, "cache-dir", "", cacheDir, "Cache directory")
	resetCmd.Flags().StringVarP(&config, "config", "c", config, fmt.Sprintf("Config filepath. (default %s)", cointop.DefaultConfigFilepath))

	return resetCmd
}
