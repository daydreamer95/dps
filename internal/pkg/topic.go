package pkg

import "time"

type Topic struct {
	Id            uint      `gorm:"column:id;type:int(6);primaryKey"`
	Name          string    `gorm:"column:name;type:varchar(3)"`
	DeliverPolicy string    `gorm:"column:deliver_policy;type:varchar(255)"`
	CreatedAt     time.Time `gorm:"column:created_at;type:timestamp"`
	UpdatedAt     time.Time `gorm:"column:updated_at;type:timestamp"`
}
