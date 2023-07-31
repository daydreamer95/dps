package pkg

import "time"

type Item struct {
	Id            string    `gorm:"column:id;primaryKey;type:char(255)"`
	TopicId       uint      `gorm:"column:topic_id;type:int(6)"`
	Priority      int32     `gorm:"column:priority;type:int"`
	DeliverAfter  time.Time `gorm:"column:deliver_after;type:timestamp"`
	Payload       []byte    `gorm:"column:pay_load;type:blob"`
	MetaData      []byte    `gorm:"column:meta_data;type:blob"`
	LeaseDuration int32     `gorm:"column:lease_duration;type:int"`
	CreatedAt     time.Time `gorm:"column:created_at;type:timestamp"`
	UpdatedAt     time.Time `gorm:"column:updated_at;type:timestamp"`
}
