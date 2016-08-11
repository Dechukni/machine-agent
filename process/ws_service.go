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
	ProcessGetLogs            = "process.getLogs"

	NoSuchProcessErrorCode = 20000
)

var OpRoutes = op.RoutesGroup{
	"Process Routes",
	[]op.Route{
		{
			ProcessStartOp,
			func(body []byte) (interface{}, error) {
				b := startBody{}
				err := json.Unmarshal(body, &b)
				return b, err
			},
			startProcessCallHF,
		},
		{
			ProcessKillOp,
			func(body []byte) (interface{}, error) {
				b := killBody{}
				err := json.Unmarshal(body, &b)
				return b, err
			},
			killProcessCallHF,
		},
		{
			ProcessSubscribeOp,
			func(body []byte) (interface{}, error) {
				b := subscribeBody{}
				err := json.Unmarshal(body, &b)
				return b, err
			},
			subscribeCallHF,
		},
		{
			ProcessUnsubscribeOp,
			func(body []byte) (interface{}, error) {
				b := unsubscribeBody{}
				err := json.Unmarshal(body, &b)
				return b, err
			},
			unsubscribeCallHF,
		},
		{
			ProcessUpdateSubscriberOp,
			func(body []byte) (interface{}, error) {
				b := updateSubscriberBody{}
				err := json.Unmarshal(body, &b)
				return b, err
			},
			updateSubscriberCallHF,
		},
	},
}

type startBody struct {
	Name        string `json:"name"`
	CommandLine string `json:"commandLine"`
	Type        string `json:"type"`
	EventTypes  string `json:"eventTypes"`
}

type killBody struct {
	Pid       uint64 `json:"pid"`
	NativePid uint64 `json:"nativePid"`
}

type subscribeBody struct {
	Pid        uint64 `json:"pid"`
	EventTypes string `json:"eventTypes"`
	After      string `json:"after"`
}

type subscribeResult struct {
	Pid        uint64 `json:"pid"`
	EventTypes string `json:"eventTypes"`
	Text       string `json:"text"`
}

type unsubscribeBody struct {
	Pid uint64 `json:"pid"`
}

type updateSubscriberBody struct {
	Pid        uint64 `json:"pid"`
	EventTypes string `json:"eventTypes"`
}

type processOpResult struct {
	Pid  uint64 `json:"pid"`
	Text string `json:"text"`
}

func startProcessCallHF(body interface{}, t op.Transmitter) error {
	startBody := body.(startBody)

	// Creating command
	command := Command{
		Name:        startBody.Name,
		CommandLine: startBody.CommandLine,
		Type:        startBody.Type,
	}
	if err := checkCommand(&command); err != nil {
		return op.NewArgsError(err)
	}

	// Detecting subscription mask
	subscriber := &Subscriber{
		Id:      t.Channel().Id,
		Mask:    parseTypes(startBody.EventTypes),
		Channel: t.Channel().Events,
	}

	process := NewProcess(command).BeforeEventsHook(func(process *MachineProcess) {
		t.Send(process)
	})
	if subscriber != nil {
		if err := process.AddSubscriber(subscriber); err != nil {
			return err
		}
	}

	return process.Start()
}

func killProcessCallHF(body interface{}, t op.Transmitter) error {
	killBody := body.(killBody)
	p, ok := Get(killBody.Pid)
	if !ok {
		return newNoSuchProcessError(killBody.Pid)
	}
	if err := p.Kill(); err != nil {
		return err
	}
	t.Send(&processOpResult{
		Pid:  killBody.Pid,
		Text: "Successfully killed",
	})
	return nil
}

func subscribeCallHF(body interface{}, t op.Transmitter) error {
	subscribeBody := body.(subscribeBody)
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
	if err := p.RestoreSubscriber(subscriber, after); err != nil {
		return err
	}
	t.Send(&subscribeResult{
		Pid:        p.Pid,
		EventTypes: subscribeBody.EventTypes,
		Text:       "Successfully unsubscribed",
	})
	return nil
}

func unsubscribeCallHF(call interface{}, t op.Transmitter) error {
	ubsubscribeBody := call.(unsubscribeBody)
	p, ok := Get(ubsubscribeBody.Pid)
	if !ok {
		return errors.New(fmt.Sprintf("Process with id '%s' doesn't exist", ubsubscribeBody.Pid))
	}
	p.RemoveSubscriber(t.Channel().Id)
	t.Send(&processOpResult{
		Pid:  p.Pid,
		Text: "Successfully unsubscribed",
	})
	return nil
}

func updateSubscriberCallHF(body interface{}, t op.Transmitter) error {
	updateBody := body.(updateSubscriberBody)
	p, ok := Get(updateBody.Pid)
	if !ok {
		return newNoSuchProcessError(updateBody.Pid)
	}
	if updateBody.EventTypes == "" {
		return op.NewArgsError(errors.New("'eventTypes' required for subscriber update"))
	}

	if err := p.UpdateSubscriber(t.Channel().Id, maskFromTypes(updateBody.EventTypes)); err != nil {
		return err
	}
	t.Send(&subscribeResult{
		Pid:        p.Pid,
		EventTypes: updateBody.EventTypes,
		Text:       "Subscriber successfully updated",
	})
	return nil
}

func newNoSuchProcessError(pid uint64) op.Error {
	return op.NewError(errors.New(fmt.Sprintf("No process with id '%d'", pid)), NoSuchProcessErrorCode)
}
