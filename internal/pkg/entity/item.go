package entity

import (
	"context"
	"dps/internal/pkg/common"
	"dps/internal/pkg/storage"
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
	UpdateStatusAndMetadata(ctx context.Context, topicId uint, assignedUniqueId string, status string, metaData []byte) error
	Delete(ctx context.Context, topicId uint, assignedUniqueId string) error
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
	if !common.ValidateBytesSize(item.Payload, 100000) {
		return Item{}, errors.New("invalid payload bytes size. Must < 100Kb")
	}
	if !common.ValidateBytesSize(item.MetaData, 65000) {
		return Item{}, errors.New("invalid payload bytes size. Must < 65Kb")
	}

	if time.Now().After(item.DeliverAfter) {
		return Item{}, errors.New(fmt.Sprintf("delivery after must be greater than now. Delivery time: %s", item.DeliverAfter.Format(time.RFC1123)))
	}

	item.Status = ItemStatusInitialize
	item, err = GetStore().CreateItems(ctx, item)
	return item, err
}

func (i *itemProcessor) UpdateStatusAndMetadata(ctx context.Context, topicId uint, assignedUniqueId string, status string, metaData []byte) error {
	err := GetStore().UpdateItemStatusAndMetaData(ctx, topicId, status, assignedUniqueId, metaData)
	if err != nil {
		return err
	}
	return nil
}

func (i *itemProcessor) Delete(ctx context.Context, topicId uint, assignedUniqueId string) error {
	err := GetStore().DeleteItem(ctx, topicId, assignedUniqueId)
	if err != nil {
		return err
	}
	return nil
}
