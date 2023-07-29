package dto

import (
	"dps/internal/pkg/entity"
)

type PushItemRequest struct {
	TopicName string        `json:"topic_name"`
	Items     []entity.Item `json:"items"`
}
