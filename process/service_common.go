package process

import "strings"

const (
	DefaultLogsLimit = 50
)

func maskFromTypes(types string) uint64 {
	var mask uint64
	for _, t := range strings.Split(types, ",") {
		switch strings.ToLower(strings.TrimSpace(t)) {
		case "stderr":
			mask |= StderrBit
		case "stdout":
			mask |= StdoutBit
		case "process_status":
			mask |= ProcessStatusBit
		}
	}
	return mask
}

func parseTypes(types string) uint64 {
	var mask uint64 = DefaultMask
	if types != "" {
		mask = maskFromTypes(types)
	}
	return mask
}
