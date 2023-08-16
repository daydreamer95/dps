package storage

import "context"

type Store interface {
	Ping() error
	GetTopicById(ctx context.Context, id uint) (TopicStore, error)
	GetTopicByName(ctx context.Context, name string) (TopicStore, error)
	CreateTopic(ctx context.Context, store TopicStore) (TopicStore, error)
	CreateItems(ctx context.Context, items ItemStore) (ItemStore, error)
	GetActiveTopic(ctx context.Context) ([]TopicStore, error)
	FetchItemReadyToDelivery(ctx context.Context, status string) ([]ItemStore, error)
}
