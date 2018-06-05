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

## Use it
This is good place to start:
```
flyte [command]
```
The commands are:
```
help        Help about any command
test        Test step execution
version     Show the flyte version information
```

### Test command
Test command is executing step with provided test file containing step, trigger event and optional context and datastore items. It should be in json or yaml format.

Example:
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
event:
  pack:
    name: Slack
  event: ReceivedMessage
  payload:
    message: flyte status
    user:
      id: johny
context:
  ChannelID: '123'
datastore:
  message: 'I''m up and running :run:' 
```  
## Abuse it
Feel free to experiment and extend it by contributing back :relaxed:
