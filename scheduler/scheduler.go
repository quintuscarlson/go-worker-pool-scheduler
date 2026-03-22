package scheduler 
import (
	"fmt"
	"sort"
	"time"
	"github.com/Knetic/govaluate"
)

//struct for Scheduler class
type Scheduler struct {
	Results chan *Job
	numWorkers int
	jobTracker map[int]*Job
	workerTracker map[int]*Worker
	WorkerUpdates chan *Worker
	jobQueue []*Job
	startTime time.Time
	notOptTime time.Duration
}

//Scheduler constructor 
func NewScheduler(workerCount int) *Scheduler {
	return &Scheduler{
		Results: make(chan *Job, 10),
		numWorkers: workerCount,
		jobTracker: make(map[int]*Job),
		workerTracker: make(map[int]*Worker),
		WorkerUpdates: make(chan *Worker, workerCount),
		jobQueue: []*Job{},
		startTime: time.Now(),
		notOptTime: time.Duration(0),
		}
}

//creates a worker for each numWorkers and starts them, basically tells them to start watching their respective job channel
func (s *Scheduler) StartWorkers() {
	for i := 1; i <= s.numWorkers; i++ {
		w := &Worker{
			ID: i,
			CurJob: nil, 
			Status: NotBusy, 
			JobChan: make(chan *Job, 1),
		}
		s.workerTracker[i] = w
		go worker(w, s.Results, s.WorkerUpdates)
	}
}

//creates a new job using info received from the user in main
func (s *Scheduler) CreateJob(eq string, jobID int){
	expr, _ := govaluate.NewEvaluableExpression(eq)
	val, _ := expr.Evaluate(nil)

	cur := &Job{
			ID: jobID,
			Name: fmt.Sprintf("Job%d", jobID),
			ExpressionValue: val.(float64),
			Duration: s.findDur(eq),
			WaitTime: time.Duration(0),
		}
 	s.SubmitJob(cur)
}

//finds the duration by counting the number of operators (and starting at one)
func (s *Scheduler) findDur(eq string)time.Duration{
	var tm int = 1
	for _, ch := range eq {
		if (ch == '+' || ch == '-' || ch == '*' || ch == '/'){
			tm++
		}
	}
	return time.Duration(tm) * time.Second
}

//add a job to the jobsQueue
func (s *Scheduler) SubmitJob(job *Job) {
	job.Status = Queued 
	s.jobTracker[job.ID] = job
	s.jobQueue = append(s.jobQueue, job)
}

//calculates the total wait time if the scheduler didnt optimize the wait times by sorting the jobs in ascending order of duration
func (s *Scheduler) calcNotOptWTime()time.Duration{
	totalTime := time.Duration(0)
	curWait := make([]time.Duration, s.numWorkers)


	for i := 0; i < len(s.jobQueue); i++ {
		curLow := curWait[0]
		curWorker := 0
		for j := 1; j < s.numWorkers; j++ {
			if curWait[j] < curLow {
				curLow = curWait[j];
				curWorker = j
			}
		}
		totalTime += curWait[curWorker]
		curWait[curWorker] += s.jobQueue[i].Duration
	}

	return totalTime
}

//tells the scheduler to start assigning jobs and have the workers start working
func (s *Scheduler) RunWorkers() {
	s.notOptTime = s.calcNotOptWTime()
	s.SortJobs()
	s.startTime = time.Now()
	s.TrackWorkers()
	s.AssignJobs()
}

//sorts the jobQueue in ascending order of duration
func (s *Scheduler) SortJobs() {
	sort.Slice(s.jobQueue, func(i, j int)bool{
		return s.jobQueue[i].Duration < s.jobQueue[j].Duration
	})
}

//Assigns the next job in the Queue to a "NotBusy" worker by sending it to their personal job channel
func (s *Scheduler) AssignJobs() {
	for i := 1; i <= s.numWorkers; i++ {
		if(len(s.jobQueue) == 0){
			return
		}

		w := s.workerTracker[i]

		if (w.Status == NotBusy){
			j := s.jobQueue[0]
			j.WaitTime = time.Since(s.startTime)
			s.jobQueue = s.jobQueue[1:]
			w.Status = Working
			w.CurJob = j
			w.JobChan <- j
		}
	}
}

//constantly reads the worker channels so that the scheduler knows what workers are doing what and when
func (s *Scheduler) TrackWorkers() {
	go func() {
		for w := range s.WorkerUpdates {
			s.workerTracker[w.ID] = w

			if (w.Status == NotBusy){
				s.AssignJobs()
			}
		}
	} ()
}

//prints the current job status of every job submitted to the scheduler 
func (s *Scheduler) PrintJobsStatus(){
	fmt.Println("Current Job Statuses: ")

	for id, job := range s.jobTracker {
		found := false
		for id2, worker := range s.workerTracker {
			if(worker.CurJob == job) {
				fmt.Println("Job", id, ": Being worked on by Worker", id2)
				found = true
			}
		}
		if (!found) {
				fmt.Println("Job", id, ": ", job.Status)
			}
	}
	fmt.Println()
}

//closes the results channel and prints final remakrs, and statistics on the efficiency of the program
func (s *Scheduler) Stop(){
	for finish := 1; finish <= len(s.jobTracker); finish++ {
		 <-s.Results
	}
	optimizedWaitTime := time.Duration(0)

	for i:=1; i <= len(s.jobTracker); i++{
		optimizedWaitTime += s.jobTracker[i].WaitTime;
	}


	elapsedTime := time.Since(s.startTime)
	fmt.Println("\nAll Jobs Completed!\nTotal Runtime:", elapsedTime.Round(time.Second).Seconds(), "seconds.\nTotal Wait Time across all jobs with optimized Solution:", optimizedWaitTime.Round(time.Second).Seconds(), "seconds.\nTotal Wait Time across all jobs with non-optimized Solution:", s.notOptTime.Seconds(), "seconds.")

	for i := 1; i <= s.numWorkers; i++ {
		close(s.workerTracker[i].JobChan)
	}
}
