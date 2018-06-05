package cmd

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/spf13/cobra"
	"bytes"
	"github.com/stretchr/testify/assert"
)

func TestTestCommand_ShouldExecuteStepAndReturnOutputForJsonInput(t *testing.T) {
	output, err := executeCommand(rootCmd, "test", "testdata/step-test.json")
	require.NoError(t, err)

	assert.Equal(t, jsonOutput, output)
}

func TestTestCommand_ShouldExecuteStepAndReturnOutputForYamlInput(t *testing.T) {
	output, err := executeCommand(rootCmd, "test", "testdata/step-test.yaml")
	require.NoError(t, err)

	assert.Equal(t, jsonOutput, output)
}

func TestTestCommand_ShouldExecuteStepAndReturnOutputForYmlInput(t *testing.T) {
	output, err := executeCommand(rootCmd, "test", "testdata/step-test.yml")
	require.NoError(t, err)

	assert.Equal(t, jsonOutput, output)
}

func TestTestCommand_ShouldExecuteStepAndReturnOutputAsYaml(t *testing.T) {
	format = "yaml"
	output, err := executeCommand(rootCmd, "test", "testdata/step-test.json")
	require.NoError(t, err)

	assert.Equal(t, yamlOutput, output)
}

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
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
}`

const yamlOutput = `context:
  ChannelID: "123"
  UserID: johnny
input:
  channelId: "123"
  message: 'Hey <@johnny>, I''m up and running :run:'
name: SendMessage
packName: Slack
`