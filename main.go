package main

import (
	"encoding/json"
	"fmt"
	"github.com/evoevodin/machine-agent/core/api"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"github.com/evoevodin/machine-agent/process"
)

var (
	ApplicationHttpRoutes = []api.HttpRoutesGroup{
		process.HttpRoutes,
	}

	ApplicationOperationRoutes = []api.OperationRoutesGroup{
	//process.OperationRoutes,
	}
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	fmt.Println("Registered Routes:\n")
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

	// TODO this is test process for testing purposes remove it from here
	//_, err := process.Start(&process.NewProcess{"ping test", "ping google.com"}, &FmtEventSubscriber{})
	//if err != nil {
	//	log.Println("Error: ", err)
	//}

	// TODO rework the mechanism of ws connections
	//router.HandleFunc("/connect", WsConnect)

	log.Fatal(http.ListenAndServe(":8080", router))
}

// TODO remove
type FmtEventSubscriber struct{}

func (sub *FmtEventSubscriber) OnEvent(event interface{}) {
	json.NewEncoder(os.Stdout).Encode(event)
}
