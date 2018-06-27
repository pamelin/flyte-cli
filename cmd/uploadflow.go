package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"github.com/HotelsDotCom/flyte/flytepath"
	"os"
	httputl "net/http/httputil"
	"github.com/HotelsDotCom/flyte/httputil"
	"errors"
)

var argsUploadFlow = struct {
	filename    string
	contentType string
}{}

func newCmdUploadFlow() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "flow -f FILENAME",
		Short: "Upload a flow from a file or from stdin",
		Long:  longUploadFlow,
		RunE:  runUploadFlow,
	}

	cmd.Flags().StringVarP(&argsUploadFlow.filename, flagFilename, "f", "", "filename of the file to use to upload flow")
	cmd.MarkFlagRequired(flagFilename)

	cmd.Flags().StringVarP(&argsUploadFlow.contentType, flagContentType, "c", "", "flow file content type (default derived from the file extension)")

	return cmd
}

const longUploadFlow = `
Upload flow from a file to a flyte API. File must be in JSON or YAML format.
Flyte API could be specified by setting $FLYTE_API or overridden by the --url option

Examples:
  # Upload a flow from my_flow.json file to flyte api specified by $FLYTE_API
  flyte upload flow ./my_flow.json

  # Upload a flow from my_flow.yaml file to flyte api at http://127.0.0.1:8080
  flyte upload flow ./my_flow.yaml --url http://127.0.0.1:8080
`

func runUploadFlow(c *cobra.Command, args []string) error {

	if argsUploadFlow.contentType == "" {
		argsUploadFlow.contentType = getContentType(argsUploadFlow.filename)
	}
	if argsUploadFlow.contentType != httputil.MediaTypeJson &&
		argsUploadFlow.contentType != httputil.MediaTypeYaml {
		return errors.New("cannot upload flow: unsupported file type it must be JSON or YAML")
	}

	body, err := os.Open(argsUploadFlow.filename)
	if err != nil {
		return err
	}

	resp, err := client.Post(flowsURL(viper.GetString(flagURL)), argsUploadFlow.contentType, body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	dump, err := httputl.DumpResponse(resp, true)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("cannot upload flow\n%s", dump)
	}

	_, err = fmt.Fprintf(c.OutOrStdout(), "%s", dump)
	return err
}

func flowsURL(apiURL string) string {
	return fmt.Sprintf("%s%s", apiURL, flytepath.FlowsPath)
}
