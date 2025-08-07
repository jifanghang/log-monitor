package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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

// ParseLogFile reads and parses the CSV log file
func (lm *LogMonitor) ParseLogFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Set comma as delimiter and allow variable number of fields
	reader.Comma = ','
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}

	for i, record := range records {
		if len(record) < 4 {
			log.Printf("Warning: Skipping line %d - insufficient fields: %v", i+1, record)
			continue
		}

		// Parse timestamp (assuming same day, format HH:MM:SS)
		timeStr := strings.TrimSpace(record[0])
		timestamp, err := time.Parse("15:04:05", timeStr)
		if err != nil {
			log.Printf("Warning: Failed to parse timestamp on line %d: %s", i+1, timeStr)
			continue
		}

		// Parse PID
		pidStr := strings.TrimSpace(record[3])
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			log.Printf("Warning: Failed to parse PID on line %d: %s", i+1, pidStr)
			continue
		}

		entry := LogEntry{
			Timestamp: timestamp,
			JobName:   strings.TrimSpace(record[1]),
			Action:    strings.TrimSpace(record[2]),
			PID:       pid,
		}

		lm.entries = append(lm.entries, entry)
	}

	log.Printf("Successfully parsed %d log entries from %s", len(lm.entries), filename)
	return nil
}

// ProcessJobs matches START and END entries to create complete job records
func (lm *LogMonitor) ProcessJobs() {
	for _, entry := range lm.entries {
		jobKey := fmt.Sprintf("%s_%d", entry.JobName, entry.PID)

		switch entry.Action {
		case "START":
			// Create new job or update existing one
			lm.jobs[jobKey] = &Job{
				Name:      entry.JobName,
				PID:       entry.PID,
				StartTime: entry.Timestamp,
				EndTime:   nil,
			}

		case "END":
			// Find corresponding job and update end time
			if job, exists := lm.jobs[jobKey]; exists {
				job.EndTime = &entry.Timestamp
				job.Duration = entry.Timestamp.Sub(job.StartTime)
			} else {
				log.Printf("Warning: Found END for job %s (PID: %d) without matching START",
					entry.JobName, entry.PID)
			}
		}
	}

	log.Printf("Processed %d unique jobs", len(lm.jobs))
}

// TODO: generate report

func main() {
	logFile := "logs.log"

	// Check if custom log file is provided
	if len(os.Args) > 1 {
		logFile = os.Args[1]
	}

	fmt.Printf("Starting log monitoring for file: %s\n", logFile)
	fmt.Println()

	// Create log monitor
	monitor := NewLogMonitor()

	// Parse log file
	if err := monitor.ParseLogFile(logFile); err != nil {
		log.Fatalf("Error parsing log file: %v", err)
	}
	// Process jobs to match START/END pairs
	monitor.ProcessJobs()

	// TODO: Generate and display report
}
