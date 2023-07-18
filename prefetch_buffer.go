package main

import (
	"context"
	"go.uber.org/zap"
)

type PrefetchBuffer struct {
	ctx           context.Context
	logger        *zap.Logger
	msgInProgress chan<- Item
	inMemPq       *MinHeap
}

func NewPrefetchBuffer(ctx context.Context,
	logger *zap.Logger) *PrefetchBuffer {
	cInProgress := make(chan Item)
	p := &PrefetchBuffer{
		ctx:           ctx,
		logger:        logger,
		msgInProgress: cInProgress,
	}
	p.inMemPq = NewMinHeap()
	p.logger.Info("Starting PrefetchBuffer. Ready to serve!")
	return p
}

func (p *PrefetchBuffer) Start() {
	worker := NewDequeueWorker(p.ctx, p.logger)
	go func() {
		worker.Start()
		for i := range worker.dequeuedChan {
			p.inMemPq.Insert(i)
		}
	}()
}
