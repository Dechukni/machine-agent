// TODO add subscribe api calls
package process

import (
	"encoding/json"
	"github.com/evoevodin/machine-agent/op"
)

const (
	PROCESS_START = "process.start"
	PROCESS_KILL  = "prcess.kill"
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
			func(apiCall interface{}, eventsChan chan interface{}) {
				startCall := apiCall.(StartProcessCall)
				Start(&Command{
					startCall.Name,
					startCall.CommandLine,
					startCall.Type,
				}, &Subscriber{
					DEFAULT_MASK,
					eventsChan,
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
			func(apiCall interface{}, eventsChan chan interface{}) {
				killCall := apiCall.(KillProcessCall)
				p, ok := Get(killCall.Pid)
				// TODO handle not ok
				if ok {
					p.Kill()
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
