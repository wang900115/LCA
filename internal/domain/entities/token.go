package entities

type UserTokenClaims struct {
	UserID      uint      `json:"user"`
	LoginStatus UserLogin `json:"login"`
	ExpiredAt   int64     `json:"expired_at"`
}

type ChannelTokenClaims struct {
	ChannelID  uint        `json:"channel"`
	JoinStatus UserChannel `json:"join"`
	ExpiredAt  int64       `json:"expired_at"`
}
