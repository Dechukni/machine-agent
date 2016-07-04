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

var (
	currentPid uint64 = 0
	processes         = &MachineProcesses{items: make(map[uint64]*MachineProcess)}
)

type Subscriber struct {
	Mask    uint64
	Channel chan interface{}
}

type subscribers struct {
	sync.RWMutex
	items []*Subscriber
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

	command    *exec.Cmd
	pumper     *LogsPumper
	fileLogger *FileLogger
	subs       *subscribers
}

type MachineProcesses struct {
	sync.RWMutex
	items map[uint64]*MachineProcess
}

func Start(newProcess *NewProcess, firstSubscriber *Subscriber) (*MachineProcess, error) {
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
		subs:        &subscribers{},
	}
	if firstSubscriber != nil {
		process.AddSubscriber(firstSubscriber)
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
	}, PROCESS_STATUS_BIT)

	// start pumping 'pumper.Pump' is blocking
	go func() {
		defer setDead(pid)
		process.pumper.Pump()
	}()

	// publish the process
	return process, nil
}

func Get(pid uint64) (*MachineProcess, bool) {
	processes.RLock()
	defer processes.RUnlock()
	item, ok := processes.items[pid]
	return item, ok
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
	return pArr
}

func (mp *MachineProcess) Kill() error {
	// workaround for killing child processes see https://github.com/golang/go/issues/8854
	return syscall.Kill(-mp.NativePid, syscall.SIGKILL)
}

func (mp *MachineProcess) ReadLogs() ([]string, error) {
	// Getting process logs
	logs, err := mp.fileLogger.ReadLogs()
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

func (mp *MachineProcess) RemoveSubscriber(subChannel chan interface{}) {
	mp.subs.Lock()
	defer mp.subs.Unlock()
	for idx, sub := range mp.subs.items {
		if sub.Channel == subChannel {
			mp.subs.items = append(mp.subs.items[0:idx], mp.subs.items[idx+1:]...)
			break
		}
	}
}

func (mp *MachineProcess) AddSubscriber(subscriber *Subscriber) error {
	mp.subs.Lock()
	defer mp.subs.Unlock()
	for _, sub := range mp.subs.items {
		if sub.Channel == subscriber.Channel {
			return errors.New("Already subscribed")
		}
	}
	mp.subs.items = append(mp.subs.items, subscriber)
	return nil
}

func (mp *MachineProcess) UpdateSubscriber(subChannel chan interface{}, newMask uint64) {
	mp.subs.Lock()
	defer mp.subs.Unlock()
	for _, sub := range mp.subs.items {
		if sub.Channel == subChannel {
			sub.Mask = newMask
			break
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
	}, STDOUT_BIT)
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
	}, STDERR_BIT)
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
	}, PROCESS_STATUS_BIT)
}

func setDead(pid uint64) {
	processes.Lock()
	defer processes.Unlock()
	process, ok := processes.items[pid]
	if ok {
		process.Alive = false
	}
}

func (mp *MachineProcess) publish(event interface{}, typeBit uint64) {
	mp.subs.RLock()
	subs := mp.subs.items
	for _, subscriber := range subs {
		// Check whether subscriber needs such kind of event and then try to notify it
		if subscriber.Mask&typeBit == typeBit && !tryWrite(subscriber.Channel, event) {
			// Impossible to write to the channel, remove the channel from the subscribers list.
			// It may happen when writing to the closed channel
			defer mp.RemoveSubscriber(subscriber.Channel)
		}
	}
	mp.subs.RUnlock()
}

// Writes to a channel and returns true if write is successful,
// otherwise if write the channel failed e.g. channel is closed then returns false
func tryWrite(eventsChan chan interface{}, event interface{}) (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()
	eventsChan <- event
	return true
}
