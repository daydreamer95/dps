package entity

import (
	"context"
	"dps/internal/pkg/common"
	"dps/internal/pkg/storage"
	"dps/logger"
	"errors"
	"fmt"
	"time"
)

const (
	ItemStatusInitialize     = "INITIALIZE"
	ItemStatusReadyToDeliver = "READY_TO_DELIVER"
	ItemStatusDelivered      = "DELIVERED"
)

type Item = storage.ItemStore

type IItemProcessor interface {
	CreateItem(ctx context.Context, item Item) (Item, error)
}

type itemProcessor struct {
}

func NewItemProcessor() *itemProcessor {
	return &itemProcessor{}
}

func (i *itemProcessor) CreateItem(ctx context.Context, item Item) (Item, error) {
	_, err := GetStore().GetTopicById(ctx, item.TopicId)
	if err != nil {
		return Item{}, err
	}

	//Validator
	if !common.ValidateBytesSize(item.Payload, 10000) {
		return Item{}, errors.New("invalid payload bytes size. Must < 10Kb")
	}
	if !common.ValidateBytesSize(item.MetaData, 100000) {
		return Item{}, errors.New("invalid payload bytes size. Must < 100Kb")
	}

	logger.Info(fmt.Sprintf("item.DeliverAfter: %s", item.DeliverAfter))
	if time.Now().After(item.DeliverAfter) {
		return Item{}, errors.New(fmt.Sprintf("delivery after must be greater than now. Delivery time: %s", item.DeliverAfter.Format(time.RFC1123)))
	}

	item.Status = ItemStatusInitialize
	item, err = GetStore().CreateItems(ctx, item)
	return item, err
}
