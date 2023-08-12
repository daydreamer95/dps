package storage

import "dps/internal/pkg/dps_util"

type TopicStore struct {
	dps_util.ModelBase
	Id            uint   `json:"id"`
	Name          string `json:"name" validate:"required,min=3"`
	Active        uint   `json:"active"`
	DeliverPolicy string `json:"deliver_policy" validate:"oneof=AT_MOST_ONCE AT_LEAST_ONCE"`
}

// TableName
func (t *TopicStore) TableName() string {
	return "topics"
}
