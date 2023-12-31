package mysql

import (
	"context"
	"dps/internal/pkg/config"
	"dps/internal/pkg/dps_util"
	"dps/internal/pkg/storage"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

var conf = &config.Config

// Store implements storage.Store, and storage with db
type Store struct {
}

func (s *Store) Ping() error {
	return dbGet().ToSQLDB().Ping()
}

func (s *Store) GetTopicById(ctx context.Context, id uint) (storage.TopicStore, error) {
	var topic storage.TopicStore
	err := dbGet().WithContext(ctx).Where("id = ?", id).First(&topic).Error
	return topic, err
}

func (s *Store) GetTopicByName(ctx context.Context, name string) (storage.TopicStore, error) {
	var topic storage.TopicStore
	err := dbGet().WithContext(ctx).Where("name = ?", name).First(&topic).Error
	return topic, err
}

func (s *Store) CreateTopic(ctx context.Context, store storage.TopicStore) (storage.TopicStore, error) {
	err := dbGet().Create(&store).Error
	return store, err
}

func (s *Store) GetActiveTopic(ctx context.Context) ([]storage.TopicStore, error) {
	var storages []storage.TopicStore
	err := dbGet().WithContext(ctx).Where("").Find(&storages).Error
	return storages, err
}

func (s *Store) CreateItems(ctx context.Context, item storage.ItemStore) (storage.ItemStore, error) {
	err := dbGet().WithContext(ctx).Create(&item).Error
	return item, err
}

func (s *Store) FetchItemReadyToDelivery(ctx context.Context, status string) ([]storage.ItemStore, error) {
	var items []storage.ItemStore
	err := dbGet().WithContext(ctx).
		Where("deliver_after <= ? and status = ?", time.Now(), status).
		Find(&items).Error
	return items, err
}

func (s *Store) GetItemByStatus(ctx context.Context, status []string) ([]storage.ItemStore, error) {
	var items []storage.ItemStore
	err := dbGet().WithContext(ctx).
		Where("deliver_after <= ? and status IN (?)", time.Now(), status).
		Find(&items).Error
	return items, err
}

func (s *Store) DeleteItem(ctx context.Context, topicId uint, itemId string) error {
	res := dbGet().WithContext(ctx).
		Delete(&storage.ItemStore{
			Id:      itemId,
			TopicId: topicId,
		})
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New(fmt.Sprintf("UpdateItemStatusAndMetaData: Failed to delete item or id not found:"))
	}
	return nil
}

func (s *Store) UpdateItemStatusAndMetaData(ctx context.Context, topicId uint, status string, id string, metaData []byte) error {
	result := dbGet().WithContext(ctx).
		Model(&storage.ItemStore{}).
		Where("topic_id = ? AND id = ? AND status != ?", topicId, id, "DELIVERED").
		Updates(map[string]interface{}{"status": status, "meta_data": metaData})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New(fmt.Sprintf("UpdateItemStatusAndMetaData: Failed to update status item id [%v]:", id))
	}
	return nil
}

func (s *Store) UpdateItemsStatusByIds(ctx context.Context, ids []string, status string) error {
	err := dbGet().WithContext(ctx).
		Model(&storage.ItemStore{}).
		Where("id IN (?)", ids).
		Update("status", status).Error
	return err
}

// SetDBConn sets db conn pool
func SetDBConn(db *gorm.DB) {
	sqldb, _ := db.DB()
	sqldb.SetMaxOpenConns(int(conf.Store.MaxOpenConns))
	sqldb.SetMaxIdleConns(int(conf.Store.MaxIdleConns))
	sqldb.SetConnMaxLifetime(time.Duration(conf.Store.ConnMaxLifeTime) * time.Minute)
}

func dbGet() *dps_util.DB {
	return dps_util.DbGet(conf.Store, SetDBConn)
}
