package pkg

import (
	"context"
	"dps/internal/pkg/dps_pb"
	"github.com/golang/protobuf/ptypes/empty"
)

type RouterGrpc struct {
	dps_pb.UnimplementedDpsServiceServer
	rpw IReplenishsesWorker
}

func NewRouterGrpc(rpw IReplenishsesWorker) *RouterGrpc {
	return &RouterGrpc{
		rpw: rpw,
	}
}

func (d *RouterGrpc) Publish(context.Context, *dps_pb.PublishReq) (*dps_pb.PublishRes, error) {
	return nil, nil
}

func (d *RouterGrpc) CreateTopic(context.Context, *dps_pb.CreateTopicReq) (*dps_pb.GetActiveTopicsRes, error) {
	return nil, nil
}
func (d *RouterGrpc) Dequeue(ctx context.Context, req *dps_pb.DequeueReq) (*dps_pb.DequeueRes, error) {
	return nil, nil
}

func (d *RouterGrpc) Ack(context.Context, *dps_pb.AckReq) (*dps_pb.CommonRes, error) {
	return nil, nil
}

func (d *RouterGrpc) NAck(context.Context, *dps_pb.NAckReq) (*dps_pb.CommonRes, error) {
	return nil, nil
}

func (d *RouterGrpc) GetActiveTopics(context.Context, *empty.Empty) (*dps_pb.GetActiveTopicsRes, error) {
	return nil, nil
}
