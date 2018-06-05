package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newVersionCommand() *cobra.Command{
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show the flyte cli version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Client version:\tv0.1\nAPI version:\tv1")
		},
	}
	return cmd
}
