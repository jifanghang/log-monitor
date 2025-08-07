# Log Monitor

A Go app that monitors log files to track job execution times and generate warnings/errors.

## Overview

This application parses CSV log files containing job start/end events and calculates job durations:
- **Warnings** (⚠️): jobs taking longer than 5 minutes
- **Errors** (🚨): jobs taking longer than 10 minutes
- **OK** (✅): jobs completing within acceptable time
- Summary statistics and incomplete jobs

## Core features:
- **CSV Log Parsing**: Reads log files in CSV format with timestamp, job description, action, and PID
- **Job Matching**: Automatically matches START/END events by job name and PID
- **Duration Calculation**: Calculates job execution times
- **Reporting**: Console output with emojis and summary file generation
- **Unit Testing**: Full test coverage 

## Installation & Usage

### Prerequisites
- Go 1.21 or later

### Running the Application

1. **Clone or download the code**
2. **Prepare your log file** (default: `logs.log`)
3. **Run the application**:

```bash
# Using default log file (logs.log)
go run main.go

# Using custom log file
go run main.go path/to/your/logfile.csv

# Build and run executable
go build -o log-monitor
./log-monitor logs.log
```

### Running Tests

```bash
# Run tests with verbose output
go test -v

# Run tests with coverage
go test -cover
```

## Log File Format

Following structure expected:

```
HH:MM:SS,job_description,action,pid
```

**Example:**
```csv
11:35:23,scheduled task 032, START,37980
11:35:56,scheduled task 032, END,37980
11:36:11,scheduled task 796, START,57672
11:36:18,scheduled task 796, END,57672
```

**Fields:**
- **Timestamp**: Time in HH:MM:SS format (24-hour)
- **Job Description**: Name/description of the job or task
- **Action**: Either "START" or "END" (spaces are automatically trimmed)
- **PID**: Process ID (integer)

## Output

### Console Report
Detailed console report includes:
- Job status with emojis (✅ OK, ⚠️ WARNING, 🚨 ERROR), with names, PIDs, and execution durations
- Running jobs (START without matching END)
- Summary stats

## Example Usage & Output
```
➜  log-monitor git:(main) ✗ go run main.go
Starting log monitoring for file: logs.log

2025/08/07 17:11:14 Successfully parsed 88 log entries from logs.log
2025/08/07 17:11:14 Processed 45 unique jobs
===== LOG MONITORING REPORT =====

🚨 ERROR:   scheduled task 051        (PID: 39547) - Duration: 11m29s (>10min)
✅ OK:      scheduled task 538        (PID: 26831) - Duration: 2m12s
⚠️  WARNING: scheduled task 811        (PID: 50295) - Duration: 6m35s (>5min)
🚨 ERROR:   scheduled task 936        (PID: 62401) - Duration: 10m24s (>10min)
✅ OK:      background job ulp        (PID: 60134) - Duration: 44s
🚨 ERROR:   scheduled task 182        (PID: 70808) - Duration: 33m43s (>10min)
🚨 ERROR:   scheduled task 064        (PID: 85742) - Duration: 12m17s (>10min)
⚠️  WARNING: scheduled task 672        (PID: 24482) - Duration: 8m36s (>5min)

SUMMARY: 8 completed, 0 running, 2 warnings, 4 errors

2025/08/07 17:11:14 Report saved to monitoring_report.txt
```

Sample output can be found at `sample_output.txt`
