package cmd

import (
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "flyte",
	Short: "Command line client for flyte",
}

func init() {
	rootCmd.AddCommand(
		newTestCommand(),
		newVersionCommand(),
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}