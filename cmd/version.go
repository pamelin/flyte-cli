package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	cliVersion = "v0.2"
	apiVersion = "v1"
)

func newVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show the flyte cli version information",
		Run: func(cmd *cobra.Command, args []string) {
			host := viper.GetString(hostFlagName)
			fmt.Printf("Client version:\t%s\nAPI version:\t%s\nAPI host:\t%s\n", cliVersion, apiVersion, host)
		},
	}
	return cmd
}
