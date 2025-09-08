package validator

type ChannelCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	ChannelType string `json:"type" binding:"required"`
}
