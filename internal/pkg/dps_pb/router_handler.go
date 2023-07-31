package dps_pb

import (
	"context"
	"dps/internal/pkg"
)

type RouterGrpc struct {
	UnimplementedDpsServiceServer
	rpw pkg.IReplenishsesWorker
}

func NewRouterGrpc(rpw pkg.IReplenishsesWorker) *RouterGrpc {
	return &RouterGrpc{
		rpw: rpw,
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
