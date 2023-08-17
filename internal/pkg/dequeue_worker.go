package pkg

import (
	"context"
	"dps/internal/pkg/entity"
	"dps/logger"
	"fmt"
	"time"
)

type DequeueWorker struct {
	ctx          context.Context
	dequeuedChan chan<- entity.Item
}

func NewDequeueWorker(
	ctx context.Context,
	dequeuedChan chan<- entity.Item) *DequeueWorker {
	logger.Info("Success init DequeueWorker")
	return &DequeueWorker{
		ctx:          ctx,
		dequeuedChan: dequeuedChan,
	}
}

func (d *DequeueWorker) Start() {
	logger.Info("Dequeue Worker start polling job")
	go func() {
		for range time.Tick(time.Second * 1) {
			items, err := d.PullItemFromSource()
			if err != nil {
				logger.Error(fmt.Sprintf("[DequeueWorker] PullItemFromSource err: %v", err))
				return
			}
			for _, item := range items {
				d.dequeuedChan <- item
			}
		}
	}()

}

// PullItemFromSource first I will fake it
func (d *DequeueWorker) PullItemFromSource() ([]entity.Item, error) {
	// TODO: implement status
	return entity.GetStore().FetchItemReadyToDelivery(d.ctx, entity.ItemStatusInitialize)
}
