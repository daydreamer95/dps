package dto

import "dps/internal/pkg"

type PushItemRequest struct {
	TopicName string     `json:"topic_name"`
	Items     []pkg.Item `json:"items"`
}
