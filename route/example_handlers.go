package route

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var ExampleRoutes = RoutesGroup{
	"ExampleRoutes",
	[]Route{
		Route{
			"GET",
			"SayHello",
			"/hello/{name}",
			HelloHandler,
		},
		Route{
			"GET",
			"Welcome",
			"/",
			WelcomeHandler,
		},
	},
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "Hello %s", vars["name"])
}

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to MachineAgent")
}
