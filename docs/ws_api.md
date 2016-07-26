
Process API
---

### Start process

- __name__ - the name of the command
- __commandLine__ - command line to execute
- __type__(optional) - command type
- __eventTypes__(optional) - comma separated types of events which will be
 received by this channel. By default all the process events will be received.

```json

{
    "operation" : "process.start",
    "name" : "build",
    "commandLine" : "mvn clean install",
    "type" : "maven",
    "eventTypes" : "stderr,stdout"
}
```

### Kill process

- __pid__ - the id of the process to kill

```json
{
    "operation" : "process.kill",
    "pid" : 123
}
```

### Subscribe to process events

- __pid__ - the id of the process to subscribe to
- __eventTypes__(optional) - comma separated types of events which will be
received by this channel. By default all the process events will be received.

```json
{
    "operation" : "process.subscribe",
    "pid" : 123,
    "eventTypes" : "stdout,stderr"
}
```
