package route

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	maxNameLen   = 15
	maxMethodLen = len("DELETE")
)

type Route struct {
	Method     string
	Name       string
	Path       string
	HandleFunc http.HandlerFunc
}

type RoutesGroup struct {
	Name  string
	Items []Route
}

func (r *Route) String() string {
	name := r.Name + strings.Repeat(" ", maxNameLen-len(r.Name))
	method := r.Method + strings.Repeat(" ", maxMethodLen-len(r.Method))
	return fmt.Sprintf("%s\t%s\t%s", name, method, r.Path)
}
