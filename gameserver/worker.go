package gameserver

import (
	"github.com/google/uuid"
)

// A worker takes a function func() from the JobQueue and runs it
type Worker struct {
	Id       string
	JobQueue JobQueue
}

// Return a new Worker object
func NewWorker(jobQueue JobQueue) *Worker {
	return &Worker{
		Id:       uuid.New().String(),
		JobQueue: jobQueue,
	}
}

// Start Worker process. The worker will take Jobs from the JobQueue and run them
func (w *Worker) Start() {
	for {
		gameProcess := <-w.JobQueue
		gameProcess.PreGameHook()
		gameProcess.MainGameProcessHook()
		gameProcess.PostGameHook()
	}
}
