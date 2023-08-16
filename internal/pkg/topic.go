package pkg

import (
	"context"
	"dps/internal/pkg/storage"
	"fmt"
	"github.com/go-playground/validator/v10"
)

type TopicStatus uint

const (
	TopicStatusInActive TopicStatus = 0
	TopicStatusActive   TopicStatus = 1
)

type TopicDeliveryPolicy int

const (
	DeliveryPolicyAtLeastOnce TopicDeliveryPolicy = 0
	DeliveryPolicyAtMostOnce  TopicDeliveryPolicy = 1
)

func (t TopicDeliveryPolicy) String() string {
	switch t {
	case DeliveryPolicyAtMostOnce:
		return "AT_MOST_ONCE"
	case DeliveryPolicyAtLeastOnce:
		return "AT_LEAST_ONCE"
	default:
		return fmt.Sprintf("%d", t)
	}
}

type Topic = storage.TopicStore

type ITopicProcessor interface {
	CreateTopic(ctx context.Context, topic Topic) (Topic, error)
	GetActiveTopic(ctx context.Context) ([]Topic, error)
	GetTopicByName(ctx context.Context, name string) (Topic, error)
}

type topicProcessor struct {
}

func NewTopicProcessor() *topicProcessor {
	return &topicProcessor{}
}

func (t *topicProcessor) CreateTopic(ctx context.Context, topic Topic) (Topic, error) {
	err := validate.Struct(topic)
	if err != nil {
		return Topic{}, err.(validator.ValidationErrors)
	}
	out, err := GetStore().CreateTopic(ctx, topic)
	if err != nil {
		return Topic{}, err
	}
	return out, nil
}

func (t *topicProcessor) GetActiveTopic(ctx context.Context) ([]Topic, error) {
	topics, err := GetStore().GetActiveTopic(ctx)
	if err != nil {
		return []Topic{}, err
	}
	return topics, nil
}

func (t *topicProcessor) GetTopicByName(ctx context.Context, name string) (Topic, error) {
	topic, err := GetStore().GetTopicByName(ctx, name)
	if err != nil {
		return Topic{}, err
	}
	return topic, nil
}
