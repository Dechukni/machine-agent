// FIXME: weak implementation of file logger
// FIXME: locks, methods publicity
package machine

import (
	"bytes"
	"fmt"
	"os"
	"sync"
)

const (
	// TODO configure with a flag
	flushThreshold = 8192
)

type FileLogger struct {
	sync.RWMutex
	filename string
	buffer   bytes.Buffer
}

func (fl *FileLogger) AcceptStdout(line string) {
	fl.acceptLine(stdoutPrefix, line)
}

func (fl *FileLogger) AcceptStderr(line string) {
	fl.acceptLine(stderrPrefix, line)
}

func (fl *FileLogger) Close() {
	fl.flush()
}

func (fl *FileLogger) acceptLine(prefix string, line string) {
	buf := &fl.buffer
	fl.Lock()
	buf.WriteString(prefix + line)
	if flushThreshold < buf.Len() {
		fl.flush()
	}
	fl.Unlock()
}

func (fl *FileLogger) flush() {
	// FIXME: remove flush print
	fmt.Println("Flushing buffer")
	f, err := os.OpenFile(fl.filename, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("ERROR: couldn't open file %s \n", err.Error())
	}
	defer f.Close()
	fl.buffer.WriteTo(f)
}
