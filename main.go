package main

import (
	"fmt"
	"os"
	"time"
)

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp time.Time
	JobName   string
	Action    string // "START" or "END"
	PID       int
}

// Job represents a complete job with start and end times
type Job struct {
	Name      string
	PID       int
	StartTime time.Time
	EndTime   *time.Time
	Duration  time.Duration
}

// LogMonitor handles parsing and monitoring of log files
type LogMonitor struct {
	entries []LogEntry
	jobs    map[string]*Job // key: "jobname_pid"
}

// NewLogMonitor creates a new log monitor instance
func NewLogMonitor() *LogMonitor {
	return &LogMonitor{
		entries: make([]LogEntry, 0),
		jobs:    make(map[string]*Job),
	}
}

// TODO: log parsing, job processing, generate report

func main() {
	logFile := "logs.log"

	// Check if custom log file is provided
	if len(os.Args) > 1 {
		logFile = os.Args[1]
	}

	fmt.Printf("Starting log monitoring for file: %s\n", logFile)
	fmt.Println()

	// TODO: Create log monitor
	// Parse log file
	// Process jobs to match START/END pairs
	// Generate and display report
}
