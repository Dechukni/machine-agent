// Exposes API for machines process management
package machine

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"sync"
	"sync/atomic"
)

const (
	// TODO configure with flag
	logsDir = "/home/jumper/go-code/src/github.com/evoevodin/machine-agent"
	stdoutPrefix = "[STDOUT] "
	stderrPrefix = "[STDERR] "
)

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
}

type MachineProcessMap struct {
	sync.RWMutex
	items map[uint64]*MachineProcess
}

var (
	currentPid uint64 = 0
	processes = &MachineProcessMap{items: make(map[uint64]*MachineProcess)}
)

func StartProcess(newProcess *NewProcess) (*MachineProcess, error) {
	cmd := exec.Command("sh", "-c", newProcess.CommandLine)

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

	fileLogger := &FileLogger{
		filename: logsDir + "/" + strconv.Itoa(int(pid)),
	}

	pumper := NewPumper(stdout, stderr)
	pumper.AddConsumer(fileLogger)
	pumper.AddConsumer(&FmtConsumer{})
	go func() {
		defer setDead(pid)
		pumper.Pump()
	}()

	// create & save process
	process := &MachineProcess{
		pid,
		newProcess.Name,
		newProcess.CommandLine,
		true,
		cmd.Process.Pid,
		cmd,
		pumper,
	}
	processes.Lock()
	processes.items[pid] = process
	processes.Unlock()

	// publish the process
	return process, nil
}

func GetProcess(pid uint64) (*MachineProcess, error) {
	processes.RLock()
	process, ok := processes.items[pid]
	processes.RUnlock()

	if !ok {
		return nil, errors.New("No process with id " + strconv.Itoa(int(pid)))
	}

	return process, nil
}

func KillProcess(pid uint64) error {
	processes.Lock()
	defer processes.Unlock()
	process, ok := processes.items[pid]
	if ok {
		return process.command.Process.Kill()
	}
	return errors.New("No process with id " + strconv.Itoa(int(pid)))
}

func setDead(pid uint64) {
	processes.Lock()
	process, ok := processes.items[pid]
	if ok {
		process.Alive = false
	}
	processes.Unlock()

}

type FmtConsumer struct{}

func (fc *FmtConsumer) AcceptStdout(line string) {
	fmt.Print(line)
}

func (fc *FmtConsumer) AcceptStderr(line string) {
	fmt.Print(line)
}

func (fc *FmtConsumer) Close() {
	fmt.Print("Closed")
}
