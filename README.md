[TODO](TODO.md)

Websocket streaming strategy
===

Notes
---

##### Websocket messages order

The order is respected
```
Message fragments MUST be delivered to the recipient in the order sent by the sender.
```
Helpful Sources
* https://tools.ietf.org/html/rfc6455 (search the sentence above)
* http://stackoverflow.com/questions/11804721/can-websocket-messages-arrive-out-of-order
* http://stackoverflow.com/questions/14287224/processing-websockets-messages-in-order-of-receiving


Events(websocket messages which are sent by agent) model
---

All the events provide at least 2 fields:
* `type` - the type of the event
* `time` - the time related to the event(e.g. if an event type is 'stdout' then the `time`
field describes the moment when the message provided by this event was pumped from stdout).
For now the format of the field is `2016-06-15T20:29:44.437650129+03:00`

#### Event examples
 When the client successfully connected to the machine-agent(handshake response received)
```json
{
    "type" : "connected",
    "channelId" : "channel0x123456789",
    "text" : "Hello!",
    "time" : "2016-06-15T20:29:44.437650129+03:00"
}
```

When a new process started(it is guaranteed that the "process_started"
 event is sent before all the other messages related to the process like "stdout")
```json
{
    "type" : "process_started",
    "pid" : 123,
    "nativePid" : 22344,
    "name" : "build",
    "commandLine" : "mvn clean install",
    "time" : "2016-06-15T20:29:44.437650129+03:00"
}
```

When the process died(normally died or killed)
```json
{
    "type" : "process_died",
    "pid" : 123,
    "nativePid" : 22344,
    "name" : "build",
    "commandLine" : "mvn clean install",
    "time" : "2016-06-15T20:29:44.437650129+03:00"
}
```


When the process produces output to the stdout
```json
{
    "type" : "stdout",
    "pid" : 123,
    "text" : "ping google.com",
    "time" : "2016-06-15T20:29:44.437650129+03:00"
}
```

When the process produces output to the stderr
```json
{
    "type" : "stderr",
    "pid" : 123,
    "text" : "Absolute path to 'ifconfig' is '/sbin/ifconfig'",
    "time" : "2016-06-15T20:29:44.437650129+03:00"
}
```

ApiCall(websocket messages which are received by agent) model
---

All the ApiCalls provide at least operation type which usually dot separated
resource and action.
- `operation` - the operation which should be performed

#### ApiCall examples


**Start a new process**
```json
{
    "operation" : "process.start",
    "name" : "build",
    "commandLine" : "mvn clean install"
}
```

**Kill an existing process**
```json
{
    "operation" : "process.kill",
    "pid" : 123
}
```

or

```json
{
    "operation" : "process.kill",
    "nativePid" : 22388
}
```

or even, where pid is in preference to nativePid

```json
{
    "operation" : "process.kill",
    "pid" : 123,
    "nativePid" : 22388
}
```



