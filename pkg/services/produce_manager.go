package services

import (
	"context"
	"sync"
	"time"
	"update-service/pkg/database"
	"update-service/pkg/models"

	"go.uber.org/zap"
)

// ProduceManager
type ProduceManager struct {
	log         *zap.Logger          // Logger to recodr internal events and errors
	serverTable database.ServerTable // Interface to store server in db
	outputChan  chan *models.Task
	inputChan   chan *models.Task
	ticker      *time.Ticker
}

func NewProduceManager(
	log *zap.Logger,
	serverTable database.ServerTable,
	outputChan chan *models.Task,
	delay time.Duration,
	inputLimmit int,
) *ProduceManager {
	return &ProduceManager{
		log:         log.With(zap.String("component", "Producer")),
		inputChan:   make(chan *models.Task, inputLimmit),
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
					inst.outputChan <- models.NewTask(&server, nil)
				}
			case Job := <-inst.inputChan:
				inst.outputChan <- Job
			case <-ctx.Done():
				close(inst.inputChan)
				close(inst.outputChan)
				inst.log.Info("Shutdown signal received. Stopping...")
				return
			}
		}
	}()
}

func (inst *ProduceManager) InputChan() chan *models.Task {
	return inst.inputChan
}
