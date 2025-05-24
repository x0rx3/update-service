package tests

import (
	"context"
	"update-service/internal/model"
)

type WorkerTest struct {
	id         string
	inputChan  chan *model.Task // Channel for receiving jobs to be processed
	outputChan chan *model.Task // Channel for sending processed jobs
}

func NewWorkerTest(id string, limit int) *WorkerTest {
	return &WorkerTest{
		id:         id,
		inputChan:  make(chan *model.Task, limit),
		outputChan: make(chan *model.Task, limit),
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
func (inst *WorkerTest) InputChan() chan *model.Task {
	return inst.inputChan
}

// OutputChan returns the output channel for sending completed jobs.
func (inst *WorkerTest) OutputChan() chan *model.Task {
	return inst.outputChan
}
