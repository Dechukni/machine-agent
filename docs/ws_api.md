
Process API
---

### Start process

```json
{
    "operation" : "process.start",
    "name" : "build",
    "commandLine" : "mvn clean install"
    "type" : "maven"
}
```
__types is optional__

### Kill process

```json
{
    "operation" : "process.kill",
    "pid" : 123
}
```

### Subscribe to process events

```json
{
    "operation" : "process.subscribe",
    "pid" : 123,
    "types" : "stdout,stderr"
}
```

__types__ is optional
