//+build !windows

package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	cssh "github.com/cointop-sh/cointop/pkg/ssh"
	"github.com/spf13/cobra"
)

// ServerCmd ...
func ServerCmd() *cobra.Command {
	var port uint = 22
	address := "0.0.0.0"
	var idleTimeout uint = 0
	var maxTimeout uint = 0
	var maxSessions uint = 0
	var executableBinary = "cointop"
	hostKeyFile := cssh.DefaultHostKeyFile
	userConfigType := cssh.UserConfigTypePublicKey
	colorsDir := os.Getenv("COINTOP_COLORS_DIR")

	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Run cointop SSH Server",
		Long:  `Run cointop SSH server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			server := cssh.NewServer(&cssh.Config{
				Address:          address,
				Port:             port,
				IdleTimeout:      time.Duration(int(idleTimeout)) * time.Second,
				MaxTimeout:       time.Duration(int(maxTimeout)) * time.Second,
				MaxSessions:      maxSessions,
				ExecutableBinary: executableBinary,
				HostKeyFile:      hostKeyFile,
				UserConfigType:   userConfigType,
				ColorsDir:        colorsDir,
			})

			fmt.Printf("Running SSH server on port %v\n", port)
			return server.ListenAndServe()
		},
	}

	serverCmd.Flags().UintVarP(&port, "port", "p", port, "Port")
	serverCmd.Flags().StringVarP(&address, "address", "a", address, "Address")
	serverCmd.Flags().UintVarP(&idleTimeout, "idle-timeout", "t", idleTimeout, "Idle timeout in seconds. Default is 0 for no idle timeout")
	serverCmd.Flags().UintVarP(&maxTimeout, "max-timeout", "m", maxTimeout, "Max timeout in seconds. Default is 0 for no max timeout")
	serverCmd.Flags().UintVarP(&maxSessions, "max-sessions", "", maxSessions, "Max number of sessions allowed. Default is 0 for unlimited.")
	serverCmd.Flags().StringVarP(&executableBinary, "binary", "b", executableBinary, "Executable binary path")
	serverCmd.Flags().StringVarP(&hostKeyFile, "host-key-file", "k", hostKeyFile, "Host key file")
	serverCmd.Flags().StringVarP(&userConfigType, "user-config-type", "", userConfigType, fmt.Sprintf("User config type. Options are: %s", strings.Join(cssh.UserConfigTypes, ",")))

	return serverCmd
}
