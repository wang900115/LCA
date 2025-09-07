package validator

type MessageReadRequest struct {
	MessageId uint `json:"id" binding:"required"`
}
