package pkg

import (
	"context"
	"dps/internal/pkg/dps_pb"
	"dps/internal/pkg/entity"
	"dps/logger"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type RouterGrpc struct {
	dps_pb.UnimplementedDpsServiceServer
	rpw            IReplenishsesWorker
	topicProcessor entity.ITopicProcessor
	itemProcessor  entity.IItemProcessor
}

func NewRouterGrpc(rpw IReplenishsesWorker,
	topicProcessor entity.ITopicProcessor,
	itemProcessor entity.IItemProcessor) *RouterGrpc {
	return &RouterGrpc{
		rpw:            rpw,
		topicProcessor: topicProcessor,
		itemProcessor:  itemProcessor,
	}
}

func (d *RouterGrpc) Publish(ctx context.Context, req *dps_pb.PublishReq) (*dps_pb.PublishRes, error) {
	topic, err := d.topicProcessor.GetTopicByName(ctx, req.GetItem().TopicName)
	if err != nil {
		logger.Error(fmt.Sprintf("GRPC/Publish occur error [%v]. Notfound topic with name [%v]", err.Error(), req.GetItem().TopicName))
		return &dps_pb.PublishRes{}, status.New(codes.Internal, fmt.Sprintf("GRPC/Publish occur error [%v]. Notfound topic with name [%v]", err.Error(), req.GetItem().TopicName)).Err()
	}

	item := entity.Item{
		TopicId:       topic.Id,
		Priority:      req.Item.Priority,
		DeliverAfter:  time.Unix(req.Item.DeliverAfter, 0),
		Payload:       req.Item.Payload,
		MetaData:      req.Item.Metadata,
		LeaseDuration: 60, // TODO: allow to push this config
	}
	createItem, err := d.itemProcessor.CreateItem(ctx, item)
	if err != nil {
		logger.Error(fmt.Sprintf("GRPC/Publish occur error [%v]", err.Error()))
		return &dps_pb.PublishRes{}, status.New(codes.Internal, fmt.Sprintf("GRPC/Publish occur error [%v]", err.Error())).Err()
	}

	return &dps_pb.PublishRes{
		Id:           createItem.Id,
		TopicId:      uint32(createItem.TopicId),
		Priority:     createItem.Priority,
		Payload:      createItem.Payload,
		Metadata:     createItem.MetaData,
		DeliverAfter: createItem.DeliverAfter.Unix(),
		Status:       createItem.Status,
	}, nil

}

func (d *RouterGrpc) Dequeue(ctx context.Context, req *dps_pb.DequeueReq) (*dps_pb.DequeueRes, error) {
	topic, err := d.topicProcessor.GetTopicByName(ctx, req.TopicName)
	if err != nil {
		logger.Error(fmt.Sprintf("GRPC/Publish occur error [%v]. Notfound topic with name [%v]", err.Error(), req.TopicName))
		return &dps_pb.DequeueRes{}, status.New(codes.Internal, err.Error()).Err()
	}

	dequeItems, err := d.rpw.Pop(topic.Id, int(req.Count))
	if err != nil {
		logger.Error(fmt.Sprintf("GRPC/Dequeue occur error [%v]", err.Error()))
		return &dps_pb.DequeueRes{}, status.New(codes.Internal, err.Error()).Err()
	}
	var out []*dps_pb.ItemRes
	for _, deqItem := range dequeItems {
		out = append(out, &dps_pb.ItemRes{
			Id:            deqItem.Id,
			TopicId:       uint32(deqItem.TopicId),
			Priority:      deqItem.Priority,
			Payload:       deqItem.Payload,
			Metadata:      deqItem.MetaData,
			DeliverAfter:  deqItem.DeliverAfter.Unix(),
			Status:        deqItem.Status,
			LeaseDuration: uint32(deqItem.LeaseDuration),
		})
	}

	return &dps_pb.DequeueRes{Items: out}, nil
}

func (d *RouterGrpc) Ack(ctx context.Context, req *dps_pb.AckReq) (*empty.Empty, error) {
	topic, err := d.topicProcessor.GetTopicByName(ctx, req.GetTopic())
	if err != nil {
		logger.Error(fmt.Sprintf("GRPC/Publish occur error [%v]. Notfound topic with name [%v]", err.Error(), req.GetTopic()))
		return nil, status.New(codes.NotFound, err.Error()).Err()
	}

	err = d.itemProcessor.Delete(ctx, topic.Id, req.GetDpsAssignedUniqueId())
	if err != nil {
		logger.Error(fmt.Sprintf("GRPC/Publish ack error error [%v]. Notfound topic with id [%v]", err.Error(), req.GetDpsAssignedUniqueId()))
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (d *RouterGrpc) NAck(ctx context.Context, req *dps_pb.NAckReq) (*empty.Empty, error) {
	topic, err := d.topicProcessor.GetTopicByName(ctx, req.GetTopic())
	if err != nil {
		logger.Error(fmt.Sprintf("GRPC/NAck occur error [%v]. Notfound topic with name [%v]", err.Error(), req.GetTopic()))
		return nil, status.New(codes.NotFound, err.Error()).Err()
	}

	err = d.itemProcessor.Update(ctx, topic.Id, req.GetDpsAssignedUniqueId(), entity.ItemStatusReadyToDeliver, req.GetMetaData())
	if err != nil {
		logger.Error(fmt.Sprintf("GRPC/NAck occur error [%v]. Notfound topic with name [%v]", err.Error(), req.GetTopic()))
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	return &empty.Empty{}, nil
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
	t := entity.Topic{
		Name:          req.TopicName,
		Active:        uint(entity.TopicStatusActive),
		DeliverPolicy: entity.TopicDeliveryPolicy(req.DeliverPolicy).String(),
	}
	topic, err := d.topicProcessor.CreateTopic(ctx, t)
	if err != nil {
		logger.Error(fmt.Sprintf("GRPC/CreateTopic occur error [%v]", err.Error()))
		return &dps_pb.CreateTopicRes{}, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	return &dps_pb.CreateTopicRes{
		TopicId:        uint32(topic.Id),
		Name:           topic.Name,
		Active:         dps_pb.TopicActive(topic.Active),
		DeliveryPolicy: topic.DeliverPolicy,
	}, nil
}
