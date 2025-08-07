package main

import (
	"os"
	"testing"
	"time"
)

func TestLogMonitor_ParseLogFile(t *testing.T) {
	// Create a temporary test log file
	testData := `11:35:23,scheduled task 032, START,37980
11:35:56,scheduled task 032, END,37980
11:36:11,scheduled task 796, START,57672
11:36:18,scheduled task 796, END,57672
11:36:58,background job wmy, START,81258
11:51:44,background job wmy, END,81258`

	tmpFile, err := os.CreateTemp("", "test_log_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(testData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tmpFile.Close()

	// Test parsing
	monitor := NewLogMonitor()
	err = monitor.ParseLogFile(tmpFile.Name())
	if err != nil {
		t.Errorf("ParseLogFile failed: %v", err)
	}

	// Verify entries were parsed correctly
	expectedEntries := 6
	if len(monitor.entries) != expectedEntries {
		t.Errorf("Expected %d entries, got %d", expectedEntries, len(monitor.entries))
	}

	// Check first entry
	if monitor.entries[0].JobName != "scheduled task 032" {
		t.Errorf("Expected job name 'scheduled task 032', got '%s'", monitor.entries[0].JobName)
	}

	if monitor.entries[0].Action != "START" {
		t.Errorf("Expected action 'START', got '%s'", monitor.entries[0].Action)
	}

	if monitor.entries[0].PID != 37980 {
		t.Errorf("Expected PID 37980, got %d", monitor.entries[0].PID)
	}
}
func TestLogMonitor_ProcessJobs(t *testing.T) {
	monitor := NewLogMonitor()

	// Add test entries manually
	startTime, _ := time.Parse("15:04:05", "11:35:23")
	endTime, _ := time.Parse("15:04:05", "11:35:56")
	longEndTime, _ := time.Parse("15:04:05", "11:51:44") // 15+ minutes later

	monitor.entries = []LogEntry{
		{Timestamp: startTime, JobName: "test task", Action: "START", PID: 12345},
		{Timestamp: endTime, JobName: "test task", Action: "END", PID: 12345},
		{Timestamp: startTime, JobName: "long task", Action: "START", PID: 54321},
		{Timestamp: longEndTime, JobName: "long task", Action: "END", PID: 54321},
		{Timestamp: startTime, JobName: "incomplete task", Action: "START", PID: 99999},
	}

	monitor.ProcessJobs()

	// Check that jobs were processed correctly
	if len(monitor.jobs) != 3 {
		t.Errorf("Expected 3 jobs, got %d", len(monitor.jobs))
	}

	// Check completed job
	jobKey := "test task_12345"
	if job, exists := monitor.jobs[jobKey]; exists {
		if job.EndTime == nil {
			t.Error("Expected job to have end time")
		}
		expectedDuration := 33 * time.Second // 11:35:56 - 11:35:23
		if job.Duration != expectedDuration {
			t.Errorf("Expected duration %v, got %v", expectedDuration, job.Duration)
		}
	} else {
		t.Error("Expected job not found")
	}

	// Check incomplete job
	incompleteKey := "incomplete task_99999"
	if job, exists := monitor.jobs[incompleteKey]; exists {
		if job.EndTime != nil {
			t.Error("Expected incomplete job to have no end time")
		}
	} else {
		t.Error("Expected incomplete job not found")
	}
}
func TestNewLogMonitor(t *testing.T) {
	monitor := NewLogMonitor()

	if monitor == nil {
		t.Error("NewLogMonitor returned nil")
	}

	if monitor.entries == nil {
		t.Error("entries slice not initialized")
	}

	if monitor.jobs == nil {
		t.Error("jobs map not initialized")
	}

	if len(monitor.entries) != 0 {
		t.Errorf("Expected empty entries slice, got length %d", len(monitor.entries))
	}

	if len(monitor.jobs) != 0 {
		t.Errorf("Expected empty jobs map, got length %d", len(monitor.jobs))
	}
}
