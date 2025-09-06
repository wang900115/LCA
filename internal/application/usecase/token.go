package usecase

import (
	"github.com/wang900115/LCA/internal/domain/entities"
	"github.com/wang900115/LCA/internal/implement"
)

type TokenUsecase struct {
	tokenRepo implement.TokenImplement
}

func NewTokenUsecase(tokenRepo implement.TokenImplement) *TokenUsecase {
	return &TokenUsecase{
		tokenRepo: tokenRepo,
	}
}

func (t *TokenUsecase) UserLoginGenerateToken(id uint, status entities.UserLogin) (string, error) {
	tokenClaims := entities.UserTokenClaims{
		UserID:      id,
		LoginStatus: status,
	}
	return t.tokenRepo.CreateUserToken(tokenClaims)
}

func (t *TokenUsecase) UserJoinGenerateToken(id uint, status entities.UserChannel) (string, error) {
	channelClaims := entities.ChannelTokenClaims{
		ChannelID:  id,
		JoinStatus: status,
	}
	return t.tokenRepo.CreateChannelToken(channelClaims)
}

func (t *TokenUsecase) DeleteUserToken(user_id uint) error {
	return t.tokenRepo.DeleteUserToken(user_id)
}

func (t *TokenUsecase) DeleteChannelToken(user_id, channel_id uint) error {
	return t.tokenRepo.DeleteChannelToken(user_id, channel_id)
}

func (t *TokenUsecase) ValidateUserToken(token string) (entities.UserTokenClaims, error) {
	return t.tokenRepo.ValidateUserToken(token)
}

func (t *TokenUsecase) ValidateChannelToken(token string) (entities.ChannelTokenClaims, error) {
	return t.tokenRepo.ValidateChannelToken(token)
}
