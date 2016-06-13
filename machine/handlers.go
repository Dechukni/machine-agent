package machine

import (
	"encoding/json"
	"fmt"
	"github.com/evoevodin/machine-agent/route"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var MachineRoutes = route.RoutesGroup{
	"MachineRoutes",
	[]route.Route{
		route.Route{
			"POST",
			"StartProcess",
			"/process",
			StartProcessHF,
		},
		route.Route{
			"GET",
			"GetProcess",
			"/process/{pid}",
			GetProcessHF,
		},
		route.Route{
			"DELETE",
			"KillProcess",
			"/process/{pid}",
			KillProcessHF,
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
	process, err := StartProcess(&newProcess)

	// writing response
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(process)
}

func GetProcessHF(w http.ResponseWriter, r *http.Request) {
	// getting & validating incoming data
	vars := mux.Vars(r)
	pid, err := strconv.Atoi(vars["pid"])
	if err != nil {
		http.Error(w, "Positive numeric pid required", http.StatusBadRequest)
		return
	}

	// getting process
	process, err := GetProcess(uint64(pid))

	// writing response
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't get machine process: %s", err.Error()), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(process)
}

func KillProcessHF(w http.ResponseWriter, r *http.Request) {
	// getting & validating incoming pid
	vars := mux.Vars(r)
	pid, err := strconv.Atoi(vars["pid"])
	if err != nil {
		http.Error(w, "Positive numeric pid required", http.StatusBadRequest)
		return
	}

	// killing process
	err = KillProcess(uint64(pid))

	// writing response
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
