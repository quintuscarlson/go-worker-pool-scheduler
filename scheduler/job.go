package scheduler

import (
	"time"
)

//JobStatus represents the current state of the job
type JobStatus int
const (
	Queued JobStatus = iota
	Running
	Completed
	Failed
)

//struct for job objects
type Job struct {
	ID int
	Name string
	ExpressionValue float64
	Duration time.Duration
	Status JobStatus
	WaitTime time.Duration
}