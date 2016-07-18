// TODO add subscribe api calls
package process

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/evoevodin/machine-agent/core"
	"github.com/evoevodin/machine-agent/op"
)

const (
	PROCESS_START     = "process.start"
	PROCESS_KILL      = "prcess.kill"
	PROCESS_SUBSCRIBE = "process.subscribe"
)

var OpRoutes = op.RoutesGroup{
	"Process Routes",
	[]op.Route{

		op.Route{
			PROCESS_START,
			func(body []byte) (interface{}, error) {
				call := StartProcessCall{}
				err := json.Unmarshal(body, &call)
				return call, err
			},
			func(apiCall interface{}, channel op.Channel) {
				startCall := apiCall.(StartProcessCall)
				Start(&Command{
					startCall.Name,
					startCall.CommandLine,
					startCall.Type,
				}, &Subscriber{
					DEFAULT_MASK,
					channel.EventsChannel,
				})
			},
		},

		op.Route{
			PROCESS_KILL,
			func(body []byte) (interface{}, error) {
				call := KillProcessCall{}
				err := json.Unmarshal(body, &call)
				return call, err
			},
			func(apiCall interface{}, channel op.Channel) {
				killCall := apiCall.(KillProcessCall)
				p, ok := Get(killCall.Pid)
				// TODO handle not ok
				if ok {
					p.Kill()
				}
			},
		},

		op.Route{
			PROCESS_SUBSCRIBE,
			func(body []byte) (interface{}, error) {
				call := SubscribeToProcess{}
				err := json.Unmarshal(body, &call)
				return call, err
			},
			func(apiCall interface{}, channel op.Channel) {
				subscribeCall := apiCall.(SubscribeToProcess)

				p, ok := Get(subscribeCall.Pid)

				if !ok {
					m := fmt.Sprintf("Process with id '%s' doesn't exist", subscribeCall.Pid)
					channel.EventsChannel <- core.NewErrorEvent(errors.New(m))
					return
				}

				// Parsing mask and adding a new subscriber
				var mask uint64 = DEFAULT_MASK
				if subscribeCall.Types != "" {
					mask = maskFromTypes(subscribeCall.Types)
				}
				err := p.AddSubscriber(&Subscriber{
					mask,
					channel.EventsChannel,
				})
				if err != nil {
					channel.EventsChannel <- core.NewErrorEvent(err)
				}
			},
		},
	},
}

type StartProcessCall struct {
	op.Call
	Name        string `json:"name"`
	CommandLine string `json:"commandLine"`
	Type        string `json:"type"`
}

type KillProcessCall struct {
	op.Call
	Pid       uint64 `json:"pid"`
	NativePid uint64 `json:"nativePid"`
}

type SubscribeToProcess struct {
	op.Call
	Pid   uint64 `json:"pid"`
	Types string `json:"types"`
}
