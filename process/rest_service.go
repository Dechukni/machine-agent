package process

import (
	"encoding/json"
	"fmt"
	"github.com/evoevodin/machine-agent/core/rest"
	"github.com/evoevodin/machine-agent/op"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var HttpRoutes = rest.HttpRoutesGroup{
	"Process Routes",
	[]rest.HttpRoute{
		rest.HttpRoute{
			"POST",
			"Start Process",
			"/process",
			StartProcessHF,
		},
		rest.HttpRoute{
			"GET",
			"Get Process",
			"/process/{pid}",
			GetProcessHF,
		},
		rest.HttpRoute{
			"DELETE",
			"Kill Process",
			"/process/{pid}",
			KillProcessHF,
		},
		rest.HttpRoute{
			"GET",
			"Get Process Logs",
			"/process/{pid}/logs",
			GetProcessLogsHF,
		},
		rest.HttpRoute{
			"GET",
			"Get Processes",
			"/process",
			GetProcessesHF,
		},
		rest.HttpRoute{
			"DELETE",
			"Unsubscribe from Process Events",
			"/process/{pid}/events/{channel}",
			UnsubscribeHF,
		},
		rest.HttpRoute{
			"POST",
			"Subscribe to Process Events",
			"/process/{pid}/events/{channel}",
			SubscribeHF,
		},
		rest.HttpRoute{
			"PUT",
			"Update Process Events Subscriber",
			"/process/{pid}/events/{channel}",
			UpdateSubscriberHF,
		},
	},
}

func StartProcessHF(w http.ResponseWriter, r *http.Request) {
	command := Command{}
	rest.ReadJson(r, &command)
	if err := checkCommand(&command); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// If channel is provided then check whether it is ready to be
	// first process subscriber and use it if it is
	var subscriber *Subscriber
	channelId := r.URL.Query().Get("channel")
	if channelId != "" {
		channel, ok := op.GetChannel(channelId)
		if !ok {
			m := fmt.Sprintf("Channel with id '%s' doesn't exist. Process won't be started", channelId)
			http.Error(w, m, http.StatusNotFound)
			return
		}

		var mask uint64 = DEFAULT_MASK
		types := r.URL.Query().Get("types")
		if types != "" {
			mask = maskFromTypes(types)
		}

		subscriber = &Subscriber{mask, channel.EventsChannel}
	}

	process, err := Start(&command, subscriber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rest.WriteJson(w, process)
}

func GetProcessHF(w http.ResponseWriter, r *http.Request) {
	pid, err := parsePid(mux.Vars(r)["pid"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	process, ok := Get(pid)

	if !ok {
		http.Error(w, fmt.Sprintf("No process with id '%d'", pid), http.StatusNotFound)
		return
	}
	rest.WriteJson(w, process)
}

func KillProcessHF(w http.ResponseWriter, r *http.Request) {
	pid, err := parsePid(mux.Vars(r)["pid"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p, ok := Get(pid)
	if !ok {
		http.Error(w, fmt.Sprintf("No process with id '%d'", pid), http.StatusNotFound)
		return
	}
	if err := p.Kill(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetProcessLogsHF(w http.ResponseWriter, r *http.Request) {
	pid, err := parsePid(mux.Vars(r)["pid"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p, ok := Get(pid)
	if !ok {
		http.Error(w, fmt.Sprintf("No process with id '%d'", pid), http.StatusNotFound)
		return
	}

	// Parse 'from', if 'from' is not specified then read all the logs from the start
	// if 'from' format is different from the DATE_TIME_FORMAT then return 400
	from, err := parseTime(r.URL.Query().Get("from"), time.Time{})
	if err != nil {
		http.Error(w, "Bad format of 'from', "+err.Error(), http.StatusBadRequest)
		return
	}

	// Parse 'till', if 'till' is not specified then 'now' is used for it
	// if 'till' format is different from the DATE_TIME_FORMAT then return 400
	till, err := parseTime(r.URL.Query().Get("till"), time.Now())
	if err != nil {
		http.Error(w, "Bad format of 'till', "+err.Error(), http.StatusBadRequest)
		return
	}

	logs, err := p.ReadLogs(from, till)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with an appropriate logs format, default json
	format := r.URL.Query().Get("format")
	switch strings.ToLower(format) {
	case "text":
		for _, item := range logs {
			line := fmt.Sprintf("[%s] %s \t %s", item.Kind, item.Time.Format(DATE_TIME_FORMAT), item.Text)
			io.WriteString(w, line)
		}
	case "json":
		fallthrough
	default:
		rest.WriteJson(w, logs)
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
	vars := mux.Vars(r)
	pid, err := parsePid(vars["pid"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Getting process
	p, ok := Get(pid)
	if !ok {
		http.Error(w, fmt.Sprintf("No process with id '%d'", pid), http.StatusNotFound)
		return
	}

	channelId := vars["channel"]

	// Getting channel
	channel, ok := op.GetChannel(channelId)
	if !ok {
		http.Error(w, fmt.Sprintf("Channel with id '%s' doesn't exist", channelId), http.StatusNotFound)
		return
	}

	p.RemoveSubscriber(channel.EventsChannel)
}

func SubscribeHF(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := parsePid(vars["pid"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Getting process
	p, ok := Get(pid)
	if !ok {
		http.Error(w, fmt.Sprintf("No process with id '%d'", pid), http.StatusNotFound)
		return
	}

	channelId := vars["channel"]

	// Getting channel
	channel, ok := op.GetChannel(channelId)
	if !ok {
		http.Error(w, fmt.Sprintf("Channel with id '%s' doesn't exist", channelId), http.StatusNotFound)
		return
	}

	subscriber := &Subscriber{DEFAULT_MASK, channel.EventsChannel}

	// Check whether subscriber should see previous process logs
	afterStr := r.URL.Query().Get("after")
	if afterStr == "" {
		p.AddSubscriber(subscriber)
	} else {
		after, err := time.Parse(DATE_TIME_FORMAT, afterStr)
		if err != nil {
			http.Error(w, "Bad format of 'after', "+err.Error(), http.StatusBadRequest)
			return
		}
		p.AddBackwardSubscriber(subscriber, after)
	}
}

func UpdateSubscriberHF(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := parsePid(vars["pid"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Getting process
	p, ok := Get(pid)
	if !ok {
		http.Error(w, fmt.Sprintf("No process with id '%d'", pid), http.StatusNotFound)
		return
	}

	channelId := vars["channel"]

	// Getting channel
	channel, ok := op.GetChannel(channelId)
	if !ok {
		http.Error(w, fmt.Sprintf("Channel with id '%s' doesn't exist", channelId), http.StatusNotFound)
		return
	}

	// Parsing mask from the level e.g. events?types=stdout,stderr
	types := r.URL.Query().Get("types")
	if types == "" {
		http.Error(w, "'types' parameter required", http.StatusBadRequest)
		return
	}

	p.UpdateSubscriber(channel.EventsChannel, maskFromTypes(types))
}

func maskFromTypes(types string) uint64 {
	var mask uint64
	for _, t := range strings.Split(types, ",") {
		switch strings.ToLower(strings.TrimSpace(t)) {
		case "stderr":
			mask |= STDERR_BIT
		case "stdout":
			mask |= STDOUT_BIT
		case "process_status":
			mask |= PROCESS_STATUS_BIT
		}
	}
	return mask
}
