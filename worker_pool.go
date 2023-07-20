package main

import (
	"context"
)

type WorkerPool struct {
	Workers []Worker
}

func NewWorkerPool(ctx context.Context, jobQueue JobQueue, poolSize int) *WorkerPool {
	workers := []Worker{}
	for i := 0; i < poolSize; i++ {
		newWorker := NewWorker(jobQueue)
		go newWorker.Start(ctx)
		workers = append(workers, *newWorker)
	}

	return &WorkerPool{
		Workers: workers,
	}
}
