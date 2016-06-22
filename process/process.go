// Exposes API for machines process management
package process

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"github.com/evoevodin/machine-agent/core/api"
)

const (
// TODO configure with flag
// logsDir = "src/github.com/evoevodin/machine-agent"
)

type ProcessSubscriber interface {
	OnEvent(event interface{})
}

type NewProcess struct {
	Name        string `json:"name"`
	CommandLine string `json:"commandLine"`
}

type MachineProcess struct {
	Pid         uint64 `json:"pid"`
	Name        string `json:"name"`
	CommandLine string `json:"commandLine"`
	Alive       bool   `json:"alive"`
	NativePid   int    `json:"nativePid"`

	command     *exec.Cmd
	pumper      *LogsPumper
	fileLogger  *FileLogger
	subscribers []ProcessSubscriber
}

type MachineProcessMap struct {
	sync.RWMutex
	items map[uint64]*MachineProcess
}

var (
	currentPid uint64 = 0
	processes = &MachineProcessMap{items: make(map[uint64]*MachineProcess)}
)

func Start(newProcess *NewProcess, subscriber ProcessSubscriber) (*MachineProcess, error) {
	// wrap command to be able to kill child processes see https://github.com/golang/go/issues/8854
	cmd := exec.Command("setsid", "sh", "-c", newProcess.CommandLine)

	// getting stdout pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	// getting stderr pipe
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	// starting a new process
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	// increment current pid & assign it to the value
	pid := atomic.AddUint64(&currentPid, 1)

	// FIXME: remove as it will be configurable with a flag
	logsDir := os.Getenv("GOPATH") + "/src/github.com/evoevodin/machine-agent/logs"
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		os.MkdirAll(logsDir, 0777)
	}

	fileLogger, err := NewLogger(logsDir + "/" + strconv.Itoa(int(pid)))
	if err != nil {
		return nil, err
	}

	// create & save process
	process := &MachineProcess{
		Pid:         pid,
		Name:        newProcess.Name,
		CommandLine: newProcess.CommandLine,
		Alive:       true,
		NativePid:   cmd.Process.Pid,
		command:     cmd,
		pumper:      NewPumper(stdout, stderr),
		fileLogger:  fileLogger,
	}
	if subscriber != nil {
		process.subscribers = append(process.subscribers, subscriber)
	}
	processes.Lock()
	processes.items[pid] = process
	processes.Unlock()

	// register logs consumers
	process.pumper.AddConsumer(fileLogger)
	process.pumper.AddConsumer(process)

	// before pumping is started publish process_started event
	process.publish(&ProcessStatusEvent{
		ProcessEvent{
			api.Event{
				PROCESS_STARTED,
				time.Now(),
			},
			process.Pid,
		},
		process.NativePid,
		process.Name,
		process.CommandLine,
	})

	// start pumping 'pumper.Pump' is blocking
	go func() {
		defer setDead(pid)
		process.pumper.Pump()
	}()

	// publish the process
	return process, nil
}

func Get(pid uint64) (*MachineProcess, error) {
	processes.RLock()
	process, ok := processes.items[pid]
	processes.RUnlock()

	if !ok {
		return nil, errors.New("No process with id " + strconv.Itoa(int(pid)))
	}

	return process, nil
}

func Kill(pid uint64) error {
	processes.Lock()
	defer processes.Unlock()
	process, ok := processes.items[pid]
	if ok {
		// workaround for killing child processes see https://github.com/golang/go/issues/8854
		return syscall.Kill(-process.NativePid, syscall.SIGKILL)
	}
	return errors.New("No process with id " + strconv.Itoa(int(pid)))
}

func ReadLogs(pid uint64) ([]string, error) {
	processes.RLock()
	defer processes.RUnlock()

	// Getting process
	process, ok := processes.items[pid]
	if !ok {
		return nil, errors.New("No process with id " + strconv.Itoa(int(pid)))
	}

	// Getting process logs
	logs, err := process.fileLogger.ReadLogs()
	if err != nil {
		return nil, err
	}

	// Transforming process logs
	formattedLogs := make([]string, len(logs))
	for idx, item := range logs {
		formattedLogs[idx] = fmt.Sprintf("[%s] %s \t %s", item.Kind, item.Time, item.Text)
	}
	return formattedLogs, nil
}

func setDead(pid uint64) {
	processes.Lock()
	defer processes.Unlock()
	process, ok := processes.items[pid]
	if ok {
		process.Alive = false
	}
}

func (process *MachineProcess) publish(event interface{}) {
	for _, subscriber := range process.subscribers {
		subscriber.OnEvent(event)
	}
}

func (process *MachineProcess) OnStdout(line string, time time.Time) {
	process.publish(&ProcessOutputEvent{
		ProcessEvent{
			api.Event{
				STDOUT,
				time,
			},
			process.Pid,
		},
		line,
	})
}

func (process *MachineProcess) OnStderr(line string, time time.Time) {
	process.publish(&ProcessOutputEvent{
		ProcessEvent{
			api.Event{
				STDERR,
				time,
			},
			process.Pid,
		},
		line,
	})
}

func (process *MachineProcess) Close() {
	process.publish(&ProcessStatusEvent{
		ProcessEvent{
			api.Event{
				PROCESS_DIED,
				time.Now(),
			},
			process.Pid,
		},
		process.NativePid,
		process.Name,
		process.CommandLine,
	})
}
