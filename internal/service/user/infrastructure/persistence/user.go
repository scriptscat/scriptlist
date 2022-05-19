package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/scriptscat/scriptlist/internal/pkg/cache"
	"github.com/scriptscat/scriptlist/internal/pkg/errs"
	"github.com/scriptscat/scriptlist/internal/service/user/domain/entity"
	"github.com/scriptscat/scriptlist/internal/service/user/domain/repository"
	"github.com/scriptscat/scriptlist/pkg/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type user struct {
	db    *gorm.DB
	redis *redis.Client
	cache cache.Cache
}

func NewUser(db *gorm.DB, redis *redis.Client, cache cache.Cache) repository.User {
	return &user{
		db:    db,
		redis: redis,
		cache: cache,
	}
}

func (u *user) userkey(id int64) string {
	return fmt.Sprintf("user:%d", id)
}

func (u *user) Find(id int64) (*entity.User, error) {
	ret := &entity.User{}
	if err := u.cache.GetOrSet(u.userkey(id), ret, func() (interface{}, error) {
		if err := u.db.First(ret, "uid=?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				archive := &entity.UserArchive{}
				if err := u.db.First(archive, "uid=?", id).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						return nil, errs.ErrUserNotFound
					}
					return nil, err
				}
				ret = (*entity.User)(archive)
				return ret, nil
			}
			return nil, err
		}
		return ret, nil
	}, cache.WithTTL(time.Minute)); err != nil {
		return nil, err
	}
	return ret, nil
}

func (u *user) FindUserToken(id int64) (string, error) {
	ret, err := u.redis.Get(context.Background(), u.tokenUserKey(id)).Result()
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
	if err := u.redis.Set(context.Background(), u.tokenUserKey(id), token, 0).Err(); err != nil {
		return err
	}
	if err := u.redis.Set(context.Background(), u.tokenKey(token), id, 0).Err(); err != nil {
		return err
	}
	if old != "" {
		if err := u.redis.Del(context.Background(), u.tokenKey(old)).Err(); err != nil {
			logrus.Errorf("setusertoken delete %s: %v", token, err)
		}
	}
	return err
}

func (u *user) FindUserByToken(token string) (int64, error) {
	ret, err := u.redis.Get(context.Background(), u.tokenKey(token)).Result()
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
	if err := u.db.First(ret, "uid=?", id).Error; err != nil {
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
	return u.db.Save(config).Error
}

func (u *user) FindByUsername(username string) (*entity.User, error) {
	ret := &entity.User{}
	if err := u.db.First(ret, "username=?", username).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}
