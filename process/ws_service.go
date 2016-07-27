// TODO add subscribe api calls
package process

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/evoevodin/machine-agent/op"
	"time"
)

const (
	PROCESS_START             = "process.start"
	PROCESS_KILL              = "process.kill"
	PROCESS_SUBSCRIBE         = "process.subscribe"
	PROCESS_UNSUBSCRIBE       = "process.unsubscribe"
	PROCESS_UPDATE_SUBSCRIBER = "process.updateSubscriber"
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
		{
			PROCESS_UNSUBSCRIBE,
			func(body []byte) (interface{}, error) {
				call := UnsubscribeFromProcessCall{}
				err := json.Unmarshal(body, &call)
				return call, err
			},
			UnsubscribeFromProcessCallHF,
		},
		{
			PROCESS_UPDATE_SUBSCRIBER,
			func (body []byte) (interface{}, error) {
				call := UpdateProcessSubscriberCall{}
				err := json.Unmarshal(body, &call)
				return call, err
			},
			UpdateProcessSubscriberCallHF,
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

type UnsubscribeFromProcessCall struct {
	op.Call
	Pid uint64 `json:"pid"`
}

type UpdateProcessSubscriberCall struct {
	op.Call
	Pid        uint64 `json:"pid"`
	EventTypes string `json:"eventTypes"`
}

func StartProcessCallHF(call interface{}, channel op.Channel) error {
	startCall := call.(StartProcessCall)

	// Creating command
	command := &Command{
		Name:        startCall.Name,
		CommandLine: startCall.CommandLine,
		Type:        startCall.Type,
	}
	if err := checkCommand(command); err != nil {
		return err
	}

	// Detecting subscription mask
	subscriber := &Subscriber{
		Mask:    parseTypes(startCall.EventTypes),
		Channel: channel.EventsChannel,
	}

	_, err := Start(command, subscriber)
	return err
}

func KillProcessCallHF(call interface{}, channel op.Channel) error {
	killCall := call.(KillProcessCall)
	p, ok := Get(killCall.Pid)
	if !ok {
		return errors.New(fmt.Sprintf("No process with id '%s'", killCall.Pid))
	}
	return p.Kill()
}

func SubscribeToProcessCallHF(call interface{}, channel op.Channel) error {
	subscribeCall := call.(SubscribeToProcessCall)

	p, ok := Get(subscribeCall.Pid)

	if !ok {
		return errors.New(fmt.Sprintf("Process with id '%s' doesn't exist", subscribeCall.Pid))
	}

	subscriber := &Subscriber{
		Mask:    parseTypes(subscribeCall.EventTypes),
		Channel: channel.EventsChannel,
	}

	// Check whether subscriber should see previous logs or not
	if subscribeCall.After == "" {
		return p.AddSubscriber(subscriber)
	}

	after, err := time.Parse(DATE_TIME_FORMAT, subscribeCall.After)
	if err != nil {
		return errors.New("Bad format of 'after', " + err.Error())
	}
	return p.RestoreSubscriber(subscriber, after)

}

func UnsubscribeFromProcessCallHF(call interface{}, channel op.Channel) error {
	unsubscribeCall := call.(UnsubscribeFromProcessCall)
	p, ok := Get(unsubscribeCall.Pid)
	if !ok {
		return errors.New(fmt.Sprintf("Process with id '%s' doesn't exist", unsubscribeCall.Pid))
	}
	p.RemoveSubscriber(channel.EventsChannel)
	return nil
}

func UpdateProcessSubscriberCallHF(call interface{}, channel op.Channel) error {
	updateCall := call.(UpdateProcessSubscriberCall)
	p, ok := Get(updateCall.Pid)
	if !ok {
		return errors.New(fmt.Sprintf("No process with id '%d'", updateCall.Pid))
	}
	if updateCall.EventTypes == "" {
		return errors.New("'eventTypes' required for subscriber update")
	}
	p.UpdateSubscriber(channel.EventsChannel, maskFromTypes(updateCall.EventTypes))
	return nil
}
