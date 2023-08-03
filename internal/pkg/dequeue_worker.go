package pkg

import (
	"context"
	"dps/logger"
	"fmt"
	"time"
)

type DequeueWorker struct {
	ctx          context.Context
	dequeuedChan chan Item
}

func NewDequeueWorker(
	ctx context.Context) *DequeueWorker {

	c := make(chan Item)
	logger.Info("Success init DequeueWorker")
	return &DequeueWorker{
		ctx:          ctx,
		dequeuedChan: c,
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
func (d *DequeueWorker) PullItemFromSource() ([]Item, error) {
	// TODO: implement status
	return GetStore().FetchItemReadyToDelivery("")
}
