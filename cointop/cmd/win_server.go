//+build windows

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func ServerCmd() *cobra.Command {
	return &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("//TODO: implement pty ssh server for Windows")
			return nil
		},
	}
}
