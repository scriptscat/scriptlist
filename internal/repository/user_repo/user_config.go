package user_repo

import (
	"context"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
)

type UserConfigRepo interface {
	Find(ctx context.Context, id int64) (*user_entity.UserConfig, error)
	FindPage(ctx context.Context, page httputils.PageRequest) ([]*user_entity.UserConfig, int64, error)
	Create(ctx context.Context, userConfig *user_entity.UserConfig) error
	Update(ctx context.Context, userConfig *user_entity.UserConfig) error
	Delete(ctx context.Context, id int64) error

	FindByUserID(ctx context.Context, userID int64) (*user_entity.UserConfig, error)
}

var defaultUserConfig UserConfigRepo

func UserConfig() UserConfigRepo {
	return defaultUserConfig
}

func RegisterUserConfig(i UserConfigRepo) {
	defaultUserConfig = i
}

type userConfigRepo struct {
}

func NewUserConfig() UserConfigRepo {
	return &userConfigRepo{}
}

func (u *userConfigRepo) Find(ctx context.Context, id int64) (*user_entity.UserConfig, error) {
	ret := &user_entity.UserConfig{ID: id}
	if err := db.Ctx(ctx).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *userConfigRepo) Create(ctx context.Context, userConfig *user_entity.UserConfig) error {
	return db.Ctx(ctx).Create(userConfig).Error
}

func (u *userConfigRepo) Update(ctx context.Context, userConfig *user_entity.UserConfig) error {
	return db.Ctx(ctx).Updates(userConfig).Error
}

func (u *userConfigRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Delete(&user_entity.UserConfig{ID: id}).Error
}

func (u *userConfigRepo) FindPage(ctx context.Context, page httputils.PageRequest) ([]*user_entity.UserConfig, int64, error) {
	var list []*user_entity.UserConfig
	var count int64
	if err := db.Ctx(ctx).Model(&user_entity.UserConfig{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Ctx(ctx).Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (u *userConfigRepo) FindByUserID(ctx context.Context, userID int64) (*user_entity.UserConfig, error) {
	ret := &user_entity.UserConfig{}
	if err := db.Ctx(ctx).First(ret, "uid=?", userID).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}
