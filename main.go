package main

import (
	"fmt"
	"github.com/evoevodin/machine-agent/op"
	"github.com/evoevodin/machine-agent/process"
	"log"
	"net/http"
	"github.com/evoevodin/machine-agent/rest"
)

var (
	AppHttpRoutes = []rest.RoutesGroup{
		process.HttpRoutes,
		op.HttpRoutes,
	}

	AppOpRoutes = []op.RoutesGroup{
		process.OpRoutes,
	}
)

func main() {
	fmt.Print("⇩ Registered HttpRoutes:\n\n")
	for _, routesGroup := range AppHttpRoutes {
		fmt.Printf("%s:\n", routesGroup.Name)
		for _, route := range routesGroup.Items {
			fmt.Printf("✓ %s\n", &route)
			rest.RegisterRoute(route)
		}
		fmt.Println()
	}

	fmt.Print("\n⇩ Registered OperationRoutes:\n\n")
	for _, routesGroup := range AppOpRoutes {
		fmt.Printf("%s:\n", routesGroup.Name)
		for _, route := range routesGroup.Items {
			fmt.Printf("✓ %s\n", route.Operation)
			op.RegisterRoute(route)
		}
	}

	log.Fatal(http.ListenAndServe(":8080", rest.Router))
}
