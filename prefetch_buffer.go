package main

import (
	"context"
	"dps/logger"
)

type PrefetchBuffer struct {
	ctx       context.Context
	topicName string
	inMemPq   *MinHeap
}

func NewPrefetchBuffer(ctx context.Context, topicName string) *PrefetchBuffer {
	p := &PrefetchBuffer{
		ctx:       ctx,
		topicName: topicName,
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
