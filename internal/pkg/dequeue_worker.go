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
	d.warmUp()
	logger.Info("Dequeue Worker start polling job")
	go func() {
		//TODO: config this tick
		for range time.Tick(time.Second * 1) {
			items, err := d.pullItemFromSource()
			if err != nil {
				logger.Error(fmt.Sprintf("[DequeueWorker] pullItemFromSource err: %v", err))
				continue
			}
			for i := 0; i < len(items); i++ {
				d.dequeuedChan <- items[i]
			}
		}
	}()

}

// pullItemFromSource
func (d *DequeueWorker) pullItemFromSource() ([]entity.Item, error) {
	out, err := entity.GetStore().GetItemByStatus(d.ctx, []string{entity.ItemStatusInitialize, entity.ItemStatusReadyToDeliver})
	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return out, nil
	}
	itemIds := make([]string, 0, len(out))
	for i := 0; i < len(out); i++ {
		itemIds = append(itemIds, out[i].Id)
	}

	err = entity.GetStore().UpdateItemsStatusByIds(d.ctx, itemIds, entity.ItemStatusDelivered)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (d *DequeueWorker) warmUp() {
	logger.Info("Start warm up dequeue worker:")
	out, err := entity.GetStore().GetItemByStatus(d.ctx, []string{entity.ItemStatusDelivered})
	if err != nil {
		return
	}
	for i := 0; i < len(out); i++ {
		d.dequeuedChan <- out[i]
	}
}
