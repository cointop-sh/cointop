package cmd

// Execute executes the program
func Execute() {
	rootCmd := RootCmd()
	rootCmd.AddCommand(
		VersionCmd(),
		CleanCmd(),
		ResetCmd(),
		HoldingsCmd(),
		PriceCmd(),
		DominanceCmd(),
		ServerCmd(),
		TestCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
