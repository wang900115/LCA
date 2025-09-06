package usecase

import (
	"context"

	"github.com/wang900115/LCA/internal/adapter/validator"
	"github.com/wang900115/LCA/internal/domain/entities"
	"github.com/wang900115/LCA/internal/implement"
)

type UserUsecase struct {
	userRepo implement.UserImplement
}

func NewUserUsecase(userRepo implement.UserImplement) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
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

func (u *UserUsecase) JoinChannel(ctx context.Context, id uint, channel_id uint, joinTime int64) (*entities.UserChannel, error) {
	return u.userRepo.UpdateChannel(ctx, id, channel_id, joinTime)
}

func (u *UserUsecase) Login(ctx context.Context, req validator.UserLoginRequest) (*uint, *entities.UserLogin, error) {
	id, err := u.userRepo.VerifyLogin(ctx, req.Username, req.Password)
	if err != nil {
		return nil, nil, err
	}
	status, err := u.userRepo.UpdateLogin(ctx, *id, req.Login)
	if err != nil {
		return nil, nil, err
	}
	return id, status, nil
}

func (u *UserUsecase) Logout(ctx context.Context, ipAddress string) (*uint, error) {
	id, err := u.userRepo.VerifyLogout(ctx, ipAddress)
	if err != nil {
		return nil, err
	}
	return id, nil
}
