package process

import (
	"errors"
	"fmt"
	"github.com/evoevodin/machine-agent/core/rest"
	"github.com/evoevodin/machine-agent/core/rest/restuitl"
	"github.com/evoevodin/machine-agent/op"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var HttpRoutes = rest.RoutesGroup{
	"Process Routes",
	[]rest.Route{
		{
			"POST",
			"Start Process",
			"/process",
			StartProcessHF,
		},
		{
			"GET",
			"Get Process",
			"/process/{pid}",
			GetProcessHF,
		},
		{
			"DELETE",
			"Kill Process",
			"/process/{pid}",
			KillProcessHF,
		},
		{
			"GET",
			"Get Process Logs",
			"/process/{pid}/logs",
			GetProcessLogsHF,
		},
		{
			"GET",
			"Get Processes",
			"/process",
			GetProcessesHF,
		},
		{
			"DELETE",
			"Unsubscribe from Process Events",
			"/process/{pid}/events/{channel}",
			UnsubscribeHF,
		},
		{
			"POST",
			"Subscribe to Process Events",
			"/process/{pid}/events/{channel}",
			SubscribeHF,
		},
		{
			"PUT",
			"Update Process Events Subscriber",
			"/process/{pid}/events/{channel}",
			UpdateSubscriberHF,
		},
	},
}

func StartProcessHF(w http.ResponseWriter, r *http.Request) error {
	command := Command{}
	restutil.ReadJson(r, &command)
	if err := checkCommand(&command); err != nil {
		return rest.BadRequest(err)
	}

	// If channel is provided then check whether it is ready to be
	// first process subscriber and use it if it is
	var subscriber *Subscriber
	channelId := r.URL.Query().Get("channel")
	if channelId != "" {
		channel, ok := op.GetChannel(channelId)
		if !ok {
			m := fmt.Sprintf("Channel with id '%s' doesn't exist. Process won't be started", channelId)
			return rest.NotFound(errors.New(m))
		}
		subscriber = &Subscriber{parseTypes(r.URL.Query().Get("types")), channel.EventsChannel}
	}

	process, err := Start(&command, subscriber)
	if err != nil {
		return err
	}
	return restutil.WriteJson(w, process)
}

func GetProcessHF(w http.ResponseWriter, r *http.Request) error {
	pid, err := parsePid(mux.Vars(r)["pid"])
	if err != nil {
		return rest.BadRequest(err)
	}

	process, ok := Get(pid)

	if !ok {
		return rest.NotFound(errors.New(fmt.Sprintf("No process with id '%d'", pid)))
	}
	return restutil.WriteJson(w, process)
}

func KillProcessHF(w http.ResponseWriter, r *http.Request) error {
	pid, err := parsePid(mux.Vars(r)["pid"])
	if err != nil {
		return rest.BadRequest(err)
	}
	p, ok := Get(pid)
	if !ok {
		return rest.NotFound(errors.New(fmt.Sprintf("No process with id '%d'", pid)))
	}
	if err := p.Kill(); err != nil {
		return err
	}
	return nil
}

func GetProcessLogsHF(w http.ResponseWriter, r *http.Request) error {
	pid, err := parsePid(mux.Vars(r)["pid"])
	if err != nil {
		return rest.BadRequest(err)
	}
	p, ok := Get(pid)
	if !ok {
		return rest.NotFound(errors.New(fmt.Sprintf("No process with id '%d'", pid)))
	}

	// Parse 'from', if 'from' is not specified then read all the logs from the start
	// if 'from' format is different from the DATE_TIME_FORMAT then return 400
	from, err := parseTime(r.URL.Query().Get("from"), time.Time{})
	if err != nil {
		return rest.BadRequest(errors.New("Bad format of 'from', " + err.Error()))
	}

	// Parse 'till', if 'till' is not specified then 'now' is used for it
	// if 'till' format is different from the DATE_TIME_FORMAT then return 400
	till, err := parseTime(r.URL.Query().Get("till"), time.Now())
	if err != nil {
		return rest.BadRequest(errors.New("Bad format of 'till', " + err.Error()))
	}

	logs, err := p.ReadLogs(from, till)
	if err != nil {
		return err
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
		return restutil.WriteJson(w, logs)
	}
	return restutil.WriteJson(w, logs)
}

func GetProcessesHF(w http.ResponseWriter, r *http.Request) error {
	all, err := strconv.ParseBool(r.URL.Query().Get("all"))
	if err != nil {
		all = false
	}
	return restutil.WriteJson(w, GetProcesses(all))
}

func UnsubscribeHF(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	pid, err := parsePid(vars["pid"])
	if err != nil {
		return rest.BadRequest(err)
	}

	// Getting process
	p, ok := Get(pid)
	if !ok {
		return rest.NotFound(errors.New(fmt.Sprintf("No process with id '%d'", pid)))
	}

	channelId := vars["channel"]

	// Getting channel
	channel, ok := op.GetChannel(channelId)
	if !ok {
		return rest.NotFound(errors.New(fmt.Sprintf("Channel with id '%s' doesn't exist", channelId)))
	}

	p.RemoveSubscriber(channel.EventsChannel)
	return nil
}

func SubscribeHF(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	pid, err := parsePid(vars["pid"])
	if err != nil {
		return rest.BadRequest(err)
	}

	// Getting process
	p, ok := Get(pid)
	if !ok {
		return rest.NotFound(errors.New(fmt.Sprintf("No process with id '%d'", pid)))
	}

	channelId := vars["channel"]

	// Getting channel
	channel, ok := op.GetChannel(channelId)
	if !ok {
		return errors.New(fmt.Sprintf("Channel with id '%s' doesn't exist", channelId))
	}

	subscriber := &Subscriber{parseTypes(r.URL.Query().Get("types")), channel.EventsChannel}

	// Check whether subscriber should see previous process logs
	afterStr := r.URL.Query().Get("after")
	if afterStr == "" {
		return p.AddSubscriber(subscriber)
	}
	after, err := time.Parse(DATE_TIME_FORMAT, afterStr)
	if err != nil {
		return rest.BadRequest(errors.New("Bad format of 'after', " + err.Error()))
	}
	return p.RestoreSubscriber(subscriber, after)

}

func UpdateSubscriberHF(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	pid, err := parsePid(vars["pid"])
	if err != nil {
		return rest.BadRequest(err)
	}

	// Getting process
	p, ok := Get(pid)
	if !ok {
		return rest.NotFound(errors.New(fmt.Sprintf("No process with id '%d'", pid)))
	}

	channelId := vars["channel"]

	// Getting channel
	channel, ok := op.GetChannel(channelId)
	if !ok {
		return rest.NotFound(errors.New(fmt.Sprintf("Channel with id '%s' doesn't exist", channelId)))
	}

	// Parsing mask from the level e.g. events?types=stdout,stderr
	types := r.URL.Query().Get("types")
	if types == "" {
		return rest.BadRequest(errors.New("'types' parameter required"))
	}

	p.UpdateSubscriber(channel.EventsChannel, maskFromTypes(types))
	return nil
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

func parseTypes(types string) uint64 {
	var mask uint64 = DEFAULT_MASK
	if types != "" {
		mask = maskFromTypes(types)
	}
	return mask
}
