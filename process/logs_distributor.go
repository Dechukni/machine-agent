package process

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

const (
	DefaultMaxDirsCount = 16
)

var (
	logsDir string
)

func init() {
	curDir, _ := os.Getwd()
	curDir += string(os.PathSeparator) + "logs"
	flag.StringVar(&logsDir, "logs-dir", curDir, "Base directory for process logs")
}

// Distributes the logs between different directories
type LogsDistributor interface {

	// The implementor must guarantee that returned file name
	// is always the same for the same pid.
	// Returns an error if it is impossible to create hierarchy of
	// logs file parent folders, otherwise returns file path
	DirForPid(pid uint64) (string, error)
}

type DefaultLogsDistributor struct {
	MaxDirsCount uint
	LogsBaseDir  *string
}

func NewLogsDistributor() LogsDistributor {
	return &DefaultLogsDistributor{
		LogsBaseDir:  &logsDir,
		MaxDirsCount: DefaultMaxDirsCount,
	}
}

func (ld *DefaultLogsDistributor) DirForPid(pid uint64) (string, error) {
	// directories from 1 to maxDirsCount inclusive
	subDirName := (pid % uint64(ld.MaxDirsCount))

	// {baseLogsDir}/{subDirName}
	pidLogsDir := fmt.Sprintf("%s%c%d", *ld.LogsBaseDir, os.PathSeparator, subDirName)

	// Create subdirectory
	if info, err := os.Stat(pidLogsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(pidLogsDir, os.ModePerm); err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	} else if !info.IsDir() {
		m := fmt.Sprintf("Couldn't create a directory '%s', the name is taken by file", pidLogsDir)
		return "", errors.New(m)
	}
	return pidLogsDir, nil
}
