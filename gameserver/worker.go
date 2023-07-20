package gameserver

import (
	"main/utils"
)

// A worker takes a function func() from the JobQueue and runs it
type Worker struct {
	utils.ID
	ActiveGame bool
	JobQueue   *utils.ModifiableQueue[GameProcess]
}

// Return a new Worker object
func NewWorker(jobQueue *utils.ModifiableQueue[GameProcess]) *Worker {
	return &Worker{
		ID:       *utils.NewId(),
		JobQueue: jobQueue,
	}
}

// Start Worker process. The worker will take Jobs from the JobQueue and run them
func (w *Worker) Start() {
	for {
		gameProcess := w.JobQueue.PopBlocking()

		w.ActiveGame = true
		gameProcess.PreGameHook()
		gameProcess.MainGameProcessHook()
		gameProcess.PostGameHook()
		w.ActiveGame = false
	}
}
