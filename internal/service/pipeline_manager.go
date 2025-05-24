package service

import (
	"context"
	"sync"
	"update-service/internal/model"

	"go.uber.org/zap"
)

// Pipeline represents a processing pipeline consisting of chained workers.
type PipelineManager struct {
	log    *zap.Logger
	worker []Worker // List of workers connected in a pipeline
}

// NewPipeline creates a new Pipeline instance with the given list of workers.
func NewPipeline(log *zap.Logger, workers []Worker) *PipelineManager {
	return &PipelineManager{log: log.With(zap.String("component", "Pipeline")), worker: workers}
}

// Build connects the output of each worker to the input of the next one using goroutines.
// It accepts a WaitGroup to wait for all connections and a context to allow graceful cancellation.
func (inst *PipelineManager) Build(wg *sync.WaitGroup, ctx context.Context) {
	inst.log.Info("Start and Wait Task...")
	// Ensure we connect each pair of workers: worker[i] -> worker[i+1]
	for i := 0; i < len(inst.worker)-1; i++ {
		wg.Add(1)

		// Get channels for the current and next worker
		outputChan := inst.worker[i].OutputChan()
		inputChan := inst.worker[i+1].InputChan()

		go func(out chan *model.Task, in chan *model.Task) {
			defer wg.Done()

			for {
				select {
				case job := <-out:
					in <- job
				case <-ctx.Done():
					inst.log.Info("Shutdown signal received. Stopping...")
					return
				}
			}
		}(outputChan, inputChan)
	}
}
