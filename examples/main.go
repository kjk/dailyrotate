package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kjk/dailyrotate"
)

func onLogClose(path string, didRotate bool) {
	fmt.Printf("we just closed a file '%s', didRotate: %v\n", path, didRotate)
	if !didRotate {
		return
	}
	// process just closed file e.g. upload to backblaze storage for backup
	go func() {
		// if processing takes a long time, do it in background
	}()
}

var (
	logFile *dailyrotate.File
)

func openLogFile(pathFormat string, onClose func(string, bool)) error {
	w, err := dailyrotate.NewFile(pathFormat, onLogClose)
	if err != nil {
		return err
	}
	logFile = w
	return nil
}

func closeLogFile() error {
	return logFile.Close()
}

func writeToLog(msg string) error {
	_, err := logFile.Write([]byte(msg))
	return err
}

func main() {
	logDir := "logs"

	// we have to ensure the directory we want to write to
	// already exists
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		log.Fatalf("os.MkdirAll()(")
	}
	// only for the purpose of the demo, cleanup the directory
	defer os.RemoveAll(logDir)

	pathFormat := filepath.Join(logDir, "2006-01-02.txt")
	err = openLogFile(pathFormat, onLogClose)
	if err != nil {
		log.Fatalf("openLogFile failed with '%s'\n", err)
	}
	defer closeLogFile()

	err = writeToLog("hello\n")
	if err != nil {
		log.Fatalf("writeToLog() failed with '%s'\n", err)
	}

	// this is the actual path of the log file
	// we're currently writing to
	path := logFile.Path()

	err = closeLogFile()
	if err != nil {
		log.Fatalf("closeLogFile() failed with '%s'\n", err)
	}

	st, err := os.Stat(path)
	if err != nil {
		log.Fatalf("os.Stat('%s') failed with '%s'\n", path, err)
	}
	fmt.Printf("We wrote %d bytes to log file %s\n", st.Size(), path)
}
