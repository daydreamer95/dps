package main

import "go.uber.org/zap"

type PrefetchBuffer struct {
	logger        *zap.Logger
	msgIn         chan<- Item
	msgInProgress chan<- Item
	inMemPq       MinHeap
}

func NewPrefetchBuffer(logger *zap.Logger) *PrefetchBuffer {
	p := &PrefetchBuffer{
		logger: logger,
	}
	cIn := make(chan Item)
	cInProgress := make(chan Item)
	p.msgIn = cIn
	p.msgInProgress = cInProgress
	p.logger.Info("Starting PrefetchBuffer. Ready to serve!")
	return p
}
