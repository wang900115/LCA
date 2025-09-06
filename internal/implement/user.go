package implement

import (
	"context"

	"github.com/redis/go-redis/v9"
	gormmodel "github.com/wang900115/LCA/internal/adapter/gorm/model"
	"github.com/wang900115/LCA/internal/domain/entities"

	"gorm.io/gorm"
)

type UserImplement interface {
	Create(context.Context, entities.User) error
	Read(context.Context, uint) (*entities.User, error)
	Update(context.Context, uint, string, any) error
	Delete(context.Context, uint) error
	VerifyLogin(context.Context, string, string) (*uint, error)
	VerifyLogout(context.Context, string) (*uint, error)
	CreateLogin(context.Context, entities.UserLogin) error
	UpdateLogin(context.Context, uint, int64) (*entities.UserLogin, error)
	UpdateChannel(context.Context, uint, uint, int64) (*entities.UserChannel, error)
}

type UserRepository struct {
	gorm  *gorm.DB
	redis *redis.Client
}

func NewUserRepository(gorm *gorm.DB, redis *redis.Client) UserImplement {
	return &UserRepository{
		gorm:  gorm,
		redis: redis,
	}
}

func (r *UserRepository) Create(ctx context.Context, user entities.User) error {
	userModel := gormmodel.User{
		Username: user.Username,
		Password: *user.Password, // !todo encrypt
		NickName: user.NickName,
		FullName: user.FullName,
		LastName: user.LastName,
		Email:    user.Email,
		Phone:    user.Phone,
		Birth:    user.Birth,
	}
	if err := r.gorm.WithContext(ctx).Create(&userModel).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Read(ctx context.Context, id uint) (*entities.User, error) {
	var user gormmodel.User
	if err := r.gorm.WithContext(ctx).Model(&user).First(id).Error; err != nil {
		return nil, err
	}
	return user.ToDomain(), nil
}

func (r *UserRepository) Update(ctx context.Context, id uint, field string, value any) error {
	var user gormmodel.User
	if err := r.gorm.WithContext(ctx).Where("id = ?", id).First(user).Update(field, value).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	var user gormmodel.User
	if err := r.gorm.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return err
	}
	if err := r.gorm.WithContext(ctx).Unscoped().Delete(&user).Error; err != nil {
		return err
	}
	return nil
}

// !todo(db)
func (r *UserRepository) VerifyLogin(ctx context.Context, username string, password string) (*uint, error) {
	return nil, nil
}

// !todo(db)
func (r *UserRepository) VerifyLogout(ctx context.Context, ipAddress string) (*uint, error) {
	return nil, nil
}

// !todo(redis, db)
func (r *UserRepository) CreateLogin(ctx context.Context, userLogin entities.UserLogin) error {
	return nil
}

// !todo(redis, db)
func (r *UserRepository) UpdateLogin(ctx context.Context, login_id uint, loginTime int64) (*entities.UserLogin, error) {
	return nil, nil
}

// !todo(redis, db)
func (r *UserRepository) UpdateChannel(ctx context.Context, id uint, channel_id uint, joinTime int64) (*entities.UserChannel, error) {
	return nil, nil
}
