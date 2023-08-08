package pkg

import (
	"dps/internal/pkg/storage"
)

type Item = storage.ItemStore
type Topic = storage.TopicStore

const (
	ItemStatusInitialize     = "INITIALIZE"
	ItemStatusReadyToDeliver = "READY_TO_DELIVER"
)

const (
	TopicStatusActive   = "ACTIVE"
	TopicStatusInActive = "INACTIVE"
)
