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
	UserID    uint
	ChannelID uint
	Role      string
	LastJoin  int64
	jwt.RegisteredClaims
}

type EventTokenClaims struct {
	UserID        uint
	EventID       uint
	Role          string
	LastParticate int64
	jwt.RegisteredClaims
}

func (ut UserTokenClaims) ToDomain() *entities.UserTokenClaims {
	return &entities.UserTokenClaims{
		UserID: ut.UserID,
		LoginStatus: &entities.UserLogin{
			LastLogin:  ut.LastLogin,
			IPAddress:  &ut.IPAddress,
			DeviceInfo: &ut.DeviceInfo,
		},
		ExpiredAt: ut.ExpiresAt.Unix(),
	}
}

func (ct ChannelTokenClaims) ToDomain() *entities.ChannelTokenClaims {
	return &entities.ChannelTokenClaims{
		UserID:    ct.UserID,
		ChannelID: ct.ChannelID,
		JoinStatus: &entities.UserJoin{
			Role:     ct.Role,
			LastJoin: ct.LastJoin,
		},
		ExpiredAt: ct.ExpiresAt.Unix(),
	}
}

func (et EventTokenClaims) ToDomain() *entities.EventTokenClaims {
	return &entities.EventTokenClaims{
		UserID:  et.UserID,
		EventID: et.EventID,
		ParticateStatus: &entities.UserParticate{
			LastParticate: et.LastParticate,
			Role:          et.Role,
		},
		ExpiredAt: et.ExpiresAt.Unix(),
	}
}
