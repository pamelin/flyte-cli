package cmd

import (
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

const hostFlagName = "host"

var (
	rootCmd = &cobra.Command{
		Use:   "flyte",
		Short: "Command line client for flyte",
	}

	client = &http.Client{
		Timeout: time.Second * 5,
	}
)

func init() {
	rootCmd.PersistentFlags().String("host", "", "Flyte host address. Overrides $FLYTE_HOST")
	viper.BindEnv(hostFlagName, "FLYTE_HOST")
	viper.BindPFlag(hostFlagName, rootCmd.PersistentFlags().Lookup(hostFlagName))
	rootCmd.AddCommand(
		newTestCommand(),
		newUploadCommand(),
		newVersionCommand(),
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
