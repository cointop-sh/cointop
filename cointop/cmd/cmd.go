package cmd

// Execute executes the program
func Execute() {
	rootCmd := RootCmd()
	rootCmd.AddCommand(
		VersionCmd(),
		CleanCmd(),
		ResetCmd(),
		PriceCmd(),
		HoldingsCmd(),
		ServerCmd(),
		TestCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
