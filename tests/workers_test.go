package tests

import (
	"context"
	"update-service/pkg/models"
)

type WorkerTest struct {
	id         string
	inputChan  chan *models.Task // Channel for receiving jobs to be processed
	outputChan chan *models.Task // Channel for sending processed jobs
}

func NewWorkerTest(id string, limit int) *WorkerTest {
	return &WorkerTest{
		id:         id,
		inputChan:  make(chan *models.Task, limit),
		outputChan: make(chan *models.Task, limit),
	}
}

func (inst *WorkerTest) Process(ctx context.Context) {
	for {
		select {
		case job := <-inst.inputChan:
			inst.outputChan <- job
		case <-ctx.Done():
			return
		}
	}
}

// InputChan returns the input channel for receiving jobs.
func (inst *WorkerTest) InputChan() chan *models.Task {
	return inst.inputChan
}

// OutputChan returns the output channel for sending completed jobs.
func (inst *WorkerTest) OutputChan() chan *models.Task {
	return inst.outputChan
}
