TODO

- Versioning
- Logs lifecycle
- Logs paging(time based, limits based)
- Flags for configuration like(logs base folder, server host:port etc)
- Plug in terminal into the infrastructure
- When to cleanup not-alive processes
- Cleanup/Restore processes when machine restored from snapshot
- Validations for REST/WS received objects
- Project structure
- Consider using RAML for API documentation

Create Dispatcher for handling api calls
---
- Defines methods for handling webscoket connections
- Manages webscoket connections state
- Exposes HttpRoutes for connections/management

Reconnect
---
Develop reconnect mechanism(e.g. remember processes which this channel subscribed to + disconnect time, and keep it for a while)
Specify 'from' time when subscribing to the processes output
