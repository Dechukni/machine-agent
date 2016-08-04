
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
    "body" : {
        "name" : "build",
        "commandLine" : "mvn clean install",
        "type" : "maven",
        "eventTypes" : "stderr,stdout"
    }
}
```

### Kill process

- __pid__ - the id of the process to kill

```json
{
    "operation" : "process.kill",
    "body" : {
        "pid" : 123
    }
}
```

### Subscribe to process events

- __pid__ - the id of the process to subscribe to
- __eventTypes__(optional) - comma separated types of events which will be
received by this channel. By default all the process events will be received.
- __after__(optional) - process logs which appeared after given time will
be republished to the channel. This parameter may be useful when reconnecting to the machine-agent

```json
{
    "operation" : "process.subscribe",
    "body" : {
        "pid" : 123,
        "eventTypes" : "stdout,stderr",
        "after" : "2016-07-26T09:36:44.920890113+03:00"
    }
}
```

### Unsubscribe from process events

- __pid__ - the id of the process to unsubscribe from

```json
{
    "operation" : "process.unsubscribe",
    "body" : {
        "pid" : 123
    }
}
```

### Update process subscriber

- __pid__ - the id of the process which subscriber should be updated
- __eventTypes__ - comma separated types of events which will be
received by this channel.

```json
{
    "operation" : "process.updateSubscriber",
    "body" : {
        "pid" : 123,
        "eventTypes": "process_status,stderr"
    }
}
```
