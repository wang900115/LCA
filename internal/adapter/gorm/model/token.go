package model

import (
	"LCA/internal/domain/entities"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserUUID    string `json:"user_uuid"`
	ChannelUUID string `json:"channel_uuid"`
	Username    string `json:"username"`

	jwt.RegisteredClaims
}

func (t TokenClaims) ToDomain() entities.TokenClaims {
	return entities.TokenClaims{
		UserUUID:    t.UserUUID,
		ChannelUUID: t.ChannelUUID,
		Username:    t.Username,
		ExpiredAt:   t.ExpiresAt.Unix(),
	}
}
