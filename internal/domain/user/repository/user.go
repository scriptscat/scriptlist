package repository

import (
	"fmt"
	"time"

	"github.com/scriptscat/scriptweb/internal/domain/user/entity"
	"github.com/scriptscat/scriptweb/internal/pkg/cache"
	"github.com/scriptscat/scriptweb/internal/pkg/db"
)

type User interface {
	Find(id int64) (*entity.User, error)
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
