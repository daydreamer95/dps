package storage

type Store interface {
	Ping() error
	GetActiveTopic() ([]TopicStore, error)
	FetchItemByTopicIds(topicId uint, status string) ([]ItemStore, error)
}
