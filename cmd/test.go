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
	"github.com/HotelsDotCom/flyte/flytepath"
	"strings"
	"github.com/HotelsDotCom/flyte/httputil"
	"github.com/spf13/viper"
)

var (
	dsLookup bool
	format   string
)

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
			output, err := runTestCmd(args[0])
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(c.OutOrStdout(), output)
			return err
		},
	}

	cmd.Flags().BoolVar(&dsLookup, "ds-lookup", true, "lookup datastore item in the flyte host if not present in test data")
	cmd.Flags().StringVarP(&format, "format", "f", "json", "Output format. One of: json|yaml")
	return cmd
}

func runTestCmd(testFilePath string) (string, error) {
	var step testStep
	if err := unmarshalFile(testFilePath, &step); err != nil {
		return "", err
	}

	action, err := step.execute()
	if err != nil {
		return "", err
	}

	out, err := marshal(action)
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
	//override flyte's default datastore function
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

func marshal(v interface{}) ([]byte, error) {
	switch format {
	case "yaml":
		return yaml.Marshal(v)
	default:
		return json.MarshalIndent(v, "", "\t")
	}
}

// datastore function which is backed by a map
// first search for item in the map, if not present try to lookup in the flyte host
func datastoreFn(datastore map[string]interface{}) func(string) interface{} {

	if datastore == nil {
		datastore = map[string]interface{}{}
	}
	return func(key string) interface{} {
		v, ok := datastore[key]
		if !ok {
			if dsLookup {
				v, err := findDatastoreItem(key)
				if err != nil {
					panic(fmt.Errorf("cannot lookup datastore item key=%s: %v", key, err))
				}
				return v
			} else {
				panic(fmt.Errorf("cannot find datastore item key=%s", key))
			}
		}
		return v
	}
}

func findDatastoreItem(key string) (interface{}, error) {

	resp, err := client.Get(datastoreItemURL(key))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid http response %d %s", resp.StatusCode, resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return unmarshalValue(b, resp.Header.Get(httputil.HeaderContentType))
}

// parse data into the expected struct
// based on the behaviour of `github.com/HotelsDotCom/flyte/datastore.GetDataStoreValue`
func unmarshalValue(b []byte, contentType string) (interface{}, error) {

	if !strings.HasPrefix(contentType, httputil.MediaTypeJson) && !strings.HasPrefix(contentType, "text/json") {
		return string(b), nil
	}

	value := map[string]interface{}{}
	if err := json.Unmarshal(b, &value); err != nil {
		return nil, err
	}
	return value, nil
}

func datastoreItemURL(key string) string {
	scheme := "http"
	return fmt.Sprintf("%s://%s%s/%s", scheme, viper.GetString(hostFlagName), flytepath.DatastorePath, key)
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
