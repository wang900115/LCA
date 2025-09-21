package entities

type UserTokenClaims struct {
	UserID      uint       `json:"user"`
	LoginStatus *UserLogin `json:"login"`
	ExpiredAt   int64      `json:"expired_at"`
}

type ChannelTokenClaims struct {
	UserID     uint      `json:"user"`
	ChannelID  uint      `json:"channel"`
	JoinStatus *UserJoin `json:"join"`
	ExpiredAt  int64     `json:"expired_at"`
}

type EventTokenClaims struct {
	UserID          uint           `json:"user"`
	EventID         uint           `json:"eventId"`
	ParticateStatus *UserParticate `json:"particate"`
	ExpiredAt       int64          `json:"expired_at"`
}
