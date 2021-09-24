package service

import (
	"github.com/scriptscat/scriptweb/internal/domain/user/repository"
	"github.com/scriptscat/scriptweb/internal/http/dto/respond"
	"github.com/scriptscat/scriptweb/internal/pkg/errs"
	"github.com/scriptscat/scriptweb/pkg/utils"
)

type User interface {
	UserInfo(id int64) (*respond.User, error)
	GetUserWebhook(uid int64) (string, error)
	RegenWebhook(uid int64) (string, error)
	GetUserByWebhook(token string) (int64, error)
}

type user struct {
	userRepo repository.User
}

func NewUser(userRepo repository.User) User {
	return &user{
		userRepo: userRepo,
	}
}

func (u *user) UserInfo(id int64) (*respond.User, error) {
	user, err := u.userRepo.Find(id)
	if err != nil {
		return nil, err
	}
	if (user.Groupid >= 4 && user.Groupid <= 9) || user.Groupid == 20 {
		// 禁止访问 禁止发言 等待验证会员 封禁用户组
		return respond.ToUser(user), errs.ErrUserIsBan
	}
	return respond.ToUser(user), nil
}

func (u *user) GetUserWebhook(uid int64) (string, error) {
	token, err := u.userRepo.FindUserToken(uid)
	if err != nil {
		return "", err
	}
	if token == "" {
		token = utils.RandString(64, 1)
		if err := u.userRepo.SetUserToken(uid, token); err != nil {
			return "", err
		}
	}
	return token, nil
}

func (u *user) RegenWebhook(uid int64) (string, error) {
	token := utils.RandString(64, 1)
	if err := u.userRepo.SetUserToken(uid, token); err != nil {
		return "", err
	}
	return token, nil
}

func (u *user) GetUserByWebhook(token string) (int64, error) {
	ret, err := u.userRepo.FindUserByToken(token)
	if err != nil {
		return 0, err
	}
	if ret == 0 {
		return 0, errs.ErrTokenNotFound
	}
	return ret, nil
}
