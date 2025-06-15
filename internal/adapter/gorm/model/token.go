package model

import (
	"github.com/wang900115/LCA/internal/domain/entities"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	User    string `json:"user"`
	Channel string `json:"channel"`

	jwt.RegisteredClaims
}

func (t TokenClaims) ToDomain() entities.TokenClaims {
	return entities.TokenClaims{
		User:      t.User,
		Channel:   t.Channel,
		ExpiredAt: t.ExpiresAt.Unix(),
	}
}
