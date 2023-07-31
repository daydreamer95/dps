package pkg

import (
	"context"
	"dps/internal/pkg/entity"
	"dps/internal/pkg/logger"
	"math/rand"
	"time"
)

type DequeueWorker struct {
	ctx          context.Context
	dequeuedChan chan entity.Item
}

func NewDequeueWorker(
	ctx context.Context) *DequeueWorker {

	c := make(chan entity.Item)
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
			items := d.PullItemFromSource()
			for _, item := range items {
				d.dequeuedChan <- item
			}
		}
	}()

}

// PullItemFromSource first I will fake it
func (d *DequeueWorker) PullItemFromSource() []entity.Item {
	return []entity.Item{
		{Priority: rand.Int31()},
		{Priority: rand.Int31()},
	}
}
