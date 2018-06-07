package cmd

import (
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

const hostFlagName = "host"

var client = &http.Client{
	Timeout: time.Second * 5,
}

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "flyte",
		Short: "Command line client for flyte",
	}

	cmd.PersistentFlags().String("host", "", "Flyte host address. Overrides $FLYTE_HOST")
	viper.BindEnv(hostFlagName, "FLYTE_HOST")
	viper.BindPFlag(hostFlagName, cmd.PersistentFlags().Lookup(hostFlagName))
	cmd.AddCommand(
		newTestCommand(),
		newUploadCommand(),
		newVersionCommand(),
	)

	return cmd
}

func init() {
}

func Execute() {

	if err := newRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
