package usecase

import (
	"context"

	"github.com/wang900115/LCA/internal/adapter/validator"
	"github.com/wang900115/LCA/internal/domain/entities"
	"github.com/wang900115/LCA/internal/implement"
)

type UserUsecase struct {
	userRepo  implement.UserImplement
	tokenRepo implement.TokenImplement
}

func NewUserUsecase(userRepo implement.UserImplement, tokenRepo implement.TokenImplement) *UserUsecase {
	return &UserUsecase{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

func (u *UserUsecase) CreateUser(ctx context.Context, req validator.UserCreateRequest) error {
	user := entities.User{
		Username: req.Username,
		Password: req.Password,
		NickName: req.NickName,
		FullName: req.FullName,
		LastName: req.LastName,
		Email:    req.Email,
		Phone:    req.Phone,
		Birth:    req.Birth,
	}
	return u.userRepo.Create(ctx, user)
}

func (u *UserUsecase) ReadUser(ctx context.Context, id uint) (*entities.User, error) {
	return u.userRepo.Read(ctx, id)
}

func (u *UserUsecase) UpdateUser(ctx context.Context, id uint, field string, value any) error {
	return u.userRepo.Update(ctx, id, field, value)
}

func (u *UserUsecase) DeleteUser(ctx context.Context, id uint) error {
	return u.userRepo.Delete(ctx, id)
}

func (u *UserUsecase) ParticateChannel(ctx context.Context, id uint, req validator.UserParticateRequest) error {
	userJoin := entities.UserJoin{
		Role:     "general",
		LastJoin: req.JoinTime,
	}
	return u.userRepo.CreateJoin(ctx, id, req.ChannelID, userJoin)
}

func (u *UserUsecase) JoinChannel(ctx context.Context, id uint, req validator.UserJoinRequest) (string, *entities.UserJoin, error) {
	channel, err := u.userRepo.UpdateJoinTime(ctx, id, req.ChannelID, req.JoinTime)
	if err != nil {
		return "", nil, err
	}
	tokenClaims := entities.ChannelTokenClaims{
		UserID:     id,
		ChannelID:  req.ChannelID,
		JoinStatus: channel,
	}
	token, err := u.tokenRepo.CreateChannelToken(ctx, tokenClaims)
	if err != nil {
		return "", nil, err
	}
	return token, channel, nil
}

func (u *UserUsecase) LeaveChannel(ctx context.Context, id uint, channel_id uint) error {
	return u.tokenRepo.DeleteChannelToken(ctx, id, channel_id)
}

func (u *UserUsecase) Login(ctx context.Context, req validator.UserLoginRequest) (string, *entities.UserLogin, error) {
	id, err := u.userRepo.VerifyLogin(ctx, req.Username, req.Password)
	if err != nil {
		return "", nil, err
	}
	status, err := u.userRepo.UpdateLoginTime(ctx, *id, req.Login)
	if err != nil {
		return "", nil, err
	}
	tokenClaims := entities.UserTokenClaims{
		UserID:      *id,
		LoginStatus: status,
	}
	token, err := u.tokenRepo.CreateUserToken(ctx, tokenClaims)
	if err != nil {
		return "", nil, err
	}
	return token, status, nil
}

func (u *UserUsecase) Logout(ctx context.Context, ipAddress string) error {
	id, err := u.userRepo.VerifyLogout(ctx, ipAddress)
	if err != nil {
		return err
	}
	if err := u.tokenRepo.DeleteUserToken(ctx, *id); err != nil {
		return err
	}
	return nil
}
