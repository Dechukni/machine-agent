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
	ProcessStartOp            = "process.start"
	ProcessKillOp             = "process.kill"
	ProcessSubscribeOp        = "process.subscribe"
	ProcessUnsubscribeOp      = "process.unsubscribe"
	ProcessUpdateSubscriberOp = "process.updateSubscriber"
)

var OpRoutes = op.RoutesGroup{
	"Process Routes",
	[]op.Route{
		{
			ProcessStartOp,
			func(body []byte) (interface{}, error) {
				call := StartProcessCallBody{}
				err := json.Unmarshal(body, &call)
				return call, err
			},
			StartProcessCallHF,
		},
		{
			ProcessKillOp,
			func(body []byte) (interface{}, error) {
				call := KillProcessCallBody{}
				err := json.Unmarshal(body, &call)
				return call, err
			},
			KillProcessCallHF,
		},
		{
			ProcessSubscribeOp,
			func(body []byte) (interface{}, error) {
				call := SubscribeToProcessCallBody{}
				err := json.Unmarshal(body, &call)
				return call, err
			},
			SubscribeToProcessCallHF,
		},
		{
			ProcessUnsubscribeOp,
			func(body []byte) (interface{}, error) {
				call := UnsubscribeFromProcessCallBody{}
				err := json.Unmarshal(body, &call)
				return call, err
			},
			UnsubscribeFromProcessCallHF,
		},
		{
			ProcessUpdateSubscriberOp,
			func(body []byte) (interface{}, error) {
				call := UpdateProcessSubscriberCallBody{}
				err := json.Unmarshal(body, &call)
				return call, err
			},
			UpdateProcessSubscriberCallHF,
		},
	},
}

type StartProcessCallBody struct {
	Name        string `json:"name"`
	CommandLine string `json:"commandLine"`
	Type        string `json:"type"`
	EventTypes  string `json:"eventTypes"`
}

type KillProcessCallBody struct {
	Pid       uint64 `json:"pid"`
	NativePid uint64 `json:"nativePid"`
}

type SubscribeToProcessCallBody struct {
	Pid        uint64 `json:"pid"`
	EventTypes string `json:"eventTypes"`
	After      string `json:"after"`
}

type UnsubscribeFromProcessCallBody struct {
	Pid uint64 `json:"pid"`
}

type UpdateProcessSubscriberCallBody struct {
	Pid        uint64 `json:"pid"`
	EventTypes string `json:"eventTypes"`
}

func StartProcessCallHF(body interface{}, channel op.Channel) error {
	startCall := body.(StartProcessCallBody)

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
	killCall := call.(KillProcessCallBody)
	p, ok := Get(killCall.Pid)
	if !ok {
		return errors.New(fmt.Sprintf("No process with id '%s'", killCall.Pid))
	}
	return p.Kill()
}

func SubscribeToProcessCallHF(call interface{}, channel op.Channel) error {
	subscribeCall := call.(SubscribeToProcessCallBody)

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

	after, err := time.Parse(DateTimeFormat, subscribeCall.After)
	if err != nil {
		return errors.New("Bad format of 'after', " + err.Error())
	}
	return p.RestoreSubscriber(subscriber, after)

}

func UnsubscribeFromProcessCallHF(call interface{}, channel op.Channel) error {
	unsubscribeCall := call.(UnsubscribeFromProcessCallBody)
	p, ok := Get(unsubscribeCall.Pid)
	if !ok {
		return errors.New(fmt.Sprintf("Process with id '%s' doesn't exist", unsubscribeCall.Pid))
	}
	p.RemoveSubscriber(channel.EventsChannel)
	return nil
}

func UpdateProcessSubscriberCallHF(call interface{}, channel op.Channel) error {
	updateCall := call.(UpdateProcessSubscriberCallBody)
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
