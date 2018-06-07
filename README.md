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
Please specify the `FLYTE_HOST` environment variable which will be used to make flyte API calls. For example:
```
FLYTE_HOST=localhost:8080
```
This can be overridden/set by optional flag `--host`

## Use it
This is good place to start:
```
flyte [command]
```
The commands are:
```
help        Help about any command
test        Test step execution
upload      Upload flow from a file
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

### Upload command
Upload flow from a file to a flyte host. File must be in JSON or YAML format.
Flyte host could be specified by setting $FLYTE_HOST or overridden by the --host option.
Please refer to flyte documentation for file layout.

Examples:
```
	# Upload a flow from my_flow.json file to flyte host specified by $FLYTE_HOST env variable
	flyte upload ./my_flow.json

	# Upload a flow from my_flow.yaml file to flyte host at 127.0.0.1:8080
	flyte upload ./my_flow.yaml --host 127.0.0.1:8080
```

## Abuse it
Feel free to experiment and extend it by contributing back :relaxed:
