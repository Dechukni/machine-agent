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

	NoSuchProcessErrorCode = 20000
)

var OpRoutes = op.RoutesGroup{
	"Process Routes",
	[]op.Route{
		{
			ProcessStartOp,
			func(body []byte) (interface{}, error) {
				b := StartProcessCallBody{}
				err := json.Unmarshal(body, &b)
				return b, err
			},
			StartProcessCallHF,
		},
		{
			ProcessKillOp,
			func(body []byte) (interface{}, error) {
				b := KillProcessCallBody{}
				err := json.Unmarshal(body, &b)
				return b, err
			},
			KillProcessCallHF,
		},
		{
			ProcessSubscribeOp,
			func(body []byte) (interface{}, error) {
				b := SubscribeToProcessCallBody{}
				err := json.Unmarshal(body, &b)
				return b, err
			},
			SubscribeToProcessCallHF,
		},
		{
			ProcessUnsubscribeOp,
			func(body []byte) (interface{}, error) {
				b := UnsubscribeFromProcessCallBody{}
				err := json.Unmarshal(body, &b)
				return b, err
			},
			UnsubscribeFromProcessCallHF,
		},
		{
			ProcessUpdateSubscriberOp,
			func(body []byte) (interface{}, error) {
				b := UpdateProcessSubscriberCallBody{}
				err := json.Unmarshal(body, &b)
				return b, err
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

func StartProcessCallHF(body interface{}, t op.Transmitter) error {
	startCall := body.(StartProcessCallBody)

	// Creating command
	command := Command{
		Name:        startCall.Name,
		CommandLine: startCall.CommandLine,
		Type:        startCall.Type,
	}
	if err := checkCommand(&command); err != nil {
		return op.NewArgsError(err)
	}

	// Detecting subscription mask
	subscriber := &Subscriber{
		Id:      t.Channel().Id,
		Mask:    parseTypes(startCall.EventTypes),
		Channel: t.Channel().Events,
	}

	process := NewProcess(command).BeforeEventsHook(func(process *MachineProcess) { t.Send(process) })
	if subscriber != nil {
		if err := process.AddSubscriber(subscriber); err != nil {
			return err
		}
	}

	return process.Start()
}

func KillProcessCallHF(body interface{}, t op.Transmitter) error {
	killBody := body.(KillProcessCallBody)
	p, ok := Get(killBody.Pid)
	if !ok {
		return newNoSuchProcessError(killBody.Pid)
	}
	return p.Kill()
}

func SubscribeToProcessCallHF(body interface{}, t op.Transmitter) error {
	subscribeBody := body.(SubscribeToProcessCallBody)
	p, ok := Get(subscribeBody.Pid)
	if !ok {
		return newNoSuchProcessError(subscribeBody.Pid)
	}

	subscriber := &Subscriber{
		Id:      t.Channel().Id,
		Mask:    parseTypes(subscribeBody.EventTypes),
		Channel: t.Channel().Events,
	}

	// Check whether subscriber should see previous logs or not
	if subscribeBody.After == "" {
		return p.AddSubscriber(subscriber)
	}

	after, err := time.Parse(DateTimeFormat, subscribeBody.After)
	if err != nil {
		return op.NewArgsError(errors.New("Bad format of 'after', " + err.Error()))
	}
	return p.RestoreSubscriber(subscriber, after)
}

func UnsubscribeFromProcessCallHF(call interface{}, t op.Transmitter) error {
	ubsubscribeBody := call.(UnsubscribeFromProcessCallBody)
	p, ok := Get(ubsubscribeBody.Pid)
	if !ok {
		return errors.New(fmt.Sprintf("Process with id '%s' doesn't exist", ubsubscribeBody.Pid))
	}

	p.RemoveSubscriber(t.Channel().Id)
	return nil
}

func UpdateProcessSubscriberCallHF(body interface{}, t op.Transmitter) error {
	updateBody := body.(UpdateProcessSubscriberCallBody)
	p, ok := Get(updateBody.Pid)
	if !ok {
		return newNoSuchProcessError(updateBody.Pid)
	}
	if updateBody.EventTypes == "" {
		return op.NewArgsError(errors.New("'eventTypes' required for subscriber update"))
	}

	p.UpdateSubscriber(t.Channel().Id, maskFromTypes(updateBody.EventTypes))
	return nil
}

func newNoSuchProcessError(pid uint64) op.Error {
	return op.NewError(errors.New(fmt.Sprintf("No process with id '%d'", pid)), NoSuchProcessErrorCode)
}
