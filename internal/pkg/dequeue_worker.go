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
		//TODO: config this tick
		for range time.Tick(time.Second * 1) {
			items, err := d.PullItemFromSource()
			if err != nil {
				logger.Error(fmt.Sprintf("[DequeueWorker] PullItemFromSource err: %v", err))
				continue
			}
			for _, item := range items {
				d.dequeuedChan <- item
			}
		}
	}()

}

// PullItemFromSource first I will fake it
func (d *DequeueWorker) PullItemFromSource() ([]entity.Item, error) {
	out, err := entity.GetStore().FetchItemReadyToDelivery(d.ctx, entity.ItemStatusInitialize)
	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return out, nil
	}
	itemIds := make([]string, len(out))
	for i := 0; i < len(out); i++ {
		itemIds = append(itemIds, out[i].Id)
	}

	err = entity.GetStore().UpdateItemsStatusByIds(d.ctx, itemIds, entity.ItemStatusDelivered)
	if err != nil {
		return nil, err
	}
	return out, nil
}
