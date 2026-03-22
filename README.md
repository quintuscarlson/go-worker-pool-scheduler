# Go Worker Pool Scheduler

A concurrent job scheduler built in Go that distributes equation-evaluation tasks across a worker pool and reduces total wait time using shortest-job-first scheduling.

## Overview

This project simulates a worker pool scheduler that accepts user-submitted arithmetic equations, converts each equation into a job, and assigns jobs to workers for concurrent execution.

Each job’s execution time is based on the complexity of the equation, measured by the number of operators it contains. Before assigning jobs, the scheduler sorts them by duration so shorter jobs run first, improving total wait time across all submitted jobs.

## Features

- Concurrent worker pool implementation in Go
- Dynamic job assignment to available workers
- Shortest-job-first scheduling strategy
- Equation parsing and evaluation using `govaluate`
- Per-job status tracking
- Wait time and runtime statistics
- Comparison between optimized and non-optimized total wait time

## How It Works

1. The user chooses how many workers to create.
2. The user enters arithmetic equations through the terminal.
3. Each equation becomes a job with:
   - a computed result
   - a duration based on operator count
   - a tracked wait time
4. Jobs are sorted by shortest duration first.
5. Available workers pull jobs from the scheduler and execute them concurrently.
6. When all jobs finish, the program prints:
   - total runtime
   - total optimized wait time
   - total non-optimized wait time

## Example Input

```
Enter how many workers you want to use: 3
Enter Equation Here: 1+2
Enter Equation Here: 5*4-3
Enter Equation Here: 10/2+6*3
Enter Equation Here: DONE
Worker 1 started Job 1
Worker 2 started Job 2
Worker 3 started Job 3
Worker 1 finished Job 1 -> Evaluated Expression = 3
...
All Jobs Completed!
Total Runtime: 5 seconds
Total Wait Time across all jobs with optimized Solution: 4 seconds
Total Wait Time across all jobs with non-optimized Solution: 7 seconds
```

## Project Structure
```
.
├── main.go
├── go.mod
├── go.sum
└── scheduler
    ├── job.go
    ├── scheduler.go
    └── worker.go
```

## Files
- main.go — handles user interaction and program startup
- scheduler.go — contains the scheduling logic, worker tracking, job sorting, and statistics
- worker.go — defines worker behavior and concurrent job execution
- job.go — defines job data and job status types


## Technologies Used
- Go
- Goroutines
- Channels
- govaluate

## Scheduling Strategy
This scheduler uses a shortest-job-first style optimization by sorting jobs in ascending order of duration before assigning them to workers.  

The goal is to reduce the total wait time across all jobs compared to a non-optimized first-come, first-served assignment strategy.

## Skills Demonstrated
- Concurrency in Go
- Goroutines and channels
- Scheduling and worker-pool design
- Data structures and state tracking
- CLI program design
- Performance-oriented thinking
## Possible Future Improvements
- Support for more complex mathematical expressions
- Better input validation and error handling
- Live status dashboard for workers and jobs
- Additional scheduling algorithms for comparison
- Unit tests for scheduler behavior

## Running the Project
1. Clone the Repository:
```
git clone https://github.com/quintuscarlson/go-worker-pool-scheduler.git
```

2. Move into the project folder:
```
cd go-worker-pool-scheduler
```

3. Run the program:
```
go run main.go
```

## Notes
This project was built as a systems/concurrency-focused scheduling simulation to demonstrate how worker pools and job prioritization can improve performance in Go applications.
