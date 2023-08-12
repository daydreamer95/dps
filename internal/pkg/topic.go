package pkg

import (
	"context"
	"dps/internal/pkg/storage"
)

type TopicStatus uint

const (
	TopicStatusInActive TopicStatus = 0
	TopicStatusActive   TopicStatus = 1
)

type Topic = storage.TopicStore

type ITopicProcessor interface {
	CreateTopic(ctx context.Context, topic Topic) (Topic, error)
	GetActiveTopic(ctx context.Context) ([]Topic, error)
}

type topicProcessor struct {
}

func NewTopicProcessor() *topicProcessor {
	return &topicProcessor{}
}

func (t *topicProcessor) CreateTopic(ctx context.Context, topic Topic) (Topic, error) {
	//Validate stuff

	topic, err := GetStore().CreateTopic(topic)
	if err != nil {
		return Topic{}, err
	}
	return topic, nil
}

func (t *topicProcessor) GetActiveTopic(ctx context.Context) ([]Topic, error) {
	topics, err := GetStore().GetActiveTopic()
	if err != nil {
		return []Topic{}, err
	}
	return topics, nil
}
