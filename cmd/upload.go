package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"errors"
	"github.com/spf13/viper"
	"net/http"
	"github.com/HotelsDotCom/flyte/flytepath"
	"path/filepath"
	"github.com/HotelsDotCom/flyte/httputil"
	"io"
	"os"
	httputl "net/http/httputil"
)

func newUploadCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "upload file",
		Short: "Upload flow from a file",
		Long:  uploadCmdLong,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("please provide a flow file")
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			filename := args[0]
			host := viper.GetString(hostFlagName)
			return runUploadCmd(filename, host, c.OutOrStdout())
		},
	}
}

func runUploadCmd(filename, host string, out io.Writer) error {

	ct, err := getContentType(filename)
	if err != nil {
		return err
	}

	body, err := os.Open(filename)
	if err != nil {
		return err
	}

	resp, err := client.Post(flowsURL(host), ct, body)
	if err != nil {
		return err
	}

	return processResponse(resp, out)
}

func getContentType(filename string) (string, error) {
	ext := filepath.Ext(filename)
	switch ext {
	case ".json":
		return httputil.MediaTypeJson, nil
	case ".yaml", ".yml":
		return httputil.MediaTypeYaml, nil
	default:
		return "", fmt.Errorf("cannot upload flow: unsupported file format %s", ext)
	}
}

// YEAH discovery will be in the next version
func flowsURL(host string) string {
	scheme := "http"
	return fmt.Sprintf("%s://%s%s", scheme, host, flytepath.FlowsPath)
}

func processResponse(resp *http.Response, out io.Writer) error {

	defer resp.Body.Close()
	dump, err := httputl.DumpResponse(resp, true)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("cannot upload flow\n%s", dump)
	}

	_, err = fmt.Fprintf(out, "%s", dump)
	return err
}

const uploadCmdLong = `

Upload flow from a file to a flyte host. File must be in JSON or YAML format.
Flyte host could be specified by setting $FLYTE_HOST or overridden by the --host option

Examples:
	# Upload a flow from my_flow.json file to flyte host specified by $FLYTE_HOST
	flyte upload ./my_flow.json

	# Upload a flow from my_flow.yaml file to flyte host at 127.0.0.1:8080
	flyte upload ./my_flow.yaml --host 127.0.0.1:8080
`
