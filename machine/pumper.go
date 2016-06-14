// TODO: consider using []byte instead of strings
// TODO: consider using go channels instead of LogsConsumers
package machine

import (
	"bufio"
	"io"
	"sync"
	"log"
)

type acceptLine func(line string)

// LogsPumper client consumes a message read by pumper
// TODO: consider channels for this purpose
type LogsConsumer interface {
	// called on each line pumped from process stdout
	AcceptStdout(line string)

	// called on each line pumped from process stderr
	AcceptStderr(line string)

	// called when pumping is finished either by normal return or by error
	Close()
}

// Pumps lines from the stdout and stderr
type LogsPumper struct {
	stdout    io.Reader
	stderr    io.Reader
	clients   []LogsConsumer
	waitGroup sync.WaitGroup
}

func NewPumper(stdout io.Reader, stderr io.Reader) *LogsPumper {
	return &LogsPumper{
		stdout: stdout,
		stderr: stderr,
	}
}

func (pumper *LogsPumper) AddConsumer(consumer LogsConsumer) {
	pumper.clients = append(pumper.clients, consumer)
}

// Start 'pumping' logs from the stdout and stderr
// The method execution is synchronous and waits for
// both stderr and stdout to complete closing all the clients after
func (pumper *LogsPumper) Pump() {
	pumper.waitGroup.Add(2)

	// reading from stdout & stderr
	go pump(pumper.stdout, pumper.AcceptStdout, &pumper.waitGroup)
	go pump(pumper.stderr, pumper.AcceptStderr, &pumper.waitGroup)

	// cleanup after pumping is complete
	pumper.waitGroup.Wait()
	pumper.Close()
}

func pump(r io.Reader, lineConsumer acceptLine, wg *sync.WaitGroup) {
	defer wg.Done()
	br := bufio.NewReader(r)
	for {
		line, err := br.ReadBytes('\n')

		if err != nil {
			// handle not normal exit
			if err != io.EOF {
				log.Println("Error pumping: " + err.Error())
			}
			return
		}

		lineConsumer(string(line))
	}
}

func (pumper *LogsPumper) AcceptStdout(line string) {
	for _, client := range pumper.clients {
		client.AcceptStdout(line)
	}
}

func (pumper *LogsPumper) AcceptStderr(line string) {
	for _, client := range pumper.clients {
		client.AcceptStderr(line)
	}
}

func (pumper *LogsPumper) Close() {
	for _, client := range pumper.clients {
		client.Close()
	}
}
