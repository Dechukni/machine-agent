package api

// Describes named operation which is called
// on the websocket client's side and usually
// performed on the machine-agent side, if appropriate OperationRoute exists.
// All the structures representing ApiCall objects must include this structure as well.
type ApiCall struct {

	// The operation which is represented by this ApiCall
	// usually dot separated resource and action
	Operation string `json:"operation"`
}

// Describes route for api calls
type OperationRoute struct {

	// The operation name like defined by ApiCall.Operation
	Operation string

	// The name of the operation route, used in logs
	// this name is unique for all the application operation routes
	// example: 'KillProcess'
	Name string

	// The decoder used for decoding the given object
	// into the special ApiCall, described by this operation route
	// the decoded value will be used by the ApiCallHandlerFunc from this
	// operation route, so it is up to the actual route - to define type safe
	// couple of ApiCallDecoderFunc & ApiCallHandlerFunc.
	// The source is a message read from the websocket channel
	ApiCallDecoderFunc func(source []byte) (interface{}, error)

	// Defines handler for the decoded ApiCall.
	// If handler function can't perform the operation then appropriate error
	// event should be written into the eventsChannel
	// The apiCall is value returned from the ApiCallDecoderFunc while
	// eventsChannel is a channel where all the events produced by
	// the operation should be written to
	ApiCallHandlerFunc func(apiCall interface{}, eventsChannel chan interface{})
}

// Named group of operation routes, those groups
// should be defined by separate apis, and than combined together
type OperationRoutesGroup struct {

	// The name of this group e.g.: 'ProcessOperationRoutes'
	Name string

	// The operation routes of this group
	Items []OperationRoute
}
