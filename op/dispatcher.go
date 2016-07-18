package op

import (
	"encoding/json"
	"fmt"
	"github.com/evoevodin/machine-agent/core"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

var (
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

func RegisterChannel(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Couldn't establish websocket connection " + err.Error())
		return nil
	}

	// Generating unique channel identifier and save the connection
	// for future interactions with the API
	chanId := "channel-" + strconv.Itoa(int(atomic.AddUint64(&prevChanId, 1)))
	connectedTime := time.Now()
	eventsChan := make(chan interface{})
	channel := Channel{chanId, connectedTime, eventsChan, conn}
	saveChannel(channel)

	// Listen for the events from the machine-agent side
	// and API calls from the channel client side
	go listenForEvents(conn, channel)
	go listenForCalls(conn, channel)

	// Say hello to the client
	eventsChan <- &ChannelEvent{
		core.Event{
			CONNECTED,
			connectedTime,
		},
		chanId,
		"Hello!",
	}
	return nil
}

func listenForCalls(conn *websocket.Conn, channel Channel) {
	for {
		// Reading the message from the client
		_, body, err := conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, 1005) {
				log.Println("Error reading message, " + err.Error())
			}
			close(channel.EventsChannel)
			break
		}

		call := &Call{}
		json.Unmarshal(body, call)
		dispatchCall(call.Operation, body, channel)
	}
}

func listenForEvents(conn *websocket.Conn, channel Channel) {
	for {
		event, ok := <-channel.EventsChannel
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

func dispatchCall(operation string, body []byte, channel Channel) {
	// Get the requested route
	opRoute, ok := routes.get(operation)
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
	opRoute.HandlerFunc(apiCall, channel)
}
