package redismodel

import (
	"github.com/wang900115/LCA/internal/domain/entities"

	"github.com/golang-jwt/jwt/v5"
)

type UserTokenClaims struct {
	UserID     uint
	LastLogin  int64
	IPAddress  string
	DeviceInfo string
	jwt.RegisteredClaims
}

type ChannelTokenClaims struct {
	ChannelID uint
	Role      string
	LastJoin  int64
	jwt.RegisteredClaims
}

func (ut UserTokenClaims) ToDomain() entities.UserTokenClaims {
	return entities.UserTokenClaims{
		UserID: ut.UserID,
		LoginStatus: entities.UserLogin{
			LastLogin:  ut.LastLogin,
			IPAddress:  &ut.IPAddress,
			DeviceInfo: &ut.DeviceInfo,
		},
		ExpiredAt: ut.ExpiresAt.Unix(),
	}
}

func (ct ChannelTokenClaims) ToDomain() entities.ChannelTokenClaims {
	return entities.ChannelTokenClaims{
		ChannelID: ct.ChannelID,
		JoinStatus: entities.UserChannel{
			Role:     ct.Role,
			LastJoin: ct.LastJoin,
		},
		ExpiredAt: ct.ExpiresAt.Unix(),
	}
}
