package script_repo

import (
	"context"

	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/utils/httputils"
)

type ScriptInviteRepo interface {
	Find(ctx context.Context, id int64) (*script_entity.ScriptInvite, error)
	FindPage(ctx context.Context, page httputils.PageRequest) ([]*script_entity.ScriptInvite, int64, error)
	Create(ctx context.Context, scriptInvite *script_entity.ScriptInvite) error
	Update(ctx context.Context, scriptInvite *script_entity.ScriptInvite) error
	Delete(ctx context.Context, id int64) error
}

var defaultScriptInvite ScriptInviteRepo

func ScriptInvite() ScriptInviteRepo {
	return defaultScriptInvite
}

func RegisterScriptInvite(i ScriptInviteRepo) {
	defaultScriptInvite = i
}

type scriptInviteRepo struct {
}

func NewScriptInvite() ScriptInviteRepo {
	return &scriptInviteRepo{}
}

func (u *scriptInviteRepo) Find(ctx context.Context, id int64) (*script_entity.ScriptInvite, error) {
	ret := &script_entity.ScriptInvite{}
	if err := db.Ctx(ctx).Where("id=? and status=?", id, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *scriptInviteRepo) Create(ctx context.Context, scriptInvite *script_entity.ScriptInvite) error {
	return db.Ctx(ctx).Create(scriptInvite).Error
}

func (u *scriptInviteRepo) Update(ctx context.Context, scriptInvite *script_entity.ScriptInvite) error {
	return db.Ctx(ctx).Updates(scriptInvite).Error
}

func (u *scriptInviteRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Model(&script_entity.ScriptInvite{}).Where("id=?", id).Update("status", consts.DELETE).Error
}

func (u *scriptInviteRepo) FindPage(ctx context.Context, page httputils.PageRequest) ([]*script_entity.ScriptInvite, int64, error) {
	var list []*script_entity.ScriptInvite
	var count int64
	find := db.Ctx(ctx).Model(&script_entity.ScriptInvite{}).Where("status=?", consts.ACTIVE)
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := find.Order("createtime desc").Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}
