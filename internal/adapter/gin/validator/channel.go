package validator

type ChannelCreateRequest struct {
}

type ChannelQueryRequest struct {
}

type ChannelQueryUserRequest struct {
	ChannelUUID string `json:"channel_uuid" binding:"required"`
}
