#ExecClient class API#

This document specifies use of client-side RESTful API of exec-agent for use in JS webapps, for general description of communication with exec-agent as well as glossary, object's structures please refer to respective documentation

##Constructor:##

**ExecClient(string address='http://localhost/', int port=9000, boolean saveHistory=true)**

-   **address(optional):** global address of our ExecAgent, default http://localhost/
-   **port(optional):** port of our ExecAgent, by default equal to 9000
-   **saveHistory(optional):** whether to save history of REST calls executed by this client, so it can be accessed later. Calls are represented as string that contains CURL command parameters

Creates instance of REST client with predefined address to our exec-agent

**ExecClient.channel(boolean openNow=true, boolean saveHistory=true)**

-	**openNow(optional):** whether to open channel on construction or later calling function open, by default equal to true
-	**saveHistory(optional):** whether to save history of subscribe/unsubscribe calls and handlers and try to restore the same state of channel upon soft reload (F5, etc.) or lose of internet connection (where it will	try to reconnect automatically every 10 seconds, upon detecting connection loss), by default equal to true

Creates instance of channel to exec-agent, optionally opening it in process of construction. For the channel API documentation please refer to [ExecClientChannelAPI.md] (ExecClientChannelAPI.md)

##Utility:##

**sendPrevious()**

This function sends last REST call once again, without any changes to parameters

**accessHistory(int callAmount=10)**

-   **callAmount(optional):** amount of last REST calls to be returned by functions

This function will return array of strings that contain representation of last callAmount calls

##REST communication##

All REST calls are represented as promises that will resolve(in case of 200) or reject(in case of other response code) with an object:

```json
{
	code: int,
	codeDescription: string,
	body: object
},
```

where code is a code of response, codeDescription is custom description of error, and body is a payload of response that is either a string, or js object (refer to REST communication doc for details of each call response body structure)

All REST functions can be called by either specifying only non-optional arguments as arguments OR providing object that contain all non-optional arguments and some optional ones.

**process.start(string name, string cmd, string type='', int channel, string eventType='')**

-   **name:** nickname of process, will be present in return message body
-   **cmd:** cmd to be executed by exec agent
-   **type(optional)**: type of the process
-   **channel(optional):** if this is specified, the mentioned channel will be subscribed to this process events
-   **eventType(optional):** event types on which channel will be subscribed (only viable when channel is passed), by default equal to 'stdout, stderr'

This function will send REST call that will result in creation of new process on server, optionally subscribing process events to channel

**process.kill(int pid)**

-   **pid:** id of process that will be killed

This function will send REST call that will result in killing the process with specified id

**process.info(int pid)**

-   **pid:** id of process about which we will obtain info

This function will send REST call that will return process information within response body

**process.getLogs(int pid, time from='', time till='', string format='json', int limit=50, int skip=0)**

-   **pid:** id of process, logs of which we will get
-   **from(optional):** time to get logs from
-   **till(optional):** time to get logs till
-   **format(optional):** the format of the response, default is json, possible values are: text, json
-   **limit(optional):** the limit of logs in result, the default value is *50*, logs are limited from the latest to the earliest
-   **skip(optional):** the logs to skip, default value is *0*

This function will send REST call that will return array of logs of certain process

**process.subscribe(int pid, int channel, string eventType='stdout, stderr', time after='')**

-   **pid:** id of process that will be subscribed to channel
-   **channel:** id of channel to which process will be subscribed
-   **eventType(optional):** event types on which channel will be subscribed, by default equal to 'stdout, stderr'
-   **after(optional):** point of time from which channel will receive events for given processed

This function will send REST call that will subscribe channel to process events

**process.unsubscribe(int pid, int channel)**

-   **pid:** id of process that will be unsubscribed from channel
-   **channel:** id of channel from which process will be unsubscribed

This function will send REST call that will unsubscribe process from channel

**process.upsubscribe(int pid, int channel, string eventType)**

-   **pid:** id of process whose subscription will be updated
-   **channel:** id of channel where subscription will be updated
-   **eventType:** event types of new subscription

This function will send REST call that will update subscription of process to specified event type at channel

**all.info(boolean all)**

-	**all(optional):** if true then all the processes including dead ones will be returned, otherwise if all is false, or not specified or invalid then only alive processes will be returned

This function will send REST call that will result in array with entries containing information about every process at server
