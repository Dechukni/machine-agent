package process

import (
	"encoding/json"
	"fmt"
	"github.com/evoevodin/machine-agent/core"
	"github.com/evoevodin/machine-agent/op"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var HttpRoutes = core.HttpRoutesGroup{
	"Process Routes",
	[]core.HttpRoute{
		core.HttpRoute{
			"POST",
			"Start Process",
			"/process",
			StartProcessHF,
		},
		core.HttpRoute{
			"GET",
			"Get Process",
			"/process/{pid}",
			GetProcessHF,
		},
		core.HttpRoute{
			"DELETE",
			"Kill Process",
			"/process/{pid}",
			KillProcessHF,
		},
		core.HttpRoute{
			"GET",
			"Get Process Logs",
			"/process/{pid}/logs",
			GetProcessLogsHF,
		},
		core.HttpRoute{
			"GET",
			"Get Processes",
			"/process",
			GetProcessesHF,
		},
		core.HttpRoute{
			"DELETE",
			"Unsubscribe from Process Events",
			"/process/{pid}/events/{channel}",
			UnsubscribeHF,
		},
		core.HttpRoute{
			"POST",
			"Subscribe to Process Events",
			"/process/{pid}/events/{channel}",
			SubscribeHF,
		},
		core.HttpRoute{
			"PUT",
			"Update Process Events Subscriber",
			"/process/{pid}/events/{channel}",
			UpdateSubscriberHF,
		},
	},
}

func StartProcessHF(w http.ResponseWriter, r *http.Request) {
	// getting & validating incoming data
	command := Command{}
	json.NewDecoder(r.Body).Decode(&command)
	if command.CommandLine == "" {
		http.Error(w, "Command line required", http.StatusBadRequest)
		return
	}
	if command.Name == "" {
		http.Error(w, "Command name required", http.StatusBadRequest)
		return
	}

	// starting the process
	process, err := Start(&command, nil)

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
		process, ok := Get(pid)

		// writing response
		if !ok {
			http.Error(w, fmt.Sprintf("No process with id '%d'", pid), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(process)
	}
}

func KillProcessHF(w http.ResponseWriter, r *http.Request) {
	if pid, ok := pidVar(w, r); ok {
		p, ok := Get(pid)
		if !ok {
			http.Error(w, fmt.Sprintf("No process with id '%d'", pid), http.StatusNotFound)
			return
		}
		if err := p.Kill(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func GetProcessLogsHF(w http.ResponseWriter, r *http.Request) {
	pid, ok := pidVar(w, r)
	if ok {
		p, ok := Get(pid)
		if !ok {
			http.Error(w, fmt.Sprintf("No process with id '%d'", pid), http.StatusNotFound)
			return
		}
		logs, err := p.ReadLogs()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, line := range logs {
			io.WriteString(w, line)
		}
	}
}

func GetProcessesHF(w http.ResponseWriter, r *http.Request) {
	all, err := strconv.ParseBool(r.URL.Query().Get("all"))
	if err != nil {
		all = false
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetProcesses(all))
}

func UnsubscribeHF(w http.ResponseWriter, r *http.Request) {
	if pid, ok := pidVar(w, r); ok {
		vars := mux.Vars(r)
		channelId := vars["channel"]

		// Getting process
		p, ok := Get(pid)
		if !ok {
			http.Error(w, fmt.Sprintf("No process with id '%d'", pid), http.StatusNotFound)
			return
		}

		// Getting channel
		channel, ok := op.GetChannel(channelId)
		if !ok {
			http.Error(w, fmt.Sprintf("Channel with id '%s' doesn't exist", channelId), http.StatusNotFound)
			return
		}

		p.RemoveSubscriber(channel.EventsChannel)
	}
}

func SubscribeHF(w http.ResponseWriter, r *http.Request) {
	if pid, ok := pidVar(w, r); ok {
		vars := mux.Vars(r)
		channelId := vars["channel"]

		// Getting process
		p, ok := Get(pid)
		if !ok {
			http.Error(w, fmt.Sprintf("No process with id '%d'", pid), http.StatusNotFound)
			return
		}

		// Getting channel
		channel, ok := op.GetChannel(channelId)
		if !ok {
			http.Error(w, fmt.Sprintf("Channel with id '%s' doesn't exist", channelId), http.StatusNotFound)
			return
		}

		p.AddSubscriber(&Subscriber{DEFAULT_MASK, channel.EventsChannel})
	}
}

func UpdateSubscriberHF(w http.ResponseWriter, r *http.Request) {
	if pid, ok := pidVar(w, r); ok {
		vars := mux.Vars(r)
		channelId := vars["channel"]

		// Parsing mask from the level e.g. events?types=stdout,stderr
		types := r.URL.Query().Get("types")
		if types == "" {
			http.Error(w, "'level' parameter required", http.StatusBadRequest)
			return
		}
		var mask uint64
		for _, t := range strings.Split(types, ",") {
			switch strings.ToLower(t) {
			case "stderr":
				mask |= STDERR_BIT
			case "stdout":
				mask |= STDOUT_BIT
			case "process_status":
				mask |= PROCESS_STATUS_BIT
			}
		}

		// Getting process
		p, ok := Get(pid)
		if !ok {
			http.Error(w, fmt.Sprintf("No process with id '%d'", pid), http.StatusNotFound)
			return
		}

		// Getting channel
		channel, ok := op.GetChannel(channelId)
		if !ok {
			http.Error(w, fmt.Sprintf("Channel with id '%s' doesn't exist", channelId), http.StatusNotFound)
			return
		}

		p.UpdateSubscriber(channel.EventsChannel, mask)
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
