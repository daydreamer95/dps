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
	rpw            IReplenishsesWorker
	topicProcessor ITopicProcessor
	itemProcessor  IItemProcessor
}

func NewRouterGrpc(rpw IReplenishsesWorker,
	topicProcessor ITopicProcessor,
	itemProcessor IItemProcessor) *RouterGrpc {
	return &RouterGrpc{
		rpw:            rpw,
		topicProcessor: topicProcessor,
		itemProcessor:  itemProcessor,
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

func (d *RouterGrpc) GetActiveTopics(ctx context.Context, req *empty.Empty) (*dps_pb.GetActiveTopicsRes, error) {
	topics, err := d.topicProcessor.GetActiveTopic(ctx)
	if err != nil {
		return &dps_pb.GetActiveTopicsRes{}, status.New(codes.Internal, err.Error()).Err()
	}
	out := &dps_pb.GetActiveTopicsRes{}
	for _, t := range topics {
		out.Topics = append(out.Topics, &dps_pb.Topic{
			Id:             uint32(t.Id),
			Name:           t.Name,
			DeliveryPolicy: t.DeliverPolicy,
		})
	}
	return out, nil
}

func (d *RouterGrpc) CreateTopic(ctx context.Context, req *dps_pb.CreateTopicReq) (*dps_pb.CreateTopicRes, error) {
	t := Topic{
		Name:          req.TopicName,
		Active:        uint(TopicStatusActive),
		DeliverPolicy: string(req.DeliverPolicy),
	}
	topic, err := d.topicProcessor.CreateTopic(ctx, t)
	if err != nil {
		return &dps_pb.CreateTopicRes{}, status.New(codes.Internal, err.Error()).Err()
	}

	return &dps_pb.CreateTopicRes{
		TopicId:        uint32(topic.Id),
		Name:           topic.Name,
		Active:         dps_pb.TopicActive(topic.Active),
		DeliveryPolicy: topic.DeliverPolicy,
	}, nil
}
