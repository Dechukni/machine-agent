// Exposes API for machines process management
package machine

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
)

const (
// TODO configure with flag
// logsDir = "src/github.com/evoevodin/machine-agent"
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
	fileLogger  *FileLogger
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
	fmt.Println(logsDir)
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		os.MkdirAll(logsDir, 0777)
	}

	fileLogger, err := NewLogger(logsDir + "/" + strconv.Itoa(int(pid)))
	if err != nil {
		return nil, err
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
		fileLogger,
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
		// workaround for killing child processes see https://github.com/golang/go/issues/8854
		return syscall.Kill(-process.NativePid, syscall.SIGKILL)
	}
	return errors.New("No process with id " + strconv.Itoa(int(pid)))
}

func ReadProcessLogs(pid uint64) ([]string, error) {
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
