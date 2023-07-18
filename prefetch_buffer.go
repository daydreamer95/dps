package main

import (
	"context"
	"dps/logger"
)

type PrefetchBuffer struct {
	ctx           context.Context
	msgInProgress chan<- Item
	inMemPq       *MinHeap
}

func NewPrefetchBuffer(ctx context.Context) *PrefetchBuffer {
	cInProgress := make(chan Item)
	p := &PrefetchBuffer{
		ctx:           ctx,
		msgInProgress: cInProgress,
	}
	p.inMemPq = NewMinHeap()
	logger.Info("Starting PrefetchBuffer. Ready to serve!")
	return p
}

func (p *PrefetchBuffer) Start() {
	worker := NewDequeueWorker(p.ctx)
	go func() {
		worker.Start()
		for i := range worker.dequeuedChan {
			p.inMemPq.Insert(i)
		}
	}()
}
