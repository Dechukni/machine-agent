package machine

import (
	"github.com/evoevodin/machine-agent/core"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	CONNECTED = "connected"
)

type ChannelEvent struct {
	core.Event
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
)

func WsConnect(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Couldn't establish websocket connection " + err.Error())
		return
	}

	// TODO generate channel instead
	channelId := "channel123"
	go registerConnection(channelId, conn)

	channel := &ChannelEvent{}
	channel.EventType = CONNECTED
	channel.Time = time.Now()
	channel.Text = "Hello!"
	channel.ChannelId = channelId

	err = conn.WriteJSON(channel)
	if err != nil {
		log.Println("Couldn't write json to the websocket channel " + err.Error())
		return
	}
}

// TODO read bytes[] -> get operation type -> decode & populate ApiCall based on the operation
// TODO handle disconnect
// TODO appropriate locks
func registerConnection(channelId string, conn *websocket.Conn) error {
	connections.items[channelId] = conn

	for {
		//call := &ApiCall{}
		//err := conn.ReadJSON(call)
		t, body, err := conn.ReadMessage()
		//if err != nil {
		//	log.Println("ERROR: " + err.Error())
		//	continue
		//}
		log.Printf("%s, %s, %s", t, body, err)
		//execute(call)
	}
}
