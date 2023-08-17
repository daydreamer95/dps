package pkg

import (
	"context"
	"dps/internal/pkg/entity"
	"dps/logger"
)

type PrefetchBuffer struct {
	ctx     context.Context
	topicId uint
	inMemPq *entity.MinHeap
}

func NewPrefetchBuffer(ctx context.Context, topicId uint) *PrefetchBuffer {
	p := &PrefetchBuffer{
		ctx:     ctx,
		topicId: topicId,
	}
	p.inMemPq = entity.NewMinHeap()
	logger.Info("Starting PrefetchBuffer. Ready to serve!")
	return p
}
