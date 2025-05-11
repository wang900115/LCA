package entities

type TokenClaims struct {
	UserUUID    string `json:"user_uuid"`
	ChannelUUID string `json:"channel_uuid"`
	Username    string `json:"username"`

	ExpiredAt int64 `json:"expired_at"`
}
