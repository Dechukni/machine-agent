package process

import (
	"strconv"
	"errors"
)

// Checks whether pid is valid and converts it to the uint64
func checkPid(strPid string) (uint64, error) {
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