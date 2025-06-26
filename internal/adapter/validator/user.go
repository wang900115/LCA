package validator

type UserCreateRequest struct {
	Username    string `json:"username" binding:"required"`
	ChannelName string `json:"channel_name" binding:"required"`
}

type UserDeleteRequest struct {
}
