// FIXME: weak implementation of file logger
// FIXME: locks, methods publicity
// TODO: consider data format, improve performance of read/write in/out lock
package machine

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"
	"encoding/json"
	"log"
)

const (
	// TODO configure with a flag
	flushThreshold = 8192
	STDOUT = "STDOUT"
	STDERR = "STDERR"
)

type LogMessage struct {
	Kind string // stderr | stdout (TODO consider using iota constants)
	Time time.Time
	Text string
}

type FileLogger struct {
	sync.RWMutex
	filename string
	buffer   bytes.Buffer
}

func (fl *FileLogger) AcceptStdout(line string) {
	fl.writeLine(&LogMessage{STDOUT, time.Now(), line})
}

func (fl *FileLogger) AcceptStderr(line string) {
	fl.writeLine(&LogMessage{STDERR, time.Now(), line})
}

func (fl *FileLogger) Close() {
	fl.flush()
}

func (fl *FileLogger) ReadLogs() ([]LogMessage, error) {
	now := time.Now()

	// Flushing all the logs available before and exactly right 'now'
	fl.Lock()
	fl.flush()
	fl.Unlock()

	// Trying to open the logs file for reading the logs
	f, err := os.Open(fl.filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// TODO
	messages := []LogMessage{}
	log.Println("Logs reading is not implemented for now %s =)", now)
	return messages, nil
}

func (fl *FileLogger) writeLine(message *LogMessage) {
	buf := &fl.buffer
	fl.Lock()
	json.NewEncoder(buf).Encode(message)
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
		log.Printf("Couldn't open file '%s' for flushing the buffer. %s \n", fl.filename, err.Error())
	}
	defer f.Close()
	fl.buffer.WriteTo(f)
}
