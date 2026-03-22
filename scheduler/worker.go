package scheduler

import (
	"fmt"
	"time"
)

//WorkerStatus represents the current state of the worker
type WorkerStatus int
const (
	NotBusy WorkerStatus = iota
	Working
)

//struct for the worker class
type Worker struct {
	ID int
	CurJob *Job
	Status WorkerStatus
	JobChan chan *Job
}

//worker gets jobs from its personal job channel, as it completes them it sends them to the results channel and updates the job
func worker(w *Worker, results chan<- *Job, updates chan<- *Worker) {
		updates <- w

	for j := range w.JobChan {
		w.CurJob = j
		w.Status = Working
		updates <- w

		fmt.Println("Worker", w.ID, "started Job", j.ID)
		j.Status = Running
		time.Sleep(j.Duration)
		fmt.Println("Worker", w.ID, "finished Job", j.ID, "-> Evaluated Expression =", j.ExpressionValue)
		j.Status = Completed

		w.CurJob = nil
		w.Status = NotBusy
		updates <- w
		results <- j

	}
}