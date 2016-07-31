package process_test

import (
	"encoding/json"
	"github.com/evoevodin/machine-agent/process"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"
)

var alphabet = []byte("abcdefgh123456789")

func TestFileLoggerCreatesFileWhenFileDoesNotExist(t *testing.T) {
	filename := os.TempDir() + string(os.PathSeparator) + randomName(10)
	if _, err := os.Stat(filename); err == nil {
		t.Fatalf("File '%s' already exists", filename)
	}

	if _, err := process.NewLogger(filename); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("Expected file '%s' was created, but it wasn't", filename)
	}
	os.Remove(filename)
}

func TestFileLoggerTruncatesFileIfFileExistsOnCreate(t *testing.T) {
	filename := os.TempDir() + string(os.PathSeparator) + randomName(10)
	if _, err := os.Create(filename); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(filename, []byte("file-content"), 0666); err != nil {
		t.Fatal(err)
	}

	if _, err := process.NewLogger(filename); err != nil {
		t.Fatal(err)
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	if len(content) != 0 {
		t.Errorf("Expected file '%s' content is empty", filename)
	}
	os.Remove(filename)
}

func TestLogsAreFlushedOnClose(t *testing.T) {
	filename := os.TempDir() + string(os.PathSeparator) + randomName(10)

	fl, err := process.NewLogger(filename)
	if err != nil {
		t.Fatal(err)
	}

	// Write something to the log
	now := time.Now()
	fl.OnStdout("stdout", now)
	fl.OnStderr("stderr", now)
	fl.Close()

	// Read file content
	f, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}

	// Read log messages
	stdout := process.LogMessage{}
	stderr := process.LogMessage{}
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&stdout); err != nil {
		t.Fatal(err)
	}
	if err := decoder.Decode(&stderr); err != nil {
		t.Fatal(err)
	}

	// Check logs are okay
	expectedStdout := process.LogMessage{
		Kind: process.STDOUT_KIND,
		Time: now,
		Text: "stdout",
	}
	if stdout != expectedStdout {
		t.Fatalf("Expected %v but found %v", expectedStdout, stdout)
	}
	expectedStderr := process.LogMessage{
		Kind: process.STDERR_KIND,
		Time: now,
		Text: "stderr",
	}
	if stdout != expectedStdout {
		t.Fatalf("Expected %v but found %v", expectedStderr, stderr)
	}

	os.Remove(filename)
}

func TestReadLogs(t *testing.T) {
	filename := os.TempDir() + string(os.PathSeparator) + randomName(10)

	fl, err := process.NewLogger(filename)
	if err != nil {
		t.Fatal(err)
	}

	// Write something to the log
	now := time.Now()
	fl.OnStdout("line1", now.Add(time.Second))
	fl.OnStdout("line2", now.Add(time.Second*2))
	fl.OnStdout("line3", now.Add(time.Second*3))
	fl.OnStdout("line4", now.Add(time.Second*4))
	fl.OnStdout("line5", now.Add(time.Second*5))
	fl.Close()

	// Read logs [2, 4]
	logs, err := fl.ReadLogs(now.Add(time.Second*2), now.Add(time.Second*4))
	if err != nil {
		t.Fatal(err)
	}

	// Check everything is okay
	expected := []process.LogMessage{
		{Kind: process.STDOUT_KIND, Time: now.Add(time.Second * 2), Text: "line2"},
		{Kind: process.STDOUT_KIND, Time: now.Add(time.Second * 3), Text: "line3"},
		{Kind: process.STDOUT_KIND, Time: now.Add(time.Second * 4), Text: "line4"},
	}
	for i := 0; i < len(logs); i++ {
		if *logs[i] != expected[i] {
			t.Fatalf("Expected: '%v' Found '%v'", expected[i], *logs[i])
		}
	}
}

func randomName(length int) string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(bytes)
}
