package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/scriptscat/scriptlist/internal/domain/user/entity"
	"github.com/scriptscat/scriptlist/internal/pkg/cache"
	"github.com/scriptscat/scriptlist/internal/pkg/db"
	"github.com/scriptscat/scriptlist/pkg/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User interface {
	Find(id int64) (*entity.User, error)
	FindUserToken(id int64) (string, error)
	FindUserByToken(token string) (int64, error)
	SetUserToken(id int64, token string) error
	FindUserConfig(id int64) (*entity.UserConfig, error)
	SaveUserNotifyConfig(id int64, notify datatypes.JSONMap) error
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
		if err := db.Db.First(ret, "uid=?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				archive := &entity.UserArchive{}
				if err := db.Db.First(archive, "uid=?", id).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						return nil, nil
					}
					return nil, err
				}
				ret = (*entity.User)(archive)
				return ret, nil
			}
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
		if err := db.Redis.Del(context.Background(), u.tokenKey(old)).Err(); err != nil {
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

func (u *user) FindUserConfig(id int64) (*entity.UserConfig, error) {
	ret := &entity.UserConfig{}
	if err := db.Db.First(ret, "uid=?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *user) SaveUserNotifyConfig(uid int64, notify datatypes.JSONMap) error {
	config, err := u.FindUserConfig(uid)
	if err != nil {
		return err
	}
	if config == nil {
		config = &entity.UserConfig{
			Uid:        uid,
			Notify:     notify,
			Createtime: time.Now().Unix(),
		}
	} else {
		config.Notify = notify
		config.Updatetime = time.Now().Unix()
	}
	return db.Db.Save(config).Error
}
