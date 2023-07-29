package dps_srv

import (
	"context"
	"dps/internal/pkg"
	"dps/internal/pkg/repository"
)

type RouterGrpc struct {
	UnimplementedDpsServiceServer
	rpw       pkg.IReplenishsesWorker
	itemRepo  *repository.ItemRepository
	topicRepo *repository.TopicRepository
}

func NewRouterGrpc(rpw pkg.IReplenishsesWorker,
	itemRepository *repository.ItemRepository,
	topicRepository *repository.TopicRepository) *RouterGrpc {
	return &RouterGrpc{
		rpw:       rpw,
		itemRepo:  itemRepository,
		topicRepo: topicRepository,
	}
}

func (d *RouterGrpc) Publish(context.Context, *PublishReq) (*PublishRes, error) {
	return nil, nil
}

func (d *RouterGrpc) CreateTopic(context.Context, *CreateTopicReq) (*CommonRes, error) {
	return nil, nil
}
func (d *RouterGrpc) Dequeue(ctx context.Context, req *DequeueReq) (*DequeueRes, error) {
	return nil, nil
}

func (d *RouterGrpc) Ack(context.Context, *AckReq) (*CommonRes, error) {
	return nil, nil
}

func (d *RouterGrpc) NAck(context.Context, *NAckReq) (*CommonRes, error) {
	return nil, nil
}
