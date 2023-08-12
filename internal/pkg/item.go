package pkg

import (
	"dps/internal/pkg/storage"
)

const (
	ItemStatusInitialize     = "INITIALIZE"
	ItemStatusReadyToDeliver = "READY_TO_DELIVER"
)

type Item = storage.ItemStore

type IItemProcessor interface {
}

type itemProcessor struct {
}

func NewItemProcessor() *itemProcessor {
	return &itemProcessor{}
}
