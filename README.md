# flyte-cli
`flyte-cli` is a command line client for [flyte](https://github.com/HotelsDotCom/flyte)

## Make it
Build:
```
$ make build
```
Build & Install:
```
$ make install
```

## Configure it
Please specify the `FLYTE_API` environment variable which will be used to make flyte API calls. For example:
```
FLYTE_API=http://localhost:8080
```
This can be overridden/set by optional flag `--url`

## Use it
This is good place to start:
```
flyte [command]
```
The commands are:
```
help        Help about any command
test        Test step execution
upload      Upload resource from a file
version     Show the flyte version information
```

### Test command
Executes the step in the provided file. Test files MUST contain the
step, and trigger event definitions, and can optionally contain context and datastore
items as required. It should be in json or yaml format.

Example yaml file:
```
step:
  id: status
  event:
    packName: Slack
    name: ReceivedMessage
  criteria: "{{ Event.Payload.message|match:'^flyte status$' }}"
  context:
    UserID: "{{ Event.Payload.user.id }}"
  command:
    packName: Slack
    name: SendMessage
    input:
      channelId: "{{ Context.ChannelID }}"
      message: 'Hey <@{{ Context.UserID }}>, {{datastore(''message'')}}'
testData:
  event:
    pack:
      name: Slack
    name: ReceivedMessage
    payload:
      message: flyte status
      user:
        id: johnny
  context:
    ChannelID: '123'
  datastore:
    message: 'I''m up and running :run:'
```
#### What is this datastore stuff?
By default test will try to find datastore items in the test data however if it is not available it will try to lookup
items in the flyte API. You can turn off lookup by passing `--ds-lookup=false` flag.

#### Upload command
Upload a resource from a file or from stdin to a flyte API. Valid resource types include:

  * datastore (aka ds)
  * flow

#### Upload flow command
Upload flow from a file or from stdin to a flyte API. File must be in JSON or YAML format.
Flyte API could be specified by setting $FLYTE_API or overridden by the --url option.
Please refer to flyte documentation for file layout.

Examples:
```
	# Upload a flow from my_flow.json file to flyte API specified by $FLYTE_API env variable
	flyte upload flow -f ./my_flow.json

	# Upload a flow from my_flow.yaml file to flyte API at http://127.0.0.1:8080
	flyte upload flow -f ./my_flow.yaml --url http://127.0.0.1:8080
```

#### Upload datastore (aka ds) item command
Upload a datastore item from a file or from stdin to a flyte API.
Flyte API could be specified by setting $FLYTE_API or overridden by the --url option

Examples:
```
	# Upload a datastore item from env.json file to flyte API specified by $FLYTE_API
	flyte upload ds -f ./env.json

	# Upload a datastore item from my-script.sh file to flyte API at http://127.0.0.1:8080
	flyte upload ds -f ./my-script.sh --url http://127.0.0.1:8080
```
	
## Abuse it
Feel free to experiment and extend it by contributing back :relaxed:
