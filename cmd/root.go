package cmd

import (
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"time"
	"path/filepath"
	"github.com/HotelsDotCom/flyte/httputil"
)

const (
	flagURL         = "url"
	flagFilename    = "filename"
	flagName        = "name"
	flagDescription = "description"
	flagContentType = "content-type"
	flagFormat      = "format"
	flagDslookup    = "ds-lookup"
)

var client = &http.Client{
	Timeout: time.Second * 5,
}

func newCmdFlyte() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "flyte",
		Short: "Command line client for flyte",
	}

	cmd.PersistentFlags().String(flagURL, "", "Flyte API URL. Overrides $FLYTE_API")
	viper.BindEnv(flagURL, "FLYTE_API")
	viper.BindPFlag(flagURL, cmd.PersistentFlags().Lookup(flagURL))
	cmd.AddCommand(
		newCmdTest(),
		newCmdUpload(),
		newCmdVersion(),
	)

	return cmd
}

func Execute() {
	if err := newCmdFlyte().Execute(); err != nil {
		os.Exit(1)
	}
}

func getContentType(filename string) string {
	switch filepath.Ext(filename) {
	case ".json":
		return httputil.MediaTypeJson
	case ".yaml", ".yml":
		return httputil.MediaTypeYaml
	case "sh":
		return "application/x-sh"
	default:
		return "application/octet-stream"
	}
}
