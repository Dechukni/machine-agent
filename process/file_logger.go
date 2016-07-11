package process

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

const (
	// TODO configure with a flag
	flushThreshold = 8192
	STDOUT_KIND    = "STDOUT"
	STDERR_KIND    = "STDERR"
)

type LogMessage struct {
	Kind string
	Time time.Time
	Text string
}

type FileLogger struct {
	sync.RWMutex
	filename string
	buffer   bytes.Buffer
	encoder  *json.Encoder
}

func NewLogger(filename string) (*FileLogger, error) {
	fl := &FileLogger{filename: filename}
	fl.encoder = json.NewEncoder(&fl.buffer)

	// Trying to create logs file
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return fl, nil
}

func (fl *FileLogger) OnStdout(line string, time time.Time) {
	fl.writeLine(&LogMessage{STDOUT_KIND, time, line})
}

func (fl *FileLogger) OnStderr(line string, time time.Time) {
	fl.writeLine(&LogMessage{STDERR_KIND, time, line})
}

func (fl *FileLogger) Close() {
	fl.flush()
}

func (fl *FileLogger) ReadLogs(from time.Time, till time.Time) ([]*LogMessage, error) {
	// Flushing all the logs available before 'till'
	fl.Lock()
	fl.flush()
	fl.Unlock()

	// Trying to open the logs file for reading logs
	logsFile, err := os.Open(fl.filename)
	if err != nil {
		return nil, err
	}
	defer logsFile.Close()

	// Reading logs
	logs := []*LogMessage{}
	decoder := json.NewDecoder(bufio.NewReader(logsFile))
	for {
		message := &LogMessage{}
		err = decoder.Decode(message)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if message.Time.Before(from) {
			continue
		}
		if message.Time.After(till) {
			break
		}
		logs = append(logs, message)
	}
	return logs, nil
}

func (fl *FileLogger) writeLine(message *LogMessage) {
	fl.Lock()
	fl.encoder.Encode(message)
	if flushThreshold < fl.buffer.Len() {
		fl.flush()
	}
	fl.Unlock()
}

func (fl *FileLogger) flush() {
	f, err := os.OpenFile(fl.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("Couldn't open file '%s' for flushing the buffer. %s \n", fl.filename, err.Error())
	}
	defer f.Close()
	fl.buffer.WriteTo(f)
}
