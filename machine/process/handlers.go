package process

import (
	"encoding/json"
	"fmt"
	"github.com/evoevodin/machine-agent/core"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

var Routes = core.RoutesGroup{
	"MachineRoutes",
	[]core.Route{
		core.Route{
			"POST",
			"StartProcess",
			"/process",
			StartProcessHF,
		},
		core.Route{
			"GET",
			"GetProcess",
			"/process/{pid}",
			GetProcessHF,
		},
		core.Route{
			"DELETE",
			"KillProcess",
			"/process/{pid}",
			KillProcessHF,
		},
		core.Route{
			"GET",
			"GetProcessLogs",
			"/process/{pid}/logs",
			GetProcessLogsHF,
		},
	},
}

func StartProcessHF(w http.ResponseWriter, r *http.Request) {
	// getting & validating incoming data
	newProcess := NewProcess{}
	json.NewDecoder(r.Body).Decode(&newProcess)
	if newProcess.CommandLine == "" {
		http.Error(w, "Command line required", http.StatusBadRequest)
		return
	}
	if newProcess.Name == "" {
		http.Error(w, "Command name required", http.StatusBadRequest)
		return
	}

	// starting the process
	process, err := Start(&newProcess, nil)

	// writing response
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(process)
}

func GetProcessHF(w http.ResponseWriter, r *http.Request) {
	pid, ok := pidVar(w, r)
	if ok {
		// getting process
		process, err := Get(pid)

		// writing response
		if err != nil {
			http.Error(w, fmt.Sprintf("Couldn't get machine process: %s", err.Error()), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(process)
	}
}

func KillProcessHF(w http.ResponseWriter, r *http.Request) {
	pid, ok := pidVar(w, r)
	if ok {
		// killing process
		err := Kill(pid)

		// writing response
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func GetProcessLogsHF(w http.ResponseWriter, r *http.Request) {
	pid, ok := pidVar(w, r)
	if ok {
		logs, err := ReadLogs(pid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		for _, line := range logs {
			io.WriteString(w, line)
		}
	}
}

func pidVar(w http.ResponseWriter, r *http.Request) (uint64, bool) {
	vars := mux.Vars(r)
	pid, err := strconv.Atoi(vars["pid"])
	if err != nil {
		http.Error(w, "Numeric pid required", http.StatusBadRequest)
		return 0, false
	}
	if pid < 0 {
		http.Error(w, "Positive pid required", http.StatusBadRequest)
		return 0, false
	}
	return uint64(pid), true
}
