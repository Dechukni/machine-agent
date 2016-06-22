// TODO add subscribe api calls
package process


import (
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
