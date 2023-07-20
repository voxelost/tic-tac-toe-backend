package gameserver

import (
	"context"
	"main/utils"
)

// WorkerPool struct represents the array of Worker objects
type WorkerPool struct {
	Workers []*Worker
}

// return new Worker Pool
func NewWorkerPool(ctx context.Context, jobQueue *utils.ModifiableQueue[GameProcess], poolSize int) *WorkerPool {
	workers := []*Worker{}
	for i := 0; i < poolSize; i++ {
		newWorker := NewWorker(jobQueue)
		go newWorker.Start()
		workers = append(workers, newWorker)
	}

	return &WorkerPool{
		Workers: workers,
	}
}

func (wp *WorkerPool) GetActiveGames() int {
	c := 0
	for _, worker := range wp.Workers {
		if worker.ActiveGame {
			c++
		}
	}

	return c
}
