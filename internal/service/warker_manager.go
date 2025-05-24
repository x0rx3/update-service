package service

import (
	"context"
	"sync"
)

type WorkerManager struct {
	limit   int
	workers []Worker
}

func NewWorkerManager(limit int, workers []Worker) *WorkerManager {
	return &WorkerManager{
		limit:   limit,
		workers: workers,
	}
}

func (inst *WorkerManager) Build(wg *sync.WaitGroup, ctx context.Context) {
	for _, worker := range inst.workers {
		w := worker
		for i := 0; i < inst.limit; i++ {
			wg.Add(1)
			go func(w Worker) {
				defer wg.Done()
				w.Process(ctx)
			}(w)
		}
	}
}
