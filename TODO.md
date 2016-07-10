TODO

- Versioning
- Logs lifecycle
- Logs paging(time based, limits based)
- Improve logging to file(see TODOs on top of logger file)
- Flags for configuration like(logs base folder, server host:port etc)
- Plug in terminal into the infrastructure
- When to cleanup not-alive processes
- Cleanup/Restore processes when machine restored from snapshot
- Add links to the REST API responses
- Validations for REST/WS received objects
- Project structure
- Consider using RAML for API documentation

Create Dispatcher for handling api calls
---
- Defines methods for handling webscoket connections
- Manages webscoket connections state
- Exposes HttpRoutes for connections/management
