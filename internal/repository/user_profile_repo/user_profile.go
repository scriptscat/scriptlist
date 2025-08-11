package user_profile_repo

import (
	"context"
	"fmt"

	"github.com/cago-frame/cago/database/cache"

	"github.com/scriptscat/scriptlist/internal/model/entity/user_profile_entity"

	"github.com/cago-frame/cago/database/db"
	"github.com/cago-frame/cago/pkg/consts"
	"github.com/cago-frame/cago/pkg/utils/httputils"
)

type UserProfileRepo interface {
	Find(ctx context.Context, id int64) (*user_profile_entity.UserProfile, error)
	FindPage(ctx context.Context, page httputils.PageRequest) ([]*user_profile_entity.UserProfile, int64, error)
	Create(ctx context.Context, userProfile *user_profile_entity.UserProfile) error
	Update(ctx context.Context, userProfile *user_profile_entity.UserProfile) error
	Delete(ctx context.Context, id int64) error
}

var defaultUserProfile UserProfileRepo

func UserProfile() UserProfileRepo {
	return defaultUserProfile
}

func RegisterUserProfile(i UserProfileRepo) {
	defaultUserProfile = i
}

type userProfileRepo struct {
}

func NewUserProfile() UserProfileRepo {
	return &userProfileRepo{}
}

func (u *userProfileRepo) Find(ctx context.Context, id int64) (*user_profile_entity.UserProfile, error) {
	ret := &user_profile_entity.UserProfile{}
	if err := db.Ctx(ctx).Where("id=? and status=?", id, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *userProfileRepo) Create(ctx context.Context, userProfile *user_profile_entity.UserProfile) error {
	return db.Ctx(ctx).Create(userProfile).Error
}

func (u *userProfileRepo) userKey(id int64) string {
	return fmt.Sprintf("user:%d", id)
}

func (u *userProfileRepo) Update(ctx context.Context, userProfile *user_profile_entity.UserProfile) error {
	// 删除缓存
	_ = cache.Ctx(ctx).Del(u.userKey(userProfile.ID))
	return db.Ctx(ctx).UpdateColumns(userProfile).Error
}

func (u *userProfileRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Model(&user_profile_entity.UserProfile{}).Where("id=?", id).Update("status", consts.DELETE).Error
}

func (u *userProfileRepo) FindPage(ctx context.Context, page httputils.PageRequest) ([]*user_profile_entity.UserProfile, int64, error) {
	var list []*user_profile_entity.UserProfile
	var count int64
	find := db.Ctx(ctx).Model(&user_profile_entity.UserProfile{}).Where("status=?", consts.ACTIVE)
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := find.Order("createtime desc").Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}
