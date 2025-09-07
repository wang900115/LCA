package validator

type ChannelCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Founder     string `json:"founder" binding:"required"`
	ChannelType string `json:"type" binding:"required"`
}

type ChannelQueryUserRequest struct {
	ChannelID uint `json:"channel" binding"required"`
}

type ChannelQueryMessageRequest struct {
	ChannelID uint `json:"channel" binding"required"`
}
