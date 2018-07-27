package cmd

import (
	"github.com/spf13/cobra"
)

func newCmdUpload() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload TYPE",
		Short: "Upload a resource from a file or from stdin",
		Long:    longUpload,
		Example: exampleUpload,
	}

	cmd.SetUsageTemplate(usageTmplUpload)
	cmd.AddCommand(newCmdUploadFlow(), newCmdUploadDs())
	return cmd
}

const longUpload = `

Upload a resource from a file or from stdin to a flyte API. Valid resource types include:

  * datastore (aka 'ds')
  * flow`

const exampleUpload = `  # Upload a flow from my_flow.json file to flyte API specified by $FLYTE_API
  flyte upload flow -f ./my_flow.json

  # Upload a flow from my_flow.yaml file to flyte API at http://127.0.0.1:8080
  flyte upload flow -f ./my_flow.yaml --url http://127.0.0.1:8080`

const usageTmplUpload = `Usage:
  {{.UseLine}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [TYPE] --help" for more information about a command.{{end}}
`
