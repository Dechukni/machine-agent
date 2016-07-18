package main

import (
	"fmt"
	"github.com/evoevodin/machine-agent/op"
	"github.com/evoevodin/machine-agent/process"
	"log"
	"net/http"
	"github.com/evoevodin/machine-agent/core/rest"
)

var (
	AppHttpRoutes = []rest.HttpRoutesGroup{
		process.HttpRoutes,
		op.HttpRoutes,
	}

	AppOpRoutes = []op.RoutesGroup{
		process.OpRoutes,
	}
)

func main() {
	fmt.Println("⇩ Registered HttpRoutes:\n")
	for _, routesGroup := range AppHttpRoutes {
		fmt.Printf("%s:\n", routesGroup.Name)
		for _, route := range routesGroup.Items {
			fmt.Printf("✓ %s\n", &route)
			rest.RegisterRoute(route)
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

	log.Fatal(http.ListenAndServe(":8080", rest.Router))
}
