package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/scriptscat/scriptweb/internal/domain/user/entity"
	"github.com/scriptscat/scriptweb/internal/pkg/cache"
	"github.com/scriptscat/scriptweb/internal/pkg/db"
	"github.com/scriptscat/scriptweb/pkg/utils"
	"github.com/sirupsen/logrus"
)

type User interface {
	Find(id int64) (*entity.User, error)
	FindUserToken(id int64) (string, error)
	FindUserByToken(token string) (int64, error)
	SetUserToken(id int64, token string) error
}

type user struct {
}

func NewUser() User {
	return &user{}
}

func (u *user) userkey(id int64) string {
	return fmt.Sprintf("user:%d", id)
}

func (u *user) Find(id int64) (*entity.User, error) {
	ret := &entity.User{}
	if err := db.Cache.GetOrSet(u.userkey(id), ret, func() (interface{}, error) {
		if err := db.Db.Find(ret, "uid=?", id).Error; err != nil {
			return nil, err
		}
		return ret, nil
	}, cache.WithTTL(time.Hour*24)); err != nil {
		return nil, err
	}
	return ret, nil
}

func (u *user) FindUserToken(id int64) (string, error) {
	ret, err := db.Redis.Get(context.Background(), u.tokenUserKey(id)).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return ret, nil
}

func (u *user) SetUserToken(id int64, token string) error {
	old, err := u.FindUserToken(id)
	if err != nil {
		return err
	}
	if err := db.Redis.Set(context.Background(), u.tokenUserKey(id), token, 0).Err(); err != nil {
		return err
	}
	if err := db.Redis.Set(context.Background(), u.tokenKey(token), id, 0).Err(); err != nil {
		return err
	}
	if old != "" {
		if err := db.Redis.Del(context.Background(), u.tokenKey(old)); err != nil {
			logrus.Errorf("setusertoken delete %s: %v", token, err)
		}
	}
	return err
}

func (u *user) FindUserByToken(token string) (int64, error) {
	ret, err := db.Redis.Get(context.Background(), u.tokenKey(token)).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return utils.StringToInt64(ret), nil
}

func (u *user) tokenUserKey(id int64) string {
	return fmt.Sprintf("user:token:user:%d", id)
}

func (u *user) tokenKey(t string) string {
	return fmt.Sprintf("user:token:token:%s", t)
}
