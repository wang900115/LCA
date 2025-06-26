package usecase

import (
	"github.com/wang900115/LCA/internal/domain/entities"
	redisinterface "github.com/wang900115/LCA/internal/domain/interface/redis"
)

type TokenUsecase struct {
	tokenRepo redisinterface.TokenImplement
}

func NewTokenUsecase(tokenRepo redisinterface.TokenImplement) *TokenUsecase {
	return &TokenUsecase{
		tokenRepo: tokenRepo,
	}
}

func (t *TokenUsecase) CreateToken(username, channelname string) (string, error) {
	tokenClaims := entities.TokenClaims{
		User:    username,
		Channel: channelname,
	}
	return t.tokenRepo.CreateToken(tokenClaims)
}

func (t *TokenUsecase) DeleteToken(username, channelname string) error {
	return t.tokenRepo.DeleteToken(username, channelname)
}

func (t *TokenUsecase) ValidateToken(token string) (entities.TokenClaims, error) {
	return t.tokenRepo.ValidateToken(token)
}
