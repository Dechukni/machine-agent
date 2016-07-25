package op

import (
	"log"
	"sync"
)

var (
	// Registered operation routes
	routes = &routesMap{items: make(map[string]Route)}
)

// Describes named operation which is called
// on the websocket client's side and usually
// performed on the machine-agent side, if appropriate Route exists.
// All the structures representing Call objects must include this structure as well.
type Call struct {

	// The operation which is represented by this Call
	// usually dot separated resource and action
	Operation string `json:"operation"`
}

// Describes route for api calls
type Route struct {

	// The operation name like defined by Call.Operation
	Operation string

	// The decoder used for decoding a given object
	// into the special Call, described by this operation route.
	// The decoded value will be used by the HandlerFunc from this
	// operation route, so it is up to the actual route - to define type safe
	// couple of DecoderFunc & HandlerFunc.
	// The source is a message read from a webscoket channel.
	DecoderFunc func(source []byte) (interface{}, error)

	// Defines handler for the decoded operation Call.
	// If handler function can't perform the operation then appropriate error
	// event should be published into the eventsChannel
	// The call is a value returned from the DecoderFunc.
	HandlerFunc func(call interface{}, channel Channel)
}

// Named group of operation routes, those groups
// should be defined by separate apis, and than combined together
type RoutesGroup struct {
	// The name of this group e.g.: 'ProcessOperationRoutes'
	Name string

	// The operation routes of this group
	Items []Route
}

// Defines lockable map for managing operation routes
type routesMap struct {
	sync.RWMutex
	items map[string]Route
}

// Gets route by the operation name
func (routes *routesMap) get(operation string) (Route, bool) {
	routes.RLock()
	defer routes.RUnlock()
	item, ok := routes.items[operation]
	return item, ok
}

// Adds a new route, if the route already registered then returns false
// and doesn't override existing route, if no such route found
// then the given route will be added and true returned
func (or *routesMap) add(r Route) bool {
	routes.Lock()
	defer routes.Unlock()
	_, ok := routes.items[r.Operation]
	if ok {
		return false
	}
	routes.items[r.Operation] = r
	return true
}

// Adds a new route, panics if such route already exists
// This is designed to be used on the app bootstrap
func RegisterRoute(route Route) {
	if !routes.add(route) {
		log.Fatalf("Couldn't register a new route, route for the operation '%s' already exists", route.Operation)
	}
}
