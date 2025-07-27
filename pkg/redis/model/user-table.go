package redismodel

import (
	"time"

	"github.com/wang900115/LCA/pkg/domain"
	rediskey "github.com/wang900115/LCA/pkg/redis/key"
)

type User struct {
	Username string

	NickName  string
	FirstName string
	LastName  string
	Birth     time.Time
	Country   string
	City      string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u User) ToHash() map[string]interface{} {
	return map[string]interface{}{
		rediskey.REDIS_FIELD_USER_USERNAME:  u.Username,
		rediskey.REDIS_FIELD_USER_NICKNAME:  u.NickName,
		rediskey.REDIS_FIELD_USER_FIRSTNAME: u.FirstName,
		rediskey.REDIS_FIELD_USER_LASTNAME:  u.LastName,
		rediskey.REDIS_FIELD_USER_BIRTH:     u.Birth.UTC().Format(time.RFC3339),

		rediskey.REDIS_FIELD_USER_CREATEDAT: u.CreatedAt.UTC().Format(time.RFC3339),
		rediskey.REDIS_FIELD_USER_UPDATEDAT: u.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func (u User) FromHash(data map[string]string) (User, error) {
	birth, err := time.Parse(time.RFC3339, data[rediskey.REDIS_FIELD_USER_BIRTH])
	if err != nil {
		return User{}, err
	}
	createdAt, err := time.Parse(time.RFC3339, data[rediskey.REDIS_FIELD_USER_CREATEDAT])
	if err != nil {
		return User{}, err
	}
	updatedAt, err := time.Parse(time.RFC3339, data[rediskey.REDIS_FIELD_USER_UPDATEDAT])
	if err != nil {
		return User{}, err
	}
	return User{
		Username:  data[rediskey.REDIS_FIELD_USER_USERNAME],
		NickName:  data[rediskey.REDIS_FIELD_USER_NICKNAME],
		FirstName: data[rediskey.REDIS_FIELD_USER_FIRSTNAME],
		LastName:  data[rediskey.REDIS_FIELD_USER_LASTNAME],

		Birth:     birth,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (u User) ToDomain(id uint) domain.User {
	return domain.User{
		ID:       id,
		Username: u.Username,

		NickName:  u.NickName,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Birth:     u.Birth,

		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u User) FromDomain(user domain.User) User {
	return User{
		Username: user.Username,

		NickName:  user.NickName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Birth:     user.Birth,
		Country:   user.Country,
		City:      user.City,

		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
