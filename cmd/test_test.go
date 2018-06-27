package cmd

import (
	"testing"
	"github.com/stretchr/testify/require"
	"bytes"
	"github.com/stretchr/testify/assert"
	"fmt"
	"net/http/httptest"
	"net/http"
	"github.com/HotelsDotCom/flyte/httputil"
	"github.com/HotelsDotCom/flyte/flytepath"
)

func TestTestCommand_ShouldExecuteStepAndReturnOutputForJsonInput(t *testing.T) {
	output, err := executeCommand("test", "-f", "testdata/step-test.json")
	require.NoError(t, err)

	assert.Equal(t, jsonOutput, output)
}

func TestTestCommand_ShouldExecuteStepAndReturnOutputForYamlInput(t *testing.T) {
	output, err := executeCommand("test", "-f", "testdata/step-test.yaml")
	require.NoError(t, err)

	assert.Equal(t, jsonOutput, output)
}

func TestTestCommand_ShouldExecuteStepAndReturnOutputForYmlInput(t *testing.T) {
	output, err := executeCommand("test", "-f", "testdata/step-test.yml")
	require.NoError(t, err)

	assert.Equal(t, jsonOutput, output)
}

func TestTestCommand_ShouldExecuteStepAndReturnOutputAsYaml(t *testing.T) {
	output, err := executeCommand("test", "-f", "testdata/step-test.json", "--format", "yaml")
	require.NoError(t, err)

	assert.Equal(t, yamlOutput, output)
}

func TestTestCommand_ShouldLookupDataItemInTheFlyteAPI(t *testing.T) {
	rec := requestRec{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.request = *r
		w.Header().Set(httputil.HeaderContentType, httputil.MediaTypeJson)
		fmt.Fprint(w, `{"flyte":{"status":"All good!!!"}}`)
	}))
	defer ts.Close()

	output, err := executeCommand("test", "-f", "testdata/step-ds.yaml", "--url", ts.URL)
	require.NoError(t, err)

	assert.Equal(t, flytepath.DatastorePath+"/env", rec.request.URL.String())
	assert.Equal(t, http.MethodGet, rec.request.Method)
	assert.Contains(t, output, "All good!!!")
}

func TestTestCommand_ShouldErrorWhenLookupDataItemFails(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	_, err := executeCommand("test", "-f", "testdata/step-ds.yaml", "--url", ts.URL)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot lookup datastore item key=env: invalid http response 404")
}

func TestTestCommand_ShouldSkipFlyteAPIDataItemLookup(t *testing.T) {
	_, err := executeCommand("test", "-f", "testdata/step-ds.yaml", "--ds-lookup=false")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot find datastore item key=env")
}

func TestTestCommand_ShouldLookupOtherTypesOfDataItemInTheFlyteAPI(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(httputil.HeaderContentType, "application/x-sh")
		fmt.Fprint(w, `echo hello`)
	}))
	defer ts.Close()

	output, err := executeCommand("test", "-f", "testdata/step-ds-non-json.yaml", "--url", ts.URL, "--ds-lookup=true")
	require.NoError(t, err)

	assert.Contains(t, output, "echo hello")
}

func executeCommand(args ...string) (output string, err error) {
	root := newCmdFlyte()
	buf := new(bytes.Buffer)
	root.SetOutput(buf)
	root.SetArgs(args)

	_, err = root.ExecuteC()
	return buf.String(), err
}

const jsonOutput = `{
	"name": "SendMessage",
	"packName": "Slack",
	"input": {
		"channelId": "123",
		"message": "Hey \u003c@johnny\u003e, I'm up and running :run:"
	},
	"context": {
		"ChannelID": "123",
		"UserID": "johnny"
	}
}
`

const yamlOutput = `context:
  ChannelID: "123"
  UserID: johnny
input:
  channelId: "123"
  message: 'Hey <@johnny>, I''m up and running :run:'
name: SendMessage
packName: Slack

`
