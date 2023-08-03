package pkg

import (
	"context"
	"dps/logger"
	"fmt"
	"sync"
)

type IReplenishsesWorker interface {
	Start()
	Push(items []Item) (bool, error)
	//Pop Get item responding to request dto. Simply get items from
	// prefetch buffers and returns
	Pop() ([]Item, error)
}

type ReplenishesWorker struct {
	ctx             context.Context
	mu              sync.Mutex
	createTopicChan <-chan Topic
	deferItemChan   chan Item
	// preBuffers Maps of prefetch buffers with topic id was key
	preBuffers map[uint]*PrefetchBuffer
}

func NewReplenishesWorker(ctx context.Context,
	createTopicChan <-chan Topic,
	deferItemChan chan Item) *ReplenishesWorker {
	out := &ReplenishesWorker{
		ctx:             ctx,
		createTopicChan: createTopicChan,
		deferItemChan:   deferItemChan,
	}
	out.preBuffers = map[uint]*PrefetchBuffer{}
	return out
}

func (r *ReplenishesWorker) Start() {
	logger.Info("Replenishes_Worker start!")
	t, err := GetStore().GetActiveTopic()
	if err != nil {
		logger.Fatal(fmt.Sprintf("[ReplenishesWorker]"))
	}

	for _, topic := range t {
		logger.Info(fmt.Sprintf("[ReplenishesWorker] Init prefetch buffer topic [%v]", topic))
		pb := NewPrefetchBuffer(r.ctx, topic.Id)
		r.preBuffers[topic.Id] = pb
	}
	logger.Info("[ReplenishesWorker] Start listen on createTopicChan and deferItemChan")
	for {
		select {
		case newTopics := <-r.createTopicChan:
			r.mu.Lock()
			logger.Info(fmt.Sprintf("[ReplenishesWorker] Receive new topic from chan [%v]", newTopics))
			topic := r.preBuffers[newTopics.Id]
			if topic != nil {
				logger.Fatal(fmt.Sprintf("[ReplenishesWorker] Topics name [%v] already exists. Something wrong", topic))
				return
			}
			logger.Info(fmt.Sprintf("[ReplenishesWorker] Init prefetch buffer topicname [%v]", newTopics))
			pb := NewPrefetchBuffer(r.ctx, newTopics.Id)
			r.preBuffers[newTopics.Id] = pb
			r.mu.Unlock()
		case deferItem := <-r.deferItemChan:
			logger.Info(fmt.Sprintf("[ReplenishesWorker] An item has defered [%+v]", deferItem))
			//TODO: do this
			return
		}
	}
}

func (r *ReplenishesWorker) Push(items []Item) (bool, error) {
	return false, nil
}

func (r *ReplenishesWorker) Pop() ([]Item, error) {
	return nil, nil
}
