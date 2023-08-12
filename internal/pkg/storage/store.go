package storage

import "context"

type Store interface {
	Ping() error
	CreateTopic(ctx context.Context, store TopicStore) (TopicStore, error)
	GetActiveTopic(ctx context.Context) ([]TopicStore, error)
	FetchItemReadyToDelivery(ctx context.Context, status string) ([]ItemStore, error)
}
