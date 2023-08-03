package pkg

import (
	"context"
	"dps/logger"
)

type PrefetchBuffer struct {
	ctx     context.Context
	topicId uint
	inMemPq *MinHeap
}

func NewPrefetchBuffer(ctx context.Context, topicId uint) *PrefetchBuffer {
	p := &PrefetchBuffer{
		ctx:     ctx,
		topicId: topicId,
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
