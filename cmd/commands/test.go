package cmd

import (
	"github.com/cointop-sh/cointop/cointop"
	"github.com/spf13/cobra"
)

// TestCmd ...
func TestCmd() *cobra.Command {
	testCmd := &cobra.Command{
		Use:   "test",
		Short: "Runs tests for Homebrew",
		Long:  `The test command runs tests for Homebrew`,
		Run: func(cmd *cobra.Command, args []string) {
			doTest()
		},
	}

	return testCmd
}

// DoTest ...
func doTest() {
	ct, err := cointop.NewCointop(&cointop.Config{
		NoPrompts: true,
	})

	if err != nil {
		panic(err)
	}

	ct.Exit()
}
