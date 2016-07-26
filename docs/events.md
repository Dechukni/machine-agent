Events
===
Messages sent via websocket connections to clients

Process Events
---

#### STDERR event

Published when process writes to stderr.
One stderr event describes one output line

```json
{
    "type" : "stderr",
    "pid" : 123,
    "text" : "Absolute path to 'ifconfig' is '/sbin/ifconfig'",
    "time" : "2016-06-15T20:29:44.437650129+03:00"
}
```

#### STDOUT event

Published when process writes to stdout.
One stdout event describes one output line

```json
{
    "type" : "stdout",
    "pid" : 123,
    "text" : "Starting server...",
    "time" : "2016-06-15T20:29:44.437650129+03:00"
}
```

#### Process started

Published when process is successfully started.
This is the first event from all the events produced by process,
it appears only once for one process

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

#### Process died

Published when process is done, or killed. This is the last event from the process,
it appears only once for one process

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

Channel Events
---

#### Connected

The first event in the channel, published when client successfully connected to the machine-agent.

```json
{
    "type" : "connected",
    "channel" : "channel-123",
    "text" : "Hello!",
    "time" : "2016-06-15T20:29:44.437650129+03:00"
}
```


Error event
---

Published when any error occurred during Call processing

```json
{
    "type" : "error",
    "time" : "2016-06-15T20:29:44.437650129+03:00"
    "message" : "No process with id '123''",
}
```
