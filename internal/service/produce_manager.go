package service

import (
	"context"
	"sync"
	"time"
	"update-service/internal/model"
	"update-service/internal/repository"

	"go.uber.org/zap"
)

// ProduceManager
type ProduceManager struct {
	log         *zap.Logger            // Logger to recodr internal events and errors
	serverTable repository.ServerTable // Interface to store server in db
	outputChan  chan *model.Task
	inputChan   chan *model.Task
	ticker      *time.Ticker
}

func NewProduceManager(
	log *zap.Logger,
	serverTable repository.ServerTable,
	outputChan chan *model.Task,
	delay time.Duration,
	inputLimmit int,
) *ProduceManager {
	return &ProduceManager{
		log:         log.With(zap.String("component", "Producer")),
		inputChan:   make(chan *model.Task, inputLimmit),
		serverTable: serverTable,
		outputChan:  outputChan,
		ticker:      time.NewTicker(delay),
	}
}

func (inst *ProduceManager) Produce(wg *sync.WaitGroup, ctx context.Context) {
	wg.Add(1)
	inst.log.Info("Start and Wait Task...")
	go func() {
		defer wg.Done()
		for {
			select {
			case <-inst.ticker.C:
				inst.log.Info("Start produce Task")
				servers, err := inst.serverTable.SelectAll()
				if err != nil {
					inst.log.Error("Failed get server", zap.Error(err))
					continue
				}

				for _, server := range servers {
					inst.outputChan <- model.NewTask(&server, nil)
				}

			case Job := <-inst.inputChan:
				inst.outputChan <- Job
			case <-ctx.Done():
				inst.log.Info("Shutdown signal received. Stopping...")
				return
			}
		}
	}()
}

func (inst *ProduceManager) InputChan() chan *model.Task {
	return inst.inputChan
}
