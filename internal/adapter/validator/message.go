package validator

type MessageCreateRequest struct {
	Content string `json:"content" binding:"required"`
}

type MessageDeleteRequest struct {
	MessageUUID string `json:"message_uuid" binding:"required"`
}

type MessageQueryRequest struct {
}
