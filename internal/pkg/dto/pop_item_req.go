package dto

type PopItemRequest struct {
	TopicName string `json:"topic_name"`
	Count     int    `json:"count"`
}
