package validator

type ChannelCreateRequest struct {
	Name string `json:"channel_name" binding:"required"`
}

type ChannelQueryRequest struct {
}

type ChannelQueryUserRequest struct {
	Name string `json:"channel_name" binding:"required"`
}
