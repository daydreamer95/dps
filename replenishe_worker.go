package main

import (
	"context"
	"dps/logger"
	"fmt"
)

type ReplenishesWorker struct {
	ctx             context.Context
	createTopicChan chan string
	deferItemChan   chan Item
	activeTopics    []string
	preBuffers      map[string]*PrefetchBuffer
}

func NewReplenishesWorker(ctx context.Context,
	createTopicChan chan string,
	deferItemChan chan Item) *ReplenishesWorker {
	out := &ReplenishesWorker{
		ctx:             ctx,
		createTopicChan: createTopicChan,
		deferItemChan:   deferItemChan,
	}
	out.preBuffers = map[string]*PrefetchBuffer{}
	return out
}

func (r *ReplenishesWorker) Start() {
	logger.Info("Replenishes_Worker start!")
	t := getActiveTopics()

	for _, topics := range t {
		pb := NewPrefetchBuffer(r.ctx)
		r.preBuffers[topics] = pb
		go pb.Start()
	}

	logger.Info("Start listen on createTopicChan and deferItemChan")
	for {
		select {
		case newTopics := <-r.createTopicChan:
			topic := r.preBuffers[newTopics]
			if topic != nil {
				logger.Fatal(fmt.Sprintf("Topics name [%v] already exists. Something wrong"))
				return
			}
		case deferItem := <-r.deferItemChan:
			logger.Info(fmt.Sprintf("An item has defered [%+v]", deferItem))
			return
		}
	}
}

// getActiveTopics Return topics name
func getActiveTopics() []string {
	return []string{"topic1", "topic2"}
}
