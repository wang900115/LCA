package irepository

import "github.com/wang900115/LCA/internal/domain/entities"

type ITokenRepository interface {
	// CreateToken
	CreateToken(entities.TokenClaims) (string, error)
	// ValidateToken with token string
	ValidateToken(string) (entities.TokenClaims, error)
	// DeleteToken with userUUID and channelUUID
	DeleteToken(string, string) error
}
