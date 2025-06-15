package usecase

import (
	"github.com/wang900115/LCA/internal/domain/entities"
	"github.com/wang900115/LCA/internal/domain/irepository"
)

type TokenUsecase struct {
	tokenRepo irepository.ITokenRepository
}

func NewTokenUsecase(tokenRepo irepository.ITokenRepository) *TokenUsecase {
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
