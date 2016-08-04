package process

import (
	"strconv"
	"errors"
	"time"
)

// Checks whether pid is valid and converts it to the uint64
func parsePid(strPid string) (uint64, error) {
	intPid, err := strconv.Atoi(strPid)
	if err != nil {
		return 0, errors.New("Pid value must be unsigned integer")
	}
	if intPid <= 0 {
		return 0, errors.New("Pid value must be unsigned integer")
	}
	return uint64(intPid), nil
}

// Checks whether command is valid
func checkCommand(command *Command) error {
	if command.Name == "" {
		return errors.New("Command name required")
	}
	if command.CommandLine == "" {
		return errors.New("Command line required")
	}
	return nil
}

// If time string is empty, then default time is returned
// If time string is invalid, then appropriate error is returned
// If time string is valid then parsed time is returned
func parseTime(timeStr string, defTime time.Time) (time.Time, error) {
	if timeStr == "" {
		return defTime, nil
	}
	return time.Parse(DateTimeFormat, timeStr)
}
