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
	"strings"
	"github.com/HotelsDotCom/flyte/httputil"
	"github.com/spf13/viper"
)

var argsTest = struct {
	filename string
	dsLookup bool
	format   string
}{}

func newCmdTest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test -f FILENAME",
		Short: "Test step execution with trigger event and optional context",
		Long:  longTest,
		RunE:  runTest,
	}

	cmd.Flags().StringVarP(&argsTest.filename, flagFilename, "f", "", "filename of the file with step and test data")
	cmd.MarkFlagRequired(flagFilename)

	cmd.Flags().BoolVar(&argsTest.dsLookup, flagDslookup, true, "lookup datastore item in the flyte API unless present in test data")
	cmd.Flags().StringVar(&argsTest.format, flagFormat, "json", "Output format. One of: json|yaml")
	return cmd
}

const longTest = `
Executes the step in the provided file. Test files MUST contain the
step, and trigger event definitions, and can optionally contain context and datastore
items as required.

Examples:
  # Test a step from my_step.yaml file
  flyte test -f ./my_step.yaml

and the yaml file can look like this:
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
testData:
  event:
    pack:
      name: Slack
    name: ReceivedMessage


You can run step test from stdin for example:
cat <<EOF | flyte test -f -
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
testData
  event:
    pack:
      name: Slack
    name: ReceivedMessage
EOF
`

func runTest(c *cobra.Command, args []string) error {
	var step testStep
	if err := unmarshalFile(argsTest.filename, &step); err != nil {
		return err
	}

	action, err := step.execute(argsTest.dsLookup, viper.GetString(flagURL))
	if err != nil {
		return err
	}

	out, err := marshal(action, argsTest.format)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(c.OutOrStdout(), string(out))
	return err
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

func (t testStep) execute(dsLookup bool, apiURL string) (*testAction, error) {
	//override flyte's default datastore function
	template.AddStaticContextEntry("datastore", datastoreFn(t.TestData.Datastore, dsLookup, apiURL))

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
	data, err := readFile(filename)
	if err != nil {
		return err
	}

	ext := detectExt(filename, data)

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

// datastore function which is backed by a map
// first search for item in the map, if not present try to lookup in the flyte API
func datastoreFn(datastore map[string]interface{}, dsLookup bool, apiURL string) func(string) interface{} {

	if datastore == nil {
		datastore = map[string]interface{}{}
	}
	return func(key string) interface{} {
		v, ok := datastore[key]
		if !ok {
			if dsLookup {
				v, err := findDatastoreItem(dsItemURL(apiURL, key))
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

func findDatastoreItem(url string) (interface{}, error) {

	resp, err := client.Get(url)
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
