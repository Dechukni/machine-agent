package main

import (
	"fmt"
	"github.com/evoevodin/machine-agent/machine"
	"github.com/evoevodin/machine-agent/route"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var (
	ApplicationRoutes = []route.RoutesGroup{
		route.ExampleRoutes,
		machine.MachineRoutes,
	}
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	fmt.Println("Registered Routes:\n")
	for _, routesGroup := range ApplicationRoutes {
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

	// TODO this is test process for testing purposes remove it from here
	//_, err := machine.StartProcess(&machine.NewProcess{"ping test", "ping google.com"})
	//if err != nil {
	//	log.Println("Error: ", err)
	//}

	// TODO rework the mechanism of ws connections
	router.HandleFunc("/connect", machine.WsConnect)

	log.Fatal(http.ListenAndServe(":8080", router))
}
