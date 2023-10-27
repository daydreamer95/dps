package storage

import (
	"dps/internal/pkg/dps_util"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ItemStore struct {
	dps_util.ModelBase
	Id            string    `json:"id"`
	TopicId       uint      `json:"topic_id"`
	Priority      int32     `json:"priority"`
	Status        string    `json:"status"`
	DeliverAfter  time.Time `json:"deliver_after"`
	Payload       []byte    `json:"payload"`
	MetaData      []byte    `json:"metaData"`
	LeaseDuration int32     `json:"lease_duration"`
	LeaseAfter    time.Time `json:"-" gorm:"-"`
}

// TableName
func (i *ItemStore) TableName() string {
	return "items"
}

func (i *ItemStore) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	i.Id = uuid.NewString()
	return
}
