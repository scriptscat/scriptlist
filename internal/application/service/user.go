package service

import (
	"github.com/scriptscat/scriptweb/internal/domain/user/service"
	"github.com/scriptscat/scriptweb/internal/interfaces/dto/respond"
)

type User interface {
	UserInfo(uid int64) (*respond.User, error)
}

type user struct {
	userSvc service.User
}

func NewUser(userSvc service.User) User {
	return &user{
		userSvc: userSvc,
	}
}

func (u *user) UserInfo(uid int64) (*respond.User, error) {
	user, err := u.userSvc.GetUser(uid)
	if err != nil {
		return nil, err
	}
	return respond.ToUser(user), nil
}
