// Exposes API for machines process management
package process

import (
	"errors"
	"fmt"
	"github.com/evoevodin/machine-agent/core"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	STDOUT_BIT         = 1 << iota
	STDERR_BIT         = 1 << iota
	PROCESS_STATUS_BIT = 1 << iota
	DEFAULT_MASK       = STDERR_BIT | STDOUT_BIT | PROCESS_STATUS_BIT
)

// TODO extend with mechanism for subscription level
type ProcessSubscriber interface {

	// Returns true if ok, and false otherwise.
	// If false is returned then this subscriber will be unsubscribed
	OnEvent(event interface{}) bool
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

type MachineProcesses struct {
	sync.RWMutex
	items map[uint64]*MachineProcess
}

var (
	currentPid uint64 = 0
	processes         = &MachineProcesses{items: make(map[uint64]*MachineProcess)}
)

func Start(newProcess *NewProcess, firstSubscriber ProcessSubscriber) (*MachineProcess, error) {
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
		err = os.MkdirAll(logsDir, 0777)
		if err != nil {
			return nil, err
		}
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
	if firstSubscriber != nil {
		process.subscribers = append(process.subscribers, firstSubscriber)
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
			core.Event{
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

func GetProcesses(all bool) []*MachineProcess {
	processes.RLock()
	defer processes.RUnlock()

	pArr := make([]*MachineProcess, 0, len(processes.items))
	for _, v := range processes.items {
		if all || v.Alive {
			pArr = append(pArr, v)
		}
	}
	return pArr;
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
	subs := process.subscribers
	for idx, subscriber := range subs {
		if !subscriber.OnEvent(event) {
			// remove subscriber from the list
			defer func() {
				process.subscribers = append(subs[:idx], subs[idx+1:]...)
			}()
		}
	}
}

func (process *MachineProcess) OnStdout(line string, time time.Time) {
	process.publish(&ProcessOutputEvent{
		ProcessEvent{
			core.Event{
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
			core.Event{
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
			core.Event{
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
