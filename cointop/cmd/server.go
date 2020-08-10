package cmd

import (
	"fmt"
	"time"

	cssh "github.com/miguelmota/cointop/cointop/ssh"
	"github.com/spf13/cobra"
)

// ServerCmd ...
func ServerCmd() *cobra.Command {
	var port uint = 22
	var address string = "0.0.0.0"
	var idleTimeout uint = 60
	var executableBinary string = "cointop"
	var hostKeyFile string = cssh.DefaultHostKeyFile

	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Run cintop SSH Server",
		Long:  `Run cointop SSH server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			server := cssh.NewServer(&cssh.Config{
				Address:          address,
				Port:             port,
				IdleTimeout:      time.Duration(int(idleTimeout)) * time.Second,
				ExecutableBinary: executableBinary,
				HostKeyFile:      hostKeyFile,
			})

			fmt.Printf("Running SSH server on port %v\n", port)
			return server.ListenAndServe()
		},
	}

	serverCmd.Flags().UintVarP(&port, "port", "p", port, "Port")
	serverCmd.Flags().StringVarP(&address, "address", "a", address, "Address")
	serverCmd.Flags().UintVarP(&idleTimeout, "idle-timeout", "t", idleTimeout, "Idle timeout in seconds")
	serverCmd.Flags().StringVarP(&executableBinary, "binary", "b", executableBinary, "Executable binary path")
	serverCmd.Flags().StringVarP(&hostKeyFile, "host-key-file", "k", hostKeyFile, "Host key file")

	return serverCmd
}
