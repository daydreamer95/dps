package pkg

import (
	"context"
	"dps/internal/pkg/dps_pb"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (d *RouterGrpc) Dequeue(ctx context.Context, req *dps_pb.DequeueReq) (*empty.Empty, error) {
	return nil, nil
}

func (d *RouterGrpc) Ack(context.Context, *dps_pb.AckReq) (*empty.Empty, error) {
	return nil, nil
}

func (d *RouterGrpc) NAck(context.Context, *dps_pb.NAckReq) (*empty.Empty, error) {
	return nil, nil
}

func (d *RouterGrpc) GetActiveTopics(context.Context, *empty.Empty) (*dps_pb.GetActiveTopicsRes, error) {
	topics, err := GetStore().GetActiveTopic()
	if err != nil {
		return &dps_pb.GetActiveTopicsRes{}, status.New(codes.Internal, err.Error()).Err()
	}
	out := &dps_pb.GetActiveTopicsRes{}
	for _, t := range topics {
		out.Topics = append(out.Topics, &dps_pb.Topic{
			Id:             string(t.Id),
			Name:           t.Name,
			DeliveryPolicy: t.DeliverPolicy,
		})
	}
	return nil, nil
}

func (d *RouterGrpc) CreateTopic(ctx context.Context, req *dps_pb.CreateTopicReq) (*dps_pb.CreateTopicRes, error) {
	t := Topic{
		Name:          req.TopicName,
		Active:        TopicStatusActive,
		DeliverPolicy: string(req.DeliverPolicy),
	}
	topic, err := GetStore().CreateTopic(t)
	if err != nil {
		return &dps_pb.CreateTopicRes{}, status.New(codes.Internal, err.Error()).Err()
	}

	var status dps_pb.TopicActive
	if topic.Active == TopicStatusActive {
		status = dps_pb.TopicActive_ACTIVE
	} else {
		status = dps_pb.TopicActive_INACTIVE
	}

	return &dps_pb.CreateTopicRes{
		TopicId:        string(topic.Id),
		Name:           topic.Name,
		Active:         status,
		DeliveryPolicy: topic.DeliverPolicy,
	}, nil
}
