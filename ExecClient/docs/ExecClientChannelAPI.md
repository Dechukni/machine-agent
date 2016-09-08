#ExecAgentChannel class API#

This document specifies client-side channel API of ExecAgent for use in JS webapps, for general description of communication with ExecAgent as well as glossary, object's structures please refer to respective documentation

##Constructor:##

**channel(boolean openNow=true, boolean saveHistory=true)**

-   **openNow(optional):** whether to open channel on construction or later calling function open, by default equal to true
-   **saveHistory(optional):** whether to save history of subscribe/unsubscribe calls and handlers and try to restore the same state of channel upon soft reload (F5, etc.) or lose of internet connection (where it will try to reconnect automatically every 10 seconds, upon detecting connection loss), by default equal to true

Creates instance of channel to ExecAgent, optionally opening it in process of construction

##Utility:##

**open()**

This function will open channel to given address:port

**close(boolean carryHistory=true)**

-   **carryHistory(optional):** if true, on next call of open will refer to saved history and try to restore previous state of channel, by default equal true

This function will manually attempt to close connection to exec agent server, so it later can be reopened

**getChannelID()**

This function returns channel id, so it can be used during REST communications with exec-agent server

##Operation calls:##

All operation calls are represented as promises that will generate id based on operation call time and by generated id will verify response, resolving promise on response without error (returning body object of response), or rejecting on return message with error (returning error object of response)
All operation calls functions can be called by either specifying only non-optional arguments as arguments OR providing object that contain all non-optional arguments and some optional ones.

**process.start(string name, string cmd, string type='', string eventType='stdout, stderr', func handler=())**

-   **name:** nickname of process, will be present in return message body
-   **cmd:** cmd to be executed by exec agent
-   **type(optional):** type of the process
-   **eventType(optional):** event types on which channel will be subscribed upon successful process creation, by default equal to 'stdout, stderr'
-   **handler(optional):** function that will be executed for all subscribed events of this process received by this channel

This function will create new process on server, subscribing channel to this process events of specified types and optionally adding handler for subscribed events

**process.kill(int pid)**

-   pid: id of process to be killed by exec agent

This function will kill the process with given pid, unsubscribing channel from process events, and clearing handlers assigned to this process

**process.subscribe(int pid, string eventType='stdout, stderr', time after='', func handler=())**

-   **pid:** id of process that will be subscribed to this channel
-   **eventType(optional):** event types on which channel will be subscribed, by default equal to 'stdout, stderr'
-   **after(optional):** point of time from which channel will receive events for given process
-   **handler(optional):** function that will be executed for all subscribed events

This function will subscribe our channel for events of process of given id, of given type, and after some point of time, optionally adding handler for subscribed events

**process.unsubscribe(int pid)**

-   **pid:** id of process that will be unsubscribed from this channel

This function will unsubscribe process from channel with given id and remove all handlers of this process events

**process.upsubscribe(int pid, string eventType, func handler=())**

-   **pid:** id of process whose subscription will be updated
-   **eventType:** event types of new subscription
-   **handler(optional):** function that will be executed for events of updated subsription

This function will update subscription of channel of a given process

**process.getLogs(int pid, time from='', time till='', string format='json', int limit=50, int skip=0)**

-   **pid:** id of process, logs of which we will get
-   **from(optional):** time to get logs from
-   **till(optional):** time to get logs till
-   **format(optional):** the format of the response, default is json, possible values are: text, json
-   **limit(optional):** the limit of logs in result, the default value is *50*, logs are limited from the latest to the earliest
-   **skip(optional):** the logs to skip, default value is *0*

This function will be resolved with array of logs' entries

##Handler operations:##

Handler is generally one-argument function which will be executed every time a suitable event will come to our client, with event data object as function argument.

**all.startEventHandling()**

This function allows to resume handling of events if it was paused, it is automatically called during constructor

**all.stopEventHandling()**

This function stops handling of incoming events, without changing any handlers by itself. Operation call responses promises will continue to be resolved/rejected. Notice: events that are incoming during pause are not logged and will not be processed by handlers after resuming

