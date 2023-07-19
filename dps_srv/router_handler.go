package dps_srv

import (
	"context"
)

type RouterGrpc struct {
	UnimplementedDpsServiceServer
}

func NewRouterGrpc() *RouterGrpc {
	return &RouterGrpc{}
}

func (d *RouterGrpc) Publish(context.Context, *PublishReq) (*PublishRes, error) {
	return nil, nil
}

func (d *RouterGrpc) CreateTopic(context.Context, *CreateTopicReq) (*CommonRes, error) {
	return nil, nil
}
func (d *RouterGrpc) Dequeue(context.Context, *DequeueReq) (*DequeueRes, error) {
	return nil, nil
}

func (d *RouterGrpc) Ack(context.Context, *AckReq) (*CommonRes, error) {
	return nil, nil
}

func (d *RouterGrpc) NAck(context.Context, *NAckReq) (*CommonRes, error) {
	return nil, nil
}
