package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/codfrm/cago/database/cache"
	"github.com/codfrm/cago/database/db"
	"github.com/scriptscat/scriptlist/internal/model/entity"
	"github.com/scriptscat/scriptlist/internal/repository"
)

type user struct {
}

func NewUser() repository.IUser {
	return &user{}
}

func (u *user) userkey(id int64) string {
	return fmt.Sprintf("user:%d", id)
}

func (u *user) Find(ctx context.Context, id int64) (*entity.User, error) {
	ret := &entity.User{}
	if err := cache.Ctx(ctx).GetOrSet(u.userkey(id), func() (interface{}, error) {
		ret := &entity.User{}
		if err := db.Ctx(ctx).First(ret, "uid=?", id).Error; err != nil {
			if db.RecordNotFound(err) {
				// 从归档表中查找
				archive := &entity.UserArchive{}
				if err := db.Ctx(ctx).First(archive, "uid=?", id).Error; err != nil {
					if db.RecordNotFound(err) {
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
	}, cache.Expiration(time.Minute)).Scan(ret); err != nil {
		return nil, err
	}
	return ret, nil
}
