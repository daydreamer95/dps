package storage

type Store interface {
	Ping() error
	GetActiveTopic() ([]TopicStore, error)
	FetchItemReadyToDelivery(status string) ([]ItemStore, error)
}
