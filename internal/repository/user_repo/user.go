package user_repo

import (
	"context"
	"fmt"
	"time"

	"github.com/scriptscat/scriptlist/internal/repository/user_profile_repo"

	"github.com/cago-frame/cago/database/cache"
	"github.com/cago-frame/cago/database/db"
	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
)

//go:generate mockgen -source=user.go -destination=mock/user.go
type UserRepo interface {
	Find(ctx context.Context, id int64) (*user_entity.User, error)
	FindByPrefix(ctx context.Context, query string) ([]*user_entity.User, error)
}

var defaultUser UserRepo

func User() UserRepo {
	return defaultUser
}

func RegisterUser(i UserRepo) {
	defaultUser = i
}

type user struct {
}

func NewUserRepo() UserRepo {
	return &user{}
}

func (u *user) userKey(id int64) string {
	return fmt.Sprintf("user:%d", id)
}

func (u *user) Find(ctx context.Context, id int64) (*user_entity.User, error) {
	ret := &user_entity.User{}
	if err := cache.Ctx(ctx).GetOrSet(u.userKey(id), func() (interface{}, error) {
		ret := &user_entity.User{}
		if err := db.Ctx(ctx).First(ret, "uid=?", id).Error; err != nil {
			if db.RecordNotFound(err) {
				// 从归档表中查找
				archive := &user_entity.UserArchive{}
				if err := db.Ctx(ctx).First(archive, "uid=?", id).Error; err != nil {
					if db.RecordNotFound(err) {
						return nil, nil
					}
					return nil, err
				} else {
					return nil, err
				}

			} else {
				return nil, err
			}
		}
		// 从profile中查找信息，后续可能会考虑独立用户系统
		profile, err := user_profile_repo.UserProfile().Find(ctx, id)
		if err != nil {
			return nil, err
		}
		if profile != nil {
			ret.ProfileAvatar = profile.Avatar
		}
		return ret, nil
	}, cache.Expiration(time.Minute)).Scan(&ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (u *user) FindByPrefix(ctx context.Context, query string) ([]*user_entity.User, error) {
	var ret []*user_entity.User
	if err := db.Ctx(ctx).Where("username LIKE ?", fmt.Sprintf("%s%%", query)).Find(&ret).Error; err != nil {
		return nil, err
	}
	if len(ret) < 5 {
		// 从归档表中查找
		var archive []*user_entity.UserArchive
		if err := db.Ctx(ctx).Where("username LIKE ?", fmt.Sprintf("%s%%", query)).Find(&archive).Error; err != nil {
			return nil, err
		}
		for _, v := range archive {
			ret = append(ret, (*user_entity.User)(v))
			if len(ret) >= 5 {
				break
			}
		}
	}
	return ret, nil
}
