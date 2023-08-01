package storage

import "dps/internal/pkg/dps_util"

type TopicStore struct {
	dps_util.ModelBase
	Name          string `json:"name"`
	Active        string `json:"active"`
	DeliverPolicy string `json:"deliver_policy"`
}

// TableName
func (t *TopicStore) TableName() string {
	return "topics"
}