package validator

type UserCreateRequest struct {
	Username    string `json:"username" binding:"required"`
	ChannelUUID string `json:"channel_uuid" binding:"required"`
}

type UserDeleteRequest struct {
}
