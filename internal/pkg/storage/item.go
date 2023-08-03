package storage

import (
	"dps/internal/pkg/dps_util"
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
}

// TableName
func (t *ItemStore) TableName() string {
	return "items"
}
