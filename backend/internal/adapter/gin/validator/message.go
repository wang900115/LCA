package validator

type MessageCreateRequest struct {
	ChannelUUID string `json:"channel_uuid" binding:"required"`
	UserUUID    string `json:"user_uuid" binding:"required"`
	Content     string `json:"content" binding:"required"`
}

type MessageDeleteRequest struct {
	MessageUUID string `json:"message_uuid" binding:"required"`
}

type MessageQueryRequest struct {
	ChannelUUID string `json:"channel_uuid" binding:"required"`
}
