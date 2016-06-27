// TODO add subscribe api calls
package process

import (
	"encoding/json"
	"github.com/evoevodin/machine-agent/op"
)

const (
	PROCESS_START = "process.start"
	PROCESS_KILL  = "process.kill"
)

var (
	OpRoutes = op.RoutesGroup{
		"Process Routes",
		[]op.Route{

			op.Route{
				PROCESS_START,
				func(body []byte) (interface{}, error) {
					call := StartProcessCall{}
					err := json.Unmarshal(body, &call)
					return call, err
				},
				func(apiCall interface{}, eventsChan chan interface{}) {
					startCall := apiCall.(StartProcessCall)
					Start(&NewProcess{
						startCall.Name,
						startCall.CommandLine,
					}, &channelProcessSubscriber{eventsChan})
				},
			},

			op.Route{
				PROCESS_KILL,
				func(body []byte) (interface{}, error) {
					call := KillProcessCall{}
					err := json.Unmarshal(body, &call)
					return call, err
				},
				func(apiCall interface{}, eventsChan chan interface{}) {
					kilLCall := apiCall.(KillProcessCall)
					Kill(kilLCall.Pid)
				},
			},
		},
	}
)

type StartProcessCall struct {
	op.Call
	Name        string `json:"name"`
	CommandLine string `json:"commandLine"`
}

type KillProcessCall struct {
	op.Call
	Pid       uint64 `json:"pid"`
	NativePid uint64 `json:"nativePid"`
}

type channelProcessSubscriber struct {
	channel chan interface{}
}

func (cps *channelProcessSubscriber) OnEvent(event interface{}) bool {
	return writeCarefully(cps.channel, event)
}

// Writes to the channel and returns true if everything is ok,
// otherwise if channel is closed then recovers and returns false
func writeCarefully(eventsChan chan interface{}, event interface{}) (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()
	eventsChan <- event
	return true
}
