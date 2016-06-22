package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	"github.com/evoevodin/machine-agent/core/api"
)

const (
	CONNECTED = "connected"
)

type ChannelEvent struct {
	api.Event
	ChannelId string `json:"channelId"`
	Text      string `json:"text"`
}

type ConnectionsMap struct {
	sync.RWMutex
	items map[string]*websocket.Conn
}

var (
	connections = &ConnectionsMap{items: make(map[string]*websocket.Conn)}

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	currentChannelId uint64 = 0
)

func WsConnect(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Couldn't establish websocket connection " + err.Error())
		return
	}

	channelId := "channel-" + strconv.Itoa(int(atomic.AddUint64(&currentChannelId, 1)))

	// TODO appropriate locks
	connections.items[channelId] = conn

	//eventsChannel := make(chan interface{});
	//go listenForEvents(conn, eventsChannel)
	//go listenForCalls(conn, eventsChannel)

	err = conn.WriteJSON(&ChannelEvent{
		api.Event{
			CONNECTED,
			time.Now(),
		},
		channelId,
		"Hello!",
	})
	if err != nil {
		log.Printf("Couldn't write json to the websocket channel '%s'. %s", channelId, err.Error())
		return
	}
}

// TODO handle disconnect
//func listenForCalls(conn *websocket.Conn, channel chan interface{}) error {
//	for {
//
//		//err := conn.ReadJSON(call)
//		t, body, err := conn.ReadMessage()
//		if err != nil {
//			log.Println("Error reading message, " + err.Error())
//			continue
//		}
//
//		call := &process.ApiCall{}
//		json.Unmarshal(body, call)
//
//		process.HandleApiCall(call.Operation, body, channel)
//
//		//if err != nil {
//		//	log.Println("ERROR: " + err.Error())
//		//	continue
//		//}
//		log.Printf("%s, %s, %s", t, body, err)
//		//execute(call)
//	}
//}
//
//func listenForEvents(conn *websocket.Conn, channel chan interface{}) {
//	for {
//		event, ok := <-channel
//		if !ok {
//			// channel is closed, means that process stopped producing output
//			break
//		}
//		err := conn.WriteJSON(event)
//		if err != nil {
//			log.Printf("Couldn't write event to the channel. Event: %T, %v", event, event)
//		}
//	}
//}
