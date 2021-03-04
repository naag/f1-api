package server

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	hclog "github.com/hashicorp/go-hclog"
)

func getLogger(scenario string, verbose bool) hclog.Logger {
	opts := &hclog.LoggerOptions{
		Output:     os.Stderr,
		Level:      hclog.Trace,
		JSONFormat: false,
	}

	if !verbose {
		logFile := redirectLoggingToFile(opts, scenario)
		fmt.Printf("Sending logs to log file %s\n", logFile)
	}

	return hclog.New(opts)
}

func directoryExists(path string) bool {
	directoryPath := filepath.Dir(path)
	if _, err := os.Stat(directoryPath); !os.IsNotExist(err) {
		return true
	}
	return false
}

func getGeneratedLogFilePath(scenario string) string {
	return filepath.Join(
		os.TempDir(),
		fmt.Sprintf("f1-%s-%s.log", scenario, time.Now().Format("2006-01-02_15-04-05")),
	)
}

func getLogFilePath(scenario string) string {
	logFilePath := os.Getenv("LOG_FILE_PATH")
	if logFilePath != "" && directoryExists(logFilePath) {
		return logFilePath
	}
	return getGeneratedLogFilePath(scenario)
}

func redirectLoggingToFile(opts *hclog.LoggerOptions, scenario string) string {
	logFilePath := getLogFilePath(scenario)
	dir := path.Dir(logFilePath)
	os.MkdirAll(dir, 0755)

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err == nil {
		opts.Output = file
	} else {
		fmt.Println("Failed to log to file, using default stderr")
	}

	return logFilePath
}
