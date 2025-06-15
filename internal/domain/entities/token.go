package entities

type TokenClaims struct {
	User    string `json:"user"`
	Channel string `json:"channel"`

	ExpiredAt int64 `json:"expired_at"`
}
