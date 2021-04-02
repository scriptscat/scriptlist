package service

import (
	"github.com/scriptscat/scriptweb/internal/domain/user/entity"
	"github.com/scriptscat/scriptweb/internal/domain/user/repository"
	"github.com/scriptscat/scriptweb/internal/pkg/errs"
)

type User interface {
	GetUser(id int64) (*entity.User, error)
}

type user struct {
	userRepo repository.User
}

func NewUser(userRepo repository.User) User {
	return &user{
		userRepo: userRepo,
	}
}

func (u *user) GetUser(id int64) (*entity.User, error) {
	user, err := u.userRepo.Find(id)
	if err != nil {
		return nil, err
	}
	if (user.Groupid >= 4 && user.Groupid <= 9) || user.Groupid == 20 {
		// 禁止访问 禁止发言 等待验证会员 封禁用户组
		return user, errs.ErrUserIsBan
	}
	return user, nil
}
