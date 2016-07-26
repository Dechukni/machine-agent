// TODO add subscribe api calls
package process

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/evoevodin/machine-agent/core"
	"github.com/evoevodin/machine-agent/op"
	"time"
)

const (
	PROCESS_START     = "process.start"
	PROCESS_KILL      = "process.kill"
	PROCESS_SUBSCRIBE = "process.subscribe"
)

var OpRoutes = op.RoutesGroup{
	"Process Routes",
	[]op.Route{
		{
			PROCESS_START,
			func(body []byte) (interface{}, error) {
				call := StartProcessCall{}
				err := json.Unmarshal(body, &call)
				return call, err
			},
			StartProcessCallHF,
		},
		{
			PROCESS_KILL,
			func(body []byte) (interface{}, error) {
				call := KillProcessCall{}
				err := json.Unmarshal(body, &call)
				return call, err
			},
			KillProcessCallHF,
		},
		{
			PROCESS_SUBSCRIBE,
			func(body []byte) (interface{}, error) {
				call := SubscribeToProcessCall{}
				err := json.Unmarshal(body, &call)
				return call, err
			},
			SubscribeToProcessCallHF,
		},
	},
}

type StartProcessCall struct {
	op.Call
	Name        string `json:"name"`
	CommandLine string `json:"commandLine"`
	Type        string `json:"type"`
	EventTypes  string `json:"eventTypes"`
}

type KillProcessCall struct {
	op.Call
	Pid       uint64 `json:"pid"`
	NativePid uint64 `json:"nativePid"`
}

type SubscribeToProcessCall struct {
	op.Call
	Pid        uint64 `json:"pid"`
	EventTypes string `json:"eventTypes"`
	After      string `json:"after"`
}

func StartProcessCallHF(call interface{}, channel op.Channel) {
	startCall := call.(StartProcessCall)

	// Creating command
	command := &Command{startCall.Name, startCall.CommandLine, startCall.Type}
	if err := checkCommand(command); err != nil {
		channel.EventsChannel <- core.NewErrorEvent(err)
		return
	}

	// Detecting subscription mask
	subscriber := &Subscriber{
		parseTypes(startCall.EventTypes),
		channel.EventsChannel,
	}

	if _, err := Start(command, subscriber); err != nil {
		channel.EventsChannel <- core.NewErrorEvent(err)
	}
}

func KillProcessCallHF(call interface{}, channel op.Channel) {
	killCall := call.(KillProcessCall)
	p, ok := Get(killCall.Pid)
	if !ok {
		channel.EventsChannel <- core.NewErrorEvent(errors.New(fmt.Sprintf("No process with id '%s'", killCall.Pid)))
		return
	}
	if err := p.Kill(); err != nil {
		channel.EventsChannel <- core.NewErrorEvent(err)
	}
}

func SubscribeToProcessCallHF(call interface{}, channel op.Channel) {
	subscribeCall := call.(SubscribeToProcessCall)

	p, ok := Get(subscribeCall.Pid)

	if !ok {
		m := fmt.Sprintf("Process with id '%s' doesn't exist", subscribeCall.Pid)
		channel.EventsChannel <- core.NewErrorEvent(errors.New(m))
		return
	}

	subscriber := &Subscriber{parseTypes(subscribeCall.EventTypes), channel.EventsChannel}

	// Check whether subscriber should see previous logs or not
	if subscribeCall.After == "" {
		if err := p.AddSubscriber(subscriber); err != nil {
			channel.EventsChannel <- core.NewErrorEvent(err)
		}
	} else {
		after, err := time.Parse(DATE_TIME_FORMAT, subscribeCall.After)
		if err != nil {
			channel.EventsChannel <- core.NewErrorEvent(errors.New("Bad format of 'after', " + err.Error()))
			return
		}
		if err := p.RestoreSubscriber(subscriber, after); err != nil {
			channel.EventsChannel <- core.NewErrorEvent(err)
		}
	}
}
