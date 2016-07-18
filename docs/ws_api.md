
Process API
---

### Start a new process

```json
{
    "operation" : "process.start",
    "name" : "build",
    "commandLine" : "mvn clean install"
}
```

### Kill an existing process

```json
{
    "operation" : "process.kill",
    "pid" : 123
}
```
