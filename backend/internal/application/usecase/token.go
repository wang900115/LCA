package usecase

import (
	"LCA/internal/domain/entities"
	"LCA/internal/domain/irepository"
)

type TokenUsecase struct {
	tokenRepo irepository.ITokenRepository
}

func NewTokenUsecase(tokenRepo irepository.ITokenRepository) *TokenUsecase {
	return &TokenUsecase{
		tokenRepo: tokenRepo,
	}
}

func (t *TokenUsecase) CreateToken(userUUID, channelUUID, username string) (string, error) {
	tokenClaims := entities.TokenClaims{
		UserUUID:    userUUID,
		ChannelUUID: channelUUID,
		Username:    username,
	}
	return t.tokenRepo.CreateToken(tokenClaims)
}

func (t *TokenUsecase) DeleteToken(userUUID, channelUUID string) error {
	return t.tokenRepo.DeleteToken(userUUID, channelUUID)
}

func (t *TokenUsecase) ValidateToken(token string) (entities.TokenClaims, error) {
	return t.tokenRepo.ValidateToken(token)
}
