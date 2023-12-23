package script_repo

import (
	"context"

	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/utils/httputils"
)

//go:generate mockgen -source=./script_access.go -destination=./mock/script_access.go
type ScriptAccessRepo interface {
	Find(ctx context.Context, scriptId int64, id int64) (*script_entity.ScriptAccess, error)
	FindPage(ctx context.Context, scriptId int64, page httputils.PageRequest) ([]*script_entity.ScriptAccess, int64, error)
	Create(ctx context.Context, scriptAccess *script_entity.ScriptAccess) error
	Update(ctx context.Context, scriptAccess *script_entity.ScriptAccess) error
	Delete(ctx context.Context, id int64) error

	FindByLinkID(ctx context.Context, scriptId int64, linkId int64, accessType script_entity.AccessType) ([]*script_entity.ScriptAccess, error)
}

var defaultScriptAccess ScriptAccessRepo

func ScriptAccess() ScriptAccessRepo {
	return defaultScriptAccess
}

func RegisterScriptAccess(i ScriptAccessRepo) {
	defaultScriptAccess = i
}

type scriptAccessRepo struct {
}

func NewScriptAccess() ScriptAccessRepo {
	return &scriptAccessRepo{}
}

func (u *scriptAccessRepo) Find(ctx context.Context, scriptId int64, id int64) (*script_entity.ScriptAccess, error) {
	ret := &script_entity.ScriptAccess{}
	if err := db.Ctx(ctx).Where("id=? and script_id=? and status=?", id, scriptId, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *scriptAccessRepo) Create(ctx context.Context, scriptAccess *script_entity.ScriptAccess) error {
	return db.Ctx(ctx).Create(scriptAccess).Error
}

func (u *scriptAccessRepo) Update(ctx context.Context, scriptAccess *script_entity.ScriptAccess) error {
	return db.Ctx(ctx).Updates(scriptAccess).Error
}

func (u *scriptAccessRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Model(&script_entity.ScriptAccess{}).Where("id=?", id).Update("status", consts.DELETE).Error
}

func (u *scriptAccessRepo) FindPage(ctx context.Context, scriptId int64, page httputils.PageRequest) ([]*script_entity.ScriptAccess, int64, error) {
	var list []*script_entity.ScriptAccess
	var count int64
	find := db.Ctx(ctx).Model(&script_entity.ScriptAccess{}).
		Where("script_id=? and status=?", scriptId, consts.ACTIVE)
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := find.Order("createtime desc").Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (u *scriptAccessRepo) FindByLinkID(ctx context.Context, scriptId int64, linkId int64, accessType script_entity.AccessType) ([]*script_entity.ScriptAccess, error) {
	var list []*script_entity.ScriptAccess
	if err := db.Ctx(ctx).Where("script_id=? and link_id=? and type=? and status=?",
		scriptId, linkId, accessType, consts.ACTIVE).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
