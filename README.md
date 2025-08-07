A Go app that monitors log files to track job execution times and generate warnings/errors.

Core features:
- **CSV Log Parsing**: Reads log files in CSV format with timestamp, job description, action, and PID
- **Job Matching**: Automatically matches START/END events by job name and PID
- **Duration Calculation**: Calculates job execution times

## How to run
```bash
# Using default log file (logs.log)
go run main.go

# Using custom log file
go run main.go path/to/your/logfile.csv

# Build and run executable
go build -o log-monitor
./log-monitor logs.log
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
