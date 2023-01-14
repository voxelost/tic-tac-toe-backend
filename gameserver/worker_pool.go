package gameserver

import (
	"context"
)

// WorkerPool struct represents the array of Worker objects
type WorkerPool struct {
	Workers []Worker
}

// return new Worker Pool
func NewWorkerPool(ctx context.Context, jobQueue JobQueue, poolSize int) *WorkerPool {
	workers := []Worker{}
	for i := 0; i < poolSize; i++ {
		newWorker := NewWorker(jobQueue)
		go newWorker.Start()
		workers = append(workers, *newWorker)
	}

	return &WorkerPool{
		Workers: workers,
	}
}
