REST API
===


Process API
---

### Start a new process

#### Request

_POST /process_

```json
{
    "name" : "build",
    "commandLine" : "mvn clean install",
    "type" : "maven"
}
```

#### Response

```json
{
    "pid": 1,
    "name": "build",
    "commandLine": "mvn clean install",
    "type" : "maven",
    "alive": true,
    "nativePid": 9186,
}
```
- `200` if successfully started
- `400` if incoming data is not valid e.g. name is empty
- `500` if any other error occurs


### Get a process

#### Request

_GET /process/{pid}_

- `pid` - the id of the process to get

#### Response

```json
{
    "pid": 1,
    "name": "build",
    "commandLine": "mvn clean install",
    "type" : "maven",
    "alive": false,
    "nativePid": 9186,
}
```

- `200` if response contains requested process
- `400` if `pid` is not valid, unsigned int required
- `404` if there is no such process
- `500` if any other error occurs

### Kill a process

#### Request

_DELETE /process/{pid}_

- `pid` - the id of the process to kill

#### Response

```json
{
    "pid": 1,
    "name": "build",
    "commandLine": "mvn clean install",
    "type" : "maven",
    "alive": true,
    "nativePid": 9186,
}
```
- `200` if successfully killed
- `400` if `pid` is not valid, unsigned int required
- `404` if there is no such process
- `500` if any other error occurs

### Get process logs

#### Request

_GET /process/{pid}/logs_

- `pid` - the id of the process to get logs

#### Response

The result logs of the process with the command line `printf "Hello\nWorld\n"`
```text
[STDOUT] 2016-07-04 08:37:56.315082296 +0300 EEST 	 Hello
[STDOUT] 2016-07-04 08:37:56.315128242 +0300 EEST 	 World
```
- `200` if logs are successfully fetched
- `404` if there is no such process
- `500` if any other error occurs

### Get processes

#### Request

_GET /process_

- `all`(optional) - if `true` then all the processes including _dead_ ones will be returned(respecting paging ofc), otherwise if `all` is `false`, or not specified or invalid then only _alive_ processes will be returned

#### Response

The result of the request _GET /process?all=true_
```json
[
    {
        "pid": 1,
        "name": "build",
        "commandLine": "mvn clean install",
        "type" : "maven",
        "alive": true,
        "nativePid": 9186,
    },
    {
        "pid": 2,
        "name": "build",
        "commandLine": "printf \"Hello World\"",
        "alive": false,
        "nativePid": 9588
    }
]
```
- `200` if processes are successfully retrieved
- `500` if any error occurs

### Subscribe to the process events

#### Request

_POST /process/{pid}/events/{channel}_

- `pid` - the id of the process to subscribe to
- `channel` - the id of the webscoket channel which is subscriber
- `types` - the types of the events separated with comma e.g. `?types=stderr,stdout`

#### Response

- `200` if successfully subscribed
- `400` if any of the parameters is not valid
- `404` if there is no such process or channel
- `500` if any other error occurs

### Unsubscribe from the process events

#### Request

_DELETE /process/{pid}/events/{channel}_

- `pid` - the id of the process to unsubscribe from
- `channel` - the id of the webscoket channel which currenly subscribed
to the process events

#### Response

- `200` if successfully unsubsribed
- `400` if any of the parameters is not valid
- `404` if there is no such process or channel
- `500` if any other error occurs

### Update the process events subscriber

#### Request

_PUT /process/{pid}/events/{channel}_

- `pid` - the id of the process
- `channel` - the id of the webscoket channel which is subscriber
- `types` - the types of the events separated with comma e.g. `?types=stderr,stdout`

#### Response

- `200` if successfully updated
- `400` if any of the parameters is not valid
- `404` if there is no such process or channel
- `500` if any other error occurs
