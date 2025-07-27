package validator

type CreateChannelRequest struct {
	ChannelName string `json:"channelName" binding:"required, min:1 max:30"`
	ChannelType string `json:"channelType" binding:"required"`
}

type DeleteChannelRequest struct {
	ChannelID uint `json:"channelID" binding:"required"`
}

type UpdateChannelRequest struct {
	ID          uint   `json:"id" binding:"required"`
	ChannelName string `json:"channelName" binding:"omitempty"`
	ChannelType string `json:"channelType" binding:"omitempty"`
}

type GetAllChannelsRequest struct {
}

type GetUserChannelsRequest struct {
}

type GetChannelUsersRequest struct {
	ChannelID uint `json:"channelID" binding:"required"`
}
