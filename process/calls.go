// TODO add subscribe api calls
package process

import (
	"encoding/json"
	"github.com/evoevodin/machine-agent/disp"
)

const (
	PROCESS_START = "process.start"
	PROCESS_KILL  = "process.kill"
)

var (
	OpRoutes = disp.OpRoutesGroup{
		"Process Routes",
		[]disp.OpRoute{

			disp.OpRoute{
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

			disp.OpRoute{
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
	disp.ApiCall
	Name        string `json:"name"`
	CommandLine string `json:"commandLine"`
}

type KillProcessCall struct {
	disp.ApiCall
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
