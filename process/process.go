// Exposes API for machines process management
package process

import (
	"errors"
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

	DATE_TIME_FORMAT = time.RFC3339Nano
)

var (
	prevPid   uint64 = 0
	processes        = &MachineProcesses{items: make(map[uint64]*MachineProcess)}
)

type Command struct {
	Name        string `json:"name"`
	CommandLine string `json:"commandLine"`
	Type        string `json:"type"`
}

// Defines machine process model
type MachineProcess struct {
	// The virtual id of the process, it is guaranteed  that pid
	// is always unique, while NativePid may occur twice or more(when including dead processes)
	Pid uint64 `json:"pid"`

	// The name of the process, it is equal to the Command.Name which this process created from.
	// It doesn't have to be unique, at least machine agent doesn't need such constraint,
	// as pid is used for identifying process
	Name string `json:"name"`

	// The command line executed by this process.
	// It is equal to the Command.CommandLine which this process created from
	CommandLine string `json:"commandLine"`

	// The type of the command line, this field is rather useful meta
	// information  than something used for functioning. It is equal
	// to the Command.Type which this process created from
	Type string `json:"type"`

	// Whether this process is alive or dead
	Alive bool `json:"alive"`

	// The native(OS) pid, it is unique per alive processes,
	// but those which are not alive, may have the same NativePid
	NativePid int `json:"nativePid"`

	// Command executed by this process.
	// If process is not alive then the command value is set to nil
	command *exec.Cmd

	// Stdout/stderr pumper.
	// If process is not alive then the pumper value is set to nil
	pumper *LogsPumper

	// Process subscribers, all the outgoing events are go through those subscribers.
	// If process is not alive then the subscribers value is set to nil
	subs *subscribers

	// Process file logger
	fileLogger *FileLogger
}

// Lockable map for storing processes
type MachineProcesses struct {
	sync.RWMutex
	items map[uint64]*MachineProcess
}

type Subscriber struct {
	Mask    uint64
	Channel chan interface{}
}

type subscribers struct {
	sync.RWMutex
	items []*Subscriber
}

func Start(newCommand *Command, firstSubscriber *Subscriber) (*MachineProcess, error) {
	// wrap command to be able to kill child processes see https://github.com/golang/go/issues/8854
	cmd := exec.Command("setsid", "sh", "-c", newCommand.CommandLine)

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
	pid := atomic.AddUint64(&prevPid, 1)

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
		Name:        newCommand.Name,
		CommandLine: newCommand.CommandLine,
		Type:        newCommand.Type,
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

func (mp *MachineProcess) ReadLogs(from time.Time, till time.Time) ([]*LogMessage, error) {
	return mp.fileLogger.ReadLogs(from, till)
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

// Adds a new process subscriber by reading all the logs between
// given 'after' and now and publishing them to the channel
func (mp *MachineProcess) AddBackwardSubscriber(subscriber *Subscriber, after time.Time) error {
	mp.subs.Lock()
	defer mp.subs.Unlock()
	for _, sub := range mp.subs.items {
		if sub.Channel == subscriber.Channel {
			return errors.New("Already subscribed")
		}
	}

	// Read logs between after and now
	logs, err := mp.ReadLogs(after, time.Now())
	if err != nil {
		return err
	}

	// Subscribe
	mp.subs.items = append(mp.subs.items, subscriber)

	// Publish all the logs between (after, now]
	for i := 1; i < len(logs); i++ {
		message := logs[i]
		subscriber.Channel <- newOutputEvent(mp.Pid, message.Kind, message.Text, message.Time)
	}

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
	process.publish(newOutputEvent(process.Pid, STDOUT, line, time), STDOUT_BIT)
}

func (process *MachineProcess) OnStderr(line string, time time.Time) {
	process.publish(newOutputEvent(process.Pid, STDERR, line, time), STDERR_BIT)
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
// otherwise if write to the channel failed e.g. channel is closed then returns false
func tryWrite(eventsChan chan interface{}, event interface{}) (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()
	eventsChan <- event
	return true
}

func newOutputEvent(pid uint64, kind string, line string, time time.Time) *ProcessOutputEvent {
	return &ProcessOutputEvent{
		ProcessEvent{
			core.Event{
				kind,
				time,
			},
			pid,
		},
		line,
	}
}
