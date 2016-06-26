package disp

import (
	"encoding/json"
	"fmt"
	"github.com/evoevodin/machine-agent/core"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	CONNECTED = "connected"
)

var (
	// Connections managed by the dispatcher
	connections = &Connections{items: make(map[string]*websocket.Conn)}

	// Registered operation routes
	opRoutes = &OpRoutes{items: make(map[string]OpRoute)}

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// TODO add authentication
			return true
		},
	}

	prevChanId uint64 = 0
)

type ChannelEvent struct {
	core.Event
	ChannelId string `json:"channelId"`
	Text      string `json:"text"`
}

// Defines lockable map for managing websocket connections
type Connections struct {
	sync.RWMutex
	items map[string]*websocket.Conn
}

// Defines lockable map for managing operation routes
type OpRoutes struct {
	sync.RWMutex
	items map[string]OpRoute
}

// Gets route by the operation name
func (opRoutes *OpRoutes) Get(operation string) (OpRoute, bool) {
	opRoutes.RLock()
	defer opRoutes.RUnlock()
	item, ok := opRoutes.items[operation]
	return item, ok
}

// Adds a new route, if the route already registered then returns false
// and doesn't override existing route, if no such route found
// then the given route will be added and true returned
func (or *OpRoutes) Add(r OpRoute) bool {
	opRoutes.Lock()
	defer opRoutes.Unlock()
	_, ok := opRoutes.items[r.Operation]
	if ok {
		return false
	}
	opRoutes.items[r.Operation] = r
	return true
}

func RegisterRoute(route OpRoute) {
	if !opRoutes.Add(route) {
		log.Fatalf("Couldn't register a new route, route for the operation '%s' already exists", route.Operation)
	}
}

func RegisterConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Couldn't establish websocket connection " + err.Error())
		return
	}

	// Generating unique channel identifier and save the connection
	// for future interactions with the API
	chanId := "channel-" + strconv.Itoa(int(atomic.AddUint64(&prevChanId, 1)))
	connections.Lock()
	connections.items[chanId] = conn
	connections.Unlock()

	// Listen for the events from the machine-agent side
	// and API calls from the channel client side
	eventsChan := make(chan interface{})
	go listenForEvents(conn, eventsChan)
	go listenForCalls(conn, eventsChan)

	// Say hello to the client
	eventsChan <- &ChannelEvent{
		core.Event{
			CONNECTED,
			time.Now(),
		},
		chanId,
		"Hello!",
	}
}

// TODO handle disconnect
func listenForCalls(conn *websocket.Conn, eventsChannel chan interface{}) {
	for {
		// Reading the message from the client
		_, body, err := conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, 1005) {
				log.Println("Error reading message, " + err.Error())
			}
			close(eventsChannel)
			break;
		}

		call := &ApiCall{}
		json.Unmarshal(body, call)
		dispatchCall(call.Operation, body, eventsChannel)
	}
}

func listenForEvents(conn *websocket.Conn, eventsChannel chan interface{}) {
	for {
		event, ok := <-eventsChannel
		if !ok {
			// channel is closed, should happen only if websocket connection is closed
			break
		}
		err := conn.WriteJSON(event)
		if err != nil {
			log.Printf("Couldn't write event to the channel. Event: %T, %v", event, event)
		}
	}
}

func dispatchCall(operation string, body []byte, eventsChannel chan interface{}) {
	// Get the requested route
	opRoute, ok := opRoutes.Get(operation)
	if !ok {
		// TODO mb respond with an error event?
		fmt.Printf("No route found for the operation '%s'", operation)
		return
	}

	// Dispatch call
	apiCall, err := opRoute.DecoderFunc(body)
	if err != nil {
		// TODO mb respond with an error event?
		fmt.Printf("Error decoding ApiCall for the operation '%s'. Error: '%s'\n", operation, err.Error())
	}
	opRoute.HandlerFunc(apiCall, eventsChannel)
}

