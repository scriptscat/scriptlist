package service

import (
	"time"

	entity3 "github.com/scriptscat/scriptlist/internal/domain/user/entity"
	"github.com/scriptscat/scriptlist/internal/domain/user/repository"
	"github.com/scriptscat/scriptlist/internal/http/dto/request"
	"github.com/scriptscat/scriptlist/internal/http/dto/respond"
	"github.com/scriptscat/scriptlist/internal/pkg/errs"
	"github.com/scriptscat/scriptlist/pkg/utils"
	"gorm.io/datatypes"
)

type User interface {
	UserInfo(id int64) (*respond.User, error)
	SelfInfo(id int64) (*respond.User, error)
	GetUserWebhook(uid int64) (string, error)
	RegenWebhook(uid int64) (string, error)
	GetUserByWebhook(token string) (int64, error)
	GetUserConfig(uid int64) (*entity3.UserConfig, error)
	SetUserNotifyConfig(uid int64, notify datatypes.JSONMap) error
	IsFollow(uid, follow int64) (*entity3.HomeFollow, error)
	Follow(uid, follow int64) error
	Unfollow(uid, follow int64) error
	FollowList(uid int64, page request.Pages) ([]*entity3.HomeFollow, error)
	FollowerList(uid int64, page request.Pages) ([]*entity3.HomeFollow, error)
}

type user struct {
	userRepo   repository.User
	followRepo repository.Follow
}

func NewUser(userRepo repository.User, followRepo repository.Follow) User {
	return &user{
		userRepo:   userRepo,
		followRepo: followRepo,
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

func (u *user) SelfInfo(id int64) (*respond.User, error) {
	user, err := u.userRepo.Find(id)
	if err != nil {
		return nil, err
	}
	if (user.Groupid >= 4 && user.Groupid <= 9) || user.Groupid == 20 {
		// 禁止访问 禁止发言 等待验证会员 封禁用户组
		return respond.ToSelfUser(user), errs.ErrUserIsBan
	}
	return respond.ToSelfUser(user), nil
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

func (u *user) GetUserConfig(uid int64) (*entity3.UserConfig, error) {
	ret, err := u.userRepo.FindUserConfig(uid)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		ret = &entity3.UserConfig{
			Uid: uid,
			Notify: map[string]interface{}{
				"score": true,
			},
		}
	}
	return ret, nil
}

func (u *user) SetUserNotifyConfig(uid int64, notify datatypes.JSONMap) error {
	return u.userRepo.SaveUserNotifyConfig(uid, notify)
}

func (u *user) IsFollow(uid, follow int64) (*entity3.HomeFollow, error) {
	return u.followRepo.Find(uid, follow)
}

func (u *user) Follow(uid, follow int64) error {
	ok, err := u.followRepo.Find(uid, follow)
	if err != nil {
		return err
	}
	if ok != nil {
		return errs.NewBadRequestError(1000, "已经关注过了")
	}
	mutual, err := u.followRepo.Find(follow, uid)
	if err != nil {
		return err
	}
	user, err := u.UserInfo(uid)
	if err != nil {
		return err
	}
	fo, err := u.UserInfo(follow)
	if err != nil {
		return err
	}
	hf := &entity3.HomeFollow{
		Uid:       uid,
		Username:  user.Username,
		Followuid: fo.UID,
		Fusername: fo.Username,
		Bkname:    "",
		Status:    0,
		Dateline:  time.Now().Unix(),
	}
	if mutual != nil {
		hf.Mutual = 1
	}
	return u.followRepo.Save(hf)
}

func (u *user) Unfollow(uid, follow int64) error {
	ok, err := u.followRepo.Find(uid, follow)
	if err != nil {
		return err
	}
	if ok == nil {
		return errs.NewBadRequestError(1000, "并未关注")
	}
	mutual, err := u.followRepo.Find(follow, uid)
	if err != nil {
		return err
	}
	if mutual != nil {
		mutual.Mutual = 0
		if err := u.followRepo.Save(mutual); err != nil {
			return err
		}
	}
	return u.followRepo.Delete(uid, follow)
}

func (u *user) FollowList(uid int64, page request.Pages) ([]*entity3.HomeFollow, error) {
	return u.followRepo.List(uid, page)
}

func (u *user) FollowerList(uid int64, page request.Pages) ([]*entity3.HomeFollow, error) {
	return u.followRepo.FollowerList(uid, page)
}
