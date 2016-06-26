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

var OpRoutes = disp.OpRoutesGroup{
	"Process Routes",
	[]disp.OpRoute{

		disp.OpRoute{
			PROCESS_START,
			func(body []byte) (interface{}, error) {
				call := StartProcessCall{}
				err := json.Unmarshal(body, &call)
				return call, err
			},
			func(apiCall interface{}, eventsChannel chan interface{}) {
				startCall := apiCall.(StartProcessCall)
				Start(&NewProcess{startCall.Name, startCall.CommandLine}, eventsChannel)
			},
		},

		disp.OpRoute{
			PROCESS_KILL,
			func(body []byte) (interface{}, error) {
				call := KillProcessCall{}
				err := json.Unmarshal(body, &call)
				return call, err
			},
			func(apiCall interface{}, eventsChannel chan interface{}) {
				kilLCall := apiCall.(KillProcessCall)
				Kill(kilLCall.Pid)
			},
		},
	},
}
