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
- Consider using ApiCall(operation, request_object) instead of ApiCall(operation, request_object_fields)



Create Dispatcher for handling api calls
---

- Defines ApiCall
- Defines ApiCallDispatcher
- Defines OperationRoutes/OperationRoutesGroup vs OperationDef/OperationDefGroup
- Defines methods for registering OperationRoutes/ApiCallRoutes
- Defines methods for handling webscoket connections
- Manages webscoket connections state
- Manages `eventsChannel`
- Exposes HttpRoutes for connections/management
