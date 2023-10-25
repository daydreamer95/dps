package pkg

import (
	"context"
	"dps/internal/pkg/common"
	"dps/internal/pkg/entity"
	"dps/logger"
	"fmt"
	"sync"
)

type IReplenishsesWorker interface {
	Start()
	Push(items []entity.Item) (bool, error)
	//Pop Get item responding to request dto. Simply get items from
	// prefetch buffers and returns
	Pop(topicId uint, count int) ([]entity.Item, error)
}

type ReplenishesWorker struct {
	ctx              context.Context
	mu               sync.Mutex
	updateTopicChan  <-chan entity.Topic
	dequeuedItemChan chan entity.Item
	// preBuffers Maps of prefetch buffers with topic id was key
	preBuffers map[uint]*PrefetchBuffer
}

func NewReplenishesWorker(ctx context.Context,
	updateTopicChan <-chan entity.Topic,
	deferItemChan chan entity.Item) *ReplenishesWorker {
	out := &ReplenishesWorker{
		ctx:              ctx,
		updateTopicChan:  updateTopicChan,
		dequeuedItemChan: deferItemChan,
	}
	out.preBuffers = map[uint]*PrefetchBuffer{}
	return out
}

func (r *ReplenishesWorker) Start() {
	logger.Info("Replenishes_Worker start!")
	t, err := entity.GetStore().GetActiveTopic(r.ctx)
	if err != nil {
		logger.Fatal(fmt.Sprintf("[ReplenishesWorker]"))
	}

	for _, topic := range t {
		logger.Info(fmt.Sprintf("[ReplenishesWorker] Init prefetch buffer topic [%+v]", topic))
		pb := NewPrefetchBuffer(r.ctx, topic.Id)
		r.preBuffers[topic.Id] = pb
	}
	logger.Info("[ReplenishesWorker] Start listen on updateTopicChan and dequeuedItemChan")
	for {
		select {
		case updatedTopic := <-r.updateTopicChan:
			r.mu.Lock()
			logger.Info(fmt.Sprintf("[ReplenishesWorker] Receive updated from chan [%v]", updatedTopic))
			topic := r.preBuffers[updatedTopic.Id]
			if topic != nil {
				logger.Error(fmt.Sprintf("[ReplenishesWorker] Topics name [%v] already exists. Something wrong", topic))
				break
			}
			logger.Info(fmt.Sprintf("[ReplenishesWorker] Init prefetch buffer topicname [%v]", updatedTopic))
			pb := NewPrefetchBuffer(r.ctx, updatedTopic.Id)
			r.preBuffers[updatedTopic.Id] = pb
			r.mu.Unlock()
		case dequeuedItem := <-r.dequeuedItemChan:
			logger.Info(fmt.Sprintf("[ReplenishesWorker] An item has pull from db [%+v]", dequeuedItem))
			pb := r.preBuffers[dequeuedItem.TopicId]
			if pb == nil {
				logger.Error(fmt.Sprintf("[ReplenishesWorker] Topics Id [%v] not exists. Something wrong", dequeuedItem.TopicId))
				break
			}
			_, err = r.Push([]entity.Item{dequeuedItem})
			if err != nil {
				logger.Error(fmt.Sprintf("[ReplenishesWorker] Push item to in-memory Priority Queues cause an err [%v]. Dequeued-Item [%v]", err, dequeuedItem.TopicId))
				break
			}
		}
	}
}

func (r *ReplenishesWorker) Push(items []entity.Item) (bool, error) {
	for _, item := range items {
		pfBuffer := r.preBuffers[item.TopicId]
		if pfBuffer == nil {
			logger.Error(fmt.Sprintf("Not existed topic with id: [%v]", item.TopicId))
			continue
		}
		pfBuffer.inMemPq.Insert(item)
		logger.Info(fmt.Sprintf("Insert item to prefetch buffer done: [%+v]", item))
	}
	return true, nil
}

func (r *ReplenishesWorker) Pop(topicId uint, count int) ([]entity.Item, error) {
	pfBuffer := r.preBuffers[topicId]
	if pfBuffer == nil {
		return nil, common.ErrNotFoundTopic
	}

	var result []entity.Item
	for i := 0; i < count; i++ {
		polled, err := pfBuffer.inMemPq.Poll()
		if err != nil {
			logger.Info(fmt.Sprintf("Error poll item from prefetch buffer: [%v] ", err))
			continue
		}
		result = append(result, polled)
	}
	return result, nil
}
