package user_repo

import (
	"context"

	"github.com/cago-frame/cago/database/db"
	"github.com/cago-frame/cago/pkg/utils/httputils"
	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
)

type FollowRepo interface {
	Find(ctx context.Context, uid, follow int64) (*user_entity.HomeFollow, error)
	List(ctx context.Context, uid int64, page httputils.PageRequest) ([]*user_entity.HomeFollow, int64, error)
	// FollowerList 关注我的人
	FollowerList(ctx context.Context, uid int64, page httputils.PageRequest) ([]*user_entity.HomeFollow, int64, error)
	Save(ctx context.Context, homeFollow *user_entity.HomeFollow) error
	Delete(ctx context.Context, uid, follow int64) error
	// UpdateMutual 更新互相关注状态
	UpdateMutual(ctx context.Context, uid, follow, mutual int64) error
}

var defaultFollow FollowRepo

func Follow() FollowRepo {
	return defaultFollow
}

func RegisterFollow(i FollowRepo) {
	defaultFollow = i
}

type follow struct {
}

func NewFollowRepo() FollowRepo {
	return &follow{}
}

func (f *follow) Find(ctx context.Context, uid, follow int64) (*user_entity.HomeFollow, error) {
	ret := &user_entity.HomeFollow{}
	if err := db.Ctx(ctx).First(ret, "uid=? and followuid=?", uid, follow).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (f *follow) List(ctx context.Context, uid int64, page httputils.PageRequest) ([]*user_entity.HomeFollow, int64, error) {
	list := make([]*user_entity.HomeFollow, 0)
	find := db.Ctx(ctx).Model(&user_entity.HomeFollow{}).Where("uid=?", uid)
	var count int64
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	find = find.Limit(page.GetLimit()).Offset(page.GetLimit())
	if err := find.Scan(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (f *follow) FollowerList(ctx context.Context, uid int64, page httputils.PageRequest) ([]*user_entity.HomeFollow, int64, error) {
	list := make([]*user_entity.HomeFollow, 0)
	find := db.Ctx(ctx).Model(&user_entity.HomeFollow{}).Where("followuid=?", uid)
	var count int64
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	find = find.Limit(page.GetLimit()).Offset(page.GetOffset())
	if err := find.Scan(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (f *follow) Save(ctx context.Context, homeFollow *user_entity.HomeFollow) error {
	return db.Ctx(ctx).Create(homeFollow).Error
}

func (f *follow) Delete(ctx context.Context, uid, follow int64) error {
	return db.Ctx(ctx).Delete(&user_entity.HomeFollow{}, "uid=? and followuid=?", uid, follow).Error
}

func (f *follow) UpdateMutual(ctx context.Context, uid, follow, mutual int64) error {
	return db.Ctx(ctx).Model(&user_entity.HomeFollow{}).Where("uid=? and followuid=?", uid, follow).Update("mutual", mutual).Error
}
