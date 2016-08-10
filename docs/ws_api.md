Websocket API
---
A message from a client to a server called _operation call_.
Each operation call must contain at least operation name, and may contain identifier and call body.
The example of the operation call.

```json
{
    "operation" : "process.start",
    "id" : "id1234567",
    "body" : {
        "name" : "build_no_tests",
        "commandLine" : "mvn clean install -Dmaven.test.skip",
        "type" : "maven"
    }
}
```

- __operation__ - the name of the operation, usually dot separated resource and action
- __id__ - is an optional parameter, if client sends it then it is guaranteed 
that the result of the operation execution will contain the same identifier
- __body__ - operation execution body, data needed to perform operation

The example of the operation execution result:

```json
{
    "id" : "id1234567",
    "body" : {
        "pid" : 1,
        "name" : "build_no_tests",
        "commandLine" : "mvn clean install -Dmaven.test.skip",
        "type" : "maven"
    },
    "error" : null
}
```

The __error__ and __body__ are mutual exclusive, if the operation
can't be executed due to an error then the response will contain the error 
message and the error code. The example of the operation response with error:

```json
{
    "id" : 12345,
    "body" : null,
    "error" : {
        "code" : 10002,
        "message" : "Command line required"
    }
}
```


### Process API

#### Start process

##### Call

- __name__ - the name of the command
- __commandLine__ - command line to execute
- __type__(optional) - command type
- __eventTypes__(optional) - comma separated types of events which will be
 received by this channel. By default all the process events will be received.

```json
{
    "operation" : "process.start",
    "id" : "0x12345",
    "body" : {
        "name" : "build",
        "commandLine" : "mvn clean install",
        "type" : "maven",
        "eventTypes" : "stderr,stdout"
    }
}
```

##### Result

```json
{
    "id" : "0x12345",
    "body" : {
        "pid" : 123,
        "name" : "build",
        "commandLine" : "mvn clean install",
        "type" : "maven"
    },
    "error" : null
}
```

#### Kill process

##### Call

- __pid__ - the id of the process to kill

```json
{
    "operation" : "process.kill",
    "id" : "0x12345",
    "body" : {
        "pid" : 123
    }
}
```

##### Result

```json
{
    "id" : "0x12345",
    "body" : {
        "pid" : 123,
        "text" : "Successfully killed"
    },
    "error" : null
}
```

#### Subscribe to process events

##### Call

- __pid__ - the id of the process to subscribe to
- __eventTypes__(optional) - comma separated types of events which will be
received by this channel. By default all the process events will be received
- __after__(optional) - process logs which appeared after given time will
be republished to the channel. This parameter may be useful when reconnecting to the machine-agent

```json
{
    "operation" : "process.subscribe",
    "id" : "0x12345",
    "body" : {
        "pid" : 123,
        "eventTypes" : "stdout,stderr",
        "after" : "2016-07-26T09:36:44.920890113+03:00"
    }
}
```

##### Result

```json
{
    "id" : "0x12345",
    "body" : {
        "pid" : 123,
        "eventTypes" : "stdout,stderr",
        "text" : "Successfully subscribed"
    },
    "error" : null
}
```


#### Unsubscribe from process events

##### Call

- __pid__ - the id of the process to unsubscribe from

```json
{
    "operation" : "process.unsubscribe",
    "id" : "0x12345",
    "body" : {
        "pid" : 123
    }
}
```

##### Result

```json
{
    "id" : "0x12345",
    "body" : {
        "pid" : 123,
        "text" : "Successfully unsubscribed"
    }
}
```

#### Update process subscriber

##### Call

- __pid__ - the id of the process which subscriber should be updated
- __eventTypes__ - comma separated types of events which will be
received by this channel.

```json
{
    "operation" : "process.updateSubscriber",
    "id" : "0x12345",
    "body" : {
        "pid" : 123,
        "eventTypes": "process_status,stderr"
    }
}
```

##### Result

```json
{
    "id" : "0x12345",
    "body" : {
        "pid" : 123,
        "eventTypes": "process_status,stderr",
        "text" : "Subscriber successfully updated"
    }
}
```