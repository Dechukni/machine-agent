TODO

- Versioning
- Logs lifecycle
- Logs paging(time based, limits based)
- Improve logging to file(see TODOs on top of logger file)
- Flags for configuration like(logs base folder, server host:port etc)
- Plug in terminal into the infrastructure
- When to cleanup not-alive processes
- Cleanup/Restore processes when machine restored from snapshot
- Websocket streaming strategy & logs/events model
- Add links to the REST API responses
- Validations for REST/WS received objects

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

* When the client successfully connected to the machine-agent(handshake response received)
```json
{
    "type" : "connected",
    "channelId" : "channel0x123456789",
    "text" : "Hello!",
    "time" : "2016-06-15T20:29:44.437650129+03:00"
}
```

* When a new process started(it is guaranteed that the "process_started"
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

* When the process died(normally died or killed)
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


* When the process produces output to the stdout
```json
{
    "type" : "stdout",
    "pid" : 123,
    "text" : "ping google.com",
    "time" : "2016-06-15T20:29:44.437650129+03:00"
}
```

* When the process produces output to the stderr
```json
{
    "type" : "stderr",
    "pid" : 123,
    "text" : "Absolute path to 'ifconfig' is '/sbin/ifconfig'",
    "time" : "2016-06-15T20:29:44.437650129+03:00"
}
```

TODO
