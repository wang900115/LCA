package validator

type ChannelCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Founder     string `json:"founder" binding:"required"`
	ChannelType string `json:"type" binding:"required"`
}

type ChannelQueryRequest struct {
}

type ChannelQueryUserRequest struct {
	Name string `json:"channel_name" binding:"required"`
}
