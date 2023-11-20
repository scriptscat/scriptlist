package script_repo

import (
	"context"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/utils/httputils"
)

type ScriptGroupRepo interface {
	Find(ctx context.Context, id int64) (*script_entity.ScriptGroup, error)
	FindPage(ctx context.Context, page httputils.PageRequest) ([]*script_entity.ScriptGroup, int64, error)
	Create(ctx context.Context, scriptGroup *script_entity.ScriptGroup) error
	Update(ctx context.Context, scriptGroup *script_entity.ScriptGroup) error
	Delete(ctx context.Context, id int64) error
}

var defaultScriptGroup ScriptGroupRepo

func ScriptGroup() ScriptGroupRepo {
	return defaultScriptGroup
}

func RegisterScriptGroup(i ScriptGroupRepo) {
	defaultScriptGroup = i
}

type scriptGroupRepo struct {
}

func NewScriptGroup() ScriptGroupRepo {
	return &scriptGroupRepo{}
}

func (u *scriptGroupRepo) Find(ctx context.Context, id int64) (*script_entity.ScriptGroup, error) {
	ret := &script_entity.ScriptGroup{}
	if err := db.Ctx(ctx).Where("id=? and status=?", id, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *scriptGroupRepo) Create(ctx context.Context, scriptGroup *script_entity.ScriptGroup) error {
	return db.Ctx(ctx).Create(scriptGroup).Error
}

func (u *scriptGroupRepo) Update(ctx context.Context, scriptGroup *script_entity.ScriptGroup) error {
	return db.Ctx(ctx).Updates(scriptGroup).Error
}

func (u *scriptGroupRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Model(&script_entity.ScriptGroup{}).Where("id=?", id).Update("status", consts.DELETE).Error
}

func (u *scriptGroupRepo) FindPage(ctx context.Context, page httputils.PageRequest) ([]*script_entity.ScriptGroup, int64, error) {
	var list []*script_entity.ScriptGroup
	var count int64
	find := db.Ctx(ctx).Model(&script_entity.ScriptGroup{}).Where("status=?", consts.ACTIVE)
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := find.Order("createtime desc").Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}
