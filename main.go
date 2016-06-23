package main

import (
	"fmt"
	"github.com/evoevodin/machine-agent/core/api"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"github.com/evoevodin/machine-agent/process"
	"github.com/evoevodin/machine-agent/ws"
)

var (
	ApplicationHttpRoutes = []api.HttpRoutesGroup{
		process.HttpRoutes,
	}

	ApplicationOperationRoutes = []api.OperationRoutesGroup{
		process.OperationRoutes,
	}
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	fmt.Println("⇩ Registered HttpRoutes:\n")
	for _, routesGroup := range ApplicationHttpRoutes {
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

	// TODO rework this code in dispatcher object way
	fmt.Println("\n⇩ Registered OperationRoutes:\n")
	for _, routesGroup  := range ApplicationOperationRoutes {
		fmt.Printf("%s:\n", routesGroup.Name)
		for _, route := range routesGroup.Items {
			fmt.Printf("✓ %s\n", route.Operation)
			api.RegisteredOperationRoutes = append(api.RegisteredOperationRoutes, route)
		}
	}

	//eventsChannel := make(chan interface{});
	//
	//go func() {
	//	for {
	//		fmt.Println(<- eventsChannel)
	//	}
	//}();
	//
	//// TODO this is test process for testing purposes remove it from here
	//_, err := process.Start(&process.NewProcess{"ping test", "ping google.com"}, eventsChannel)
	//if err != nil {
	//	log.Println("Error: ", err)
	//}

	// TODO rework the mechanism of ws connections
	router.HandleFunc("/connect", ws.WsConnect)

	log.Fatal(http.ListenAndServe(":8080", router))
}
