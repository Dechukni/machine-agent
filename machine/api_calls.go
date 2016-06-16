// TODO add subscribe api calls
package machine

const (
	PROCESS_START = "process.start"
	PROCESS_KILL  = "process.kill"
)

type ApiCall struct {
	Operation string `json:"operation"`
}

type StartProcessCall struct {
	ApiCall
	Name        string `json:"name"`
	CommandLine string `json:"commandLine"`
}

type KillProcessCall struct {
	ApiCall
	Pid       string `json:"pid"`
	NativePid string `json:"nativePid"`
}
