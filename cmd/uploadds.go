package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"net/http"
	"github.com/HotelsDotCom/flyte/flytepath"
	"io"
	"os"
	httputl "net/http/httputil"
	"github.com/spf13/viper"
	"bytes"
	"mime/multipart"
	"net/textproto"
	"strings"
	"path/filepath"
)

type dsItem struct {
	name        string
	description string
	contentType string
	filename    string
}

var argsUploadDs dsItem

func newCmdUploadDs() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "datastore -f FILENAME",
		Aliases: []string{"ds"},
		Short:   "Upload a datastore item from a file",
		Long:    longUploadDs,
		RunE:    runUploadDs,
	}

	cmd.Flags().StringVarP(&argsUploadDs.filename, flagFilename, "f", "", "filename of the file to use to upload resource")
	cmd.MarkFlagRequired(flagFilename)

	cmd.Flags().StringVarP(&argsUploadDs.name, flagName, "n", "", "item's name (default derived from the file name)")
	cmd.Flags().StringVarP(&argsUploadDs.description, flagDescription, "d", "", "item's description")
	cmd.Flags().StringVarP(&argsUploadDs.contentType, flagContentType, "c", "", "item's content type (default derived from the file extension)")

	return cmd
}

const longUploadDs = `
Upload a datastore item from a file or from stdin to a flyte API.
Flyte API could be specified by setting $FLYTE_API or overridden by the --url option

Examples:
  # Upload a datastore item from env.json file to flyte API specified by $FLYTE_API
  flyte upload ds -f ./env.json

  # Upload a datastore item from my-script.sh file to flyte API at http://127.0.0.1:8080
  flyte upload ds -f ./my-script.sh --url http://127.0.0.1:8080
`

func runUploadDs(c *cobra.Command, args []string) error {
	if argsUploadDs.name == "" {
		base := filepath.Base(argsUploadDs.filename)
		ext := filepath.Ext(argsUploadDs.filename)
		argsUploadDs.name = strings.TrimSuffix(base, ext)
	}

	if argsUploadDs.contentType == "" {
		argsUploadDs.contentType = getContentType(argsUploadDs.filename)
	}

	req, err := newDsRequest(viper.GetString(flagURL), argsUploadDs)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	dump, err := httputl.DumpResponse(resp, true)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("cannot upload datastore\n%s", dump)
	}

	_, err = fmt.Fprintf(c.OutOrStdout(), "%sLocation: %s\n", dump, resp.Request.URL.String())
	return err
}

func newDsRequest(apiURL string, item dsItem) (*http.Request, error) {
	file, err := os.Open(item.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)

	h := newFormFileHeader("value", filepath.Base(item.filename), item.contentType)
	part, err := w.CreatePart(h)
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(part, file); err != nil {
		return nil, err
	}

	if item.description != "" {
		w.WriteField("description", item.description)
	}

	if err = w.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, dsItemURL(apiURL, item.name), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	return req, nil
}

func newFormFileHeader(fieldname, filename, contentType string) textproto.MIMEHeader {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldname, filename))
	h.Set("Content-Type", contentType)
	return h
}

func dsItemURL(apiURL, name string) string {
	return fmt.Sprintf("%s%s/%s", apiURL, flytepath.DatastorePath, name)
}
