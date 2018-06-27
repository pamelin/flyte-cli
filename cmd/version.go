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

func newCmdVersion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show the flyte cli version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Client version:\t%s\nAPI version:\t%s\nAPI URL:\t%s\n",
				cliVersion, apiVersion, viper.GetString(flagURL))
		},
	}
	return cmd
}
