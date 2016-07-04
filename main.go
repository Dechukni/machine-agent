package main

import (
	"fmt"
	"github.com/evoevodin/machine-agent/core"
	"github.com/evoevodin/machine-agent/op"
	"github.com/evoevodin/machine-agent/process"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var (
	AppHttpRoutes = []core.HttpRoutesGroup{
		process.HttpRoutes,
		op.HttpRoutes,
	}

	AppOpRoutes = []op.RoutesGroup{
		process.OpRoutes,
	}
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	fmt.Println("⇩ Registered HttpRoutes:\n")
	for _, routesGroup := range AppHttpRoutes {
		fmt.Printf("%s:\n", routesGroup.Name)
		for _, route := range routesGroup.Items {
			fmt.Printf("✓ %s\n", &route)
			router.
				Methods(route.Method).
				Path(route.Path).
				Name(route.Name).
				HandlerFunc(route.HandleFunc)
		}
		fmt.Println()
	}

	fmt.Println("\n⇩ Registered OperationRoutes:\n")
	for _, routesGroup := range AppOpRoutes {
		fmt.Printf("%s:\n", routesGroup.Name)
		for _, route := range routesGroup.Items {
			fmt.Printf("✓ %s\n", route.Operation)
			op.RegisterRoute(route)
		}
	}

	// Registering documentation routes
	router.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs/"))))

	log.Fatal(http.ListenAndServe(":8080", router))
}
