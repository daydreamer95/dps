package storage

type Store interface {
	Ping() error
	CreateTopic(store TopicStore) (TopicStore, error)
	GetActiveTopic() ([]TopicStore, error)
	FetchItemReadyToDelivery(status string) ([]ItemStore, error)
}
