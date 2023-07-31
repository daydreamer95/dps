package storage

import "dps/internal/pkg/dps_util"

type TopicStore struct {
	dps_util.ModelBase
	Name          string `json:"name"`
	DeliverPolicy string `json:"deliver_policy"`
}

// TableName
func (t *TopicStore) TableName() string {
	return "topics"
}
