package rest

import (
	"net/http"
	"strings"
	"fmt"
)

const (
	maxNameLen   = 40
	maxMethodLen = len("DELETE")
)

// Describes route for http requests
type HttpRoute struct {

	// Http method e.g. 'GET'
	Method string

	// The name of the http route, used in logs
	// this name is unique for all the application http routes
	// example: 'StartProcess'
	Name string

	// The path of the http route which this route is mapped to
	// example: '/process'
	Path string

	// The function used for handling http request
	HandleFunc http.HandlerFunc
}

// Named group of http routes, those groups
// should be defined by separate apis, and then combined together
type HttpRoutesGroup struct {

	// The name of this group e.g.: 'ProcessRoutes'
	Name string

	// The http routes of this group
	Items []HttpRoute
}

func (r *HttpRoute) String() string {
	name := r.Name + " " + strings.Repeat(".", maxNameLen-len(r.Name))
	method := r.Method + strings.Repeat(" ", maxMethodLen-len(r.Method))
	return fmt.Sprintf("%s %s %s", name, method, r.Path)
}

