package main

import (
	"context"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

type DequeueWorker struct {
	ctx          context.Context
	logger       *zap.Logger
	dequeuedChan chan Item
}

func NewDequeueWorker(
	ctx context.Context,
	logger *zap.Logger) *DequeueWorker {

	c := make(chan Item)
	logger.Info("Success init DequeueWorker")
	return &DequeueWorker{
		ctx:          ctx,
		logger:       logger,
		dequeuedChan: c,
	}
}

func (d *DequeueWorker) Start() {
	d.logger.Info("Dequeue Worker start polling job")
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
func (d *DequeueWorker) PullItemFromSource() []Item {
	return []Item{
		{Priority: rand.Int31()},
		{Priority: rand.Int31()},
	}
}
