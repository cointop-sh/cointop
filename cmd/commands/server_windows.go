//+build windows

package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// ServerCmd ...
// TODO: implement pty ssh server for Windows
func ServerCmd() *cobra.Command {
	return &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Not implemented")
		},
	}
}
