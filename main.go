package main

import (
	"flag"
	"fmt"
	"github.com/evoevodin/machine-agent/op"
	"github.com/evoevodin/machine-agent/process"
	"github.com/evoevodin/machine-agent/rest"
	"github.com/evoevodin/machine-agent/terminal"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

var (
	AppHttpRoutes = []rest.RoutesGroup{
		process.HttpRoutes,
		op.HttpRoutes,
		terminal.HttpRoutes,
	}

	AppOpRoutes = []op.RoutesGroup{
		process.OpRoutes,
	}

	serverAddress string
)

func init() {
	flag.StringVar(&serverAddress, "addr", ":9000", "IP:PORT or :PORT the address to start the server on")
}

func main() {
	flag.Parse()

	router := mux.NewRouter().StrictSlash(true)
	fmt.Print("⇩ Registered HttpRoutes:\n\n")
	for _, routesGroup := range AppHttpRoutes {
		fmt.Printf("%s:\n", routesGroup.Name)
		for _, route := range routesGroup.Items {
			fmt.Printf("✓ %s\n", &route)
			router.
				Methods(route.Method).
				Path(route.Path).
				Name(route.Name).
				HandlerFunc(rest.ToHttpHandlerFunc(route.HandleFunc))
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

	go terminal.Activity.StartTracking()

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	http.Handle("/", router)
	server := &http.Server{
		Handler:      router,
		Addr:         serverAddress,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
