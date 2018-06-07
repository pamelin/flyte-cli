package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"github.com/HotelsDotCom/flyte/template"
	"github.com/HotelsDotCom/flyte/execution"
	jsont "github.com/HotelsDotCom/flyte/json"
	"io/ioutil"
	"encoding/json"
	"github.com/ghodss/yaml"
	"errors"
	"path/filepath"
)

var format string

func newTestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test file",
		Short: "Test step execution with trigger event and optional context",
		Long:  testCmdLong,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("you need to provide a test file")
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			output, err := runTestCmd(args[0], format)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(c.OutOrStdout(), output)
			return err
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "json", "Output format. One of: json|yaml")
	return cmd
}

func runTestCmd(testFilePath, format string) (string, error) {
	var step testStep
	if err := unmarshalFile(testFilePath, &step); err != nil {
		return "", err
	}

	action, err := step.execute()
	if err != nil {
		return "", err
	}

	out, err := marshal(action, format)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

type testStep struct {
	Step     execution.Step
	TestData testData
}

type testData struct {
	Event     event
	Context   map[string]string
	Datastore map[string]interface{}
}

// not sure why execution.Event replaces name with json tag event
// this is here so we can use `name` in the files instead of `event`
type event struct {
	Name    string         `json:"name"`
	Pack    execution.Pack `json:"pack"`
	Payload jsont.Json     `json:"payload,omitempty"`
}

func (t testStep) execute() (*testAction, error) {
	//override default datastore func which is using mongo to get data item
	//use static map instead which can be passed in the input file
	template.AddStaticContextEntry("datastore", datastoreFn(t.TestData.Datastore))

	e := execution.Event{
		Pack:    t.TestData.Event.Pack,
		Name:    t.TestData.Event.Name,
		Payload: t.TestData.Event.Payload,
	}

	action, err := t.Step.Execute(e, t.TestData.Context)
	if err != nil {
		return nil, err
	}
	if action == nil {
		return nil, nil
	}
	return &testAction{
		Name:       action.Name,
		PackName:   action.PackName,
		PackLabels: action.PackLabels,
		Input:      action.Input,
		Context:    action.Context,
	}, nil
}

type testAction struct {
	Name       string            `json:"name"`
	PackName   string            `json:"packName"`
	PackLabels map[string]string `json:"packLabels,omitempty"`
	Input      jsont.Json        `json:"input,omitempty"`
	Context    map[string]string `json:"context,omitempty"`
}

func unmarshalFile(filename string, v interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	ext := filepath.Ext(filename)
	switch ext {
	case ".json":
		return json.Unmarshal(data, v)
	case ".yaml", ".yml":
		return yaml.Unmarshal(data, v)
	default:
		return fmt.Errorf("cannot unmarshal: unsuported file %s", ext)
	}
}

func marshal(v interface{}, format string) ([]byte, error) {
	switch format {
	case "yaml":
		return yaml.Marshal(v)
	default:
		return json.MarshalIndent(v, "", "\t")
	}
}

func datastoreFn(datastore map[string]interface{}) func(string) interface{} {
	if datastore == nil {
		datastore = map[string]interface{}{}
	}
	return func(key string) interface{} {
		v, ok := datastore[key]
		if !ok {
			panic(fmt.Errorf("cannot find datastore item key=%s", key))
		}
		return v
	}
}

const testCmdLong = `

Executes the step in the provided file. Test files MUST contain the
step, and trigger event definitions, and can optionally contain context and datastore
items as required.

Example yaml file:
---
step:
  id: status
  event:
    packName: Slack
    name: ReceivedMessage
  command:
    packName: Slack
    name: SendMessage
    input:
      message: 'Hello'
event:
  pack:
    name: Slack
  name: ReceivedMessage
`
