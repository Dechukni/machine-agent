// TODO add subscribe api calls
package process

import (
	"encoding/json"
	"github.com/evoevodin/machine-agent/core/api"
)

const (
	PROCESS_START = "process.start"
	PROCESS_KILL  = "process.kill"
)

type StartProcessCall struct {
	api.ApiCall
	Name        string `json:"name"`
	CommandLine string `json:"commandLine"`
}

type KillProcessCall struct {
	api.ApiCall
	Pid       uint64 `json:"pid"`
	NativePid uint64 `json:"nativePid"`
}

var OperationRoutes = api.OperationRoutesGroup{
	"Process Routes",
	[]api.OperationRoute{

		api.OperationRoute{
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

		api.OperationRoute{
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
