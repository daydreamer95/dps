package dto

type PopItemRequest struct {
	requestedItems []ItemRequest `json:"requested_item"`
}

type ItemRequest struct {
	TopicName string `json:"topic_name"`
	Count     int    `json:"count"`
}
