package pkg

import (
	"context"
	"dps/internal/pkg/common"
	"dps/internal/pkg/entity"
	"dps/logger"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// time interval for calling randomExpire and rotateExpire methods
const expireInterval = 100 * time.Millisecond

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
	expiredChan      chan entity.Item
	// preBuffers Maps of prefetch buffers with topic id was key
	preBuffers map[uint]*PrefetchBuffer
}

func NewReplenishesWorker(ctx context.Context,
	updateTopicChan <-chan entity.Topic,
	expireChan chan entity.Item,
	deferItemChan chan entity.Item) *ReplenishesWorker {
	out := &ReplenishesWorker{
		ctx:              ctx,
		updateTopicChan:  updateTopicChan,
		dequeuedItemChan: deferItemChan,
		expiredChan:      expireChan,
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

	r.mustStartLeaseWorker()

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
			logger.Info(fmt.Sprintf("[ReplenishesWorker] An item Id [%+v] has pull from db.", dequeuedItem.Id))
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
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := 0; i < len(items); i++ {
		pfBuffer := r.preBuffers[items[i].TopicId]
		if pfBuffer == nil {
			logger.Error(fmt.Sprintf("Not existed topic with id: [%v]", items[i].TopicId))
			continue
		}
		items[i].LeaseAfter = time.Now().Add(time.Duration(items[i].LeaseDuration) * time.Second)
		pfBuffer.inMemPq.Insert(items[i])
		logger.Info(fmt.Sprintf("Insert item to prefetch buffer done: Id [%v] TopicId :[%v] Priority [%v] Status [%v] and Leash after [%v]",
			items[i].Id, items[i].TopicId, items[i].Priority, items[i].Status, items[i].LeaseAfter))
	}
	return true, nil
}

func (r *ReplenishesWorker) Pop(topicId uint, count int) ([]entity.Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	pfBuffer := r.preBuffers[topicId]
	if pfBuffer == nil {
		return nil, common.ErrNotFoundTopic
	}

	var result []entity.Item
	for i := 0; i < count; i++ {
		polled, err := pfBuffer.inMemPq.Poll()
		if err != nil {
			logger.Info(fmt.Sprintf("Error poll item from prefetch buffer: [%v] ", err))
			if len(result) != 0 {
				return result, nil
			}
			break
		}

		if polled.LeaseAfter.After(time.Now()) {
			logger.Info(fmt.Sprintf("Passive lease item id [%v] cause lease time [%v]", polled.Id, polled.LeaseAfter))
		}

		logger.Info(fmt.Sprintf("Popped message Id [%v] Priority [%v] DeliveryAfter: [%v]", polled.Id, polled.Priority, polled.DeliverAfter))
		result = append(result, polled)
	}
	return result, nil
}

func (r *ReplenishesWorker) mustStartLeaseWorker() {
	for topicId, _ := range r.preBuffers {
		id := topicId
		go func() {
			for {
				start := time.Now()
				for i := 0; i < 10; i++ {
					if !r.randomExpire(id) {
						break
					}
				}
				diff := time.Since(start)
				time.Sleep(expireInterval - diff)
			}
		}()
	}
}

func (r *ReplenishesWorker) randomExpire(topicId uint) bool {
	start := time.Now()
	defer logger.Debug(fmt.Sprintf("RandomExpire topicId %v stop after %v", topicId, time.Since(start)))
	logger.Debug(fmt.Sprintf("Stop the world for random expire key for topic id %v", topicId))
	const totalChecks = 20
	r.mu.Lock()
	defer r.mu.Unlock()

	expiredFound := 0
	for i := 0; i < totalChecks; i++ {
		sz := len(r.preBuffers[topicId].inMemPq.Data)
		if sz == 0 {
			return false
		}
		item := r.preBuffers[topicId].inMemPq.Data[rand.Intn(sz)]

		if time.Now().After(item.LeaseAfter) {
			r.preBuffers[topicId].inMemPq.Delete(item)
			logger.Info(fmt.Sprintf("Random expire delete item [%v] cause lease time [%v]", item.Id, item.LeaseAfter))
			r.expiredChan <- item
			expiredFound++
		}
	}

	if expiredFound*4 >= totalChecks {
		return true
	}
	return false
}
