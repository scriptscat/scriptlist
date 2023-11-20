package script_repo

import (
	"context"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/utils/httputils"
)

type ScriptGroupMemberRepo interface {
	Find(ctx context.Context, id int64) (*script_entity.ScriptGroupMember, error)
	FindPage(ctx context.Context, page httputils.PageRequest) ([]*script_entity.ScriptGroupMember, int64, error)
	Create(ctx context.Context, scriptGroupMember *script_entity.ScriptGroupMember) error
	Update(ctx context.Context, scriptGroupMember *script_entity.ScriptGroupMember) error
	Delete(ctx context.Context, id int64) error
}

var defaultScriptGroupMember ScriptGroupMemberRepo

func ScriptGroupMember() ScriptGroupMemberRepo {
	return defaultScriptGroupMember
}

func RegisterScriptGroupMember(i ScriptGroupMemberRepo) {
	defaultScriptGroupMember = i
}

type scriptGroupMemberRepo struct {
}

func NewScriptGroupMember() ScriptGroupMemberRepo {
	return &scriptGroupMemberRepo{}
}

func (u *scriptGroupMemberRepo) Find(ctx context.Context, id int64) (*script_entity.ScriptGroupMember, error) {
	ret := &script_entity.ScriptGroupMember{}
	if err := db.Ctx(ctx).Where("id=? and status=?", id, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *scriptGroupMemberRepo) Create(ctx context.Context, scriptGroupMember *script_entity.ScriptGroupMember) error {
	return db.Ctx(ctx).Create(scriptGroupMember).Error
}

func (u *scriptGroupMemberRepo) Update(ctx context.Context, scriptGroupMember *script_entity.ScriptGroupMember) error {
	return db.Ctx(ctx).Updates(scriptGroupMember).Error
}

func (u *scriptGroupMemberRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Model(&script_entity.ScriptGroupMember{}).Where("id=?", id).Update("status", consts.DELETE).Error
}

func (u *scriptGroupMemberRepo) FindPage(ctx context.Context, page httputils.PageRequest) ([]*script_entity.ScriptGroupMember, int64, error) {
	var list []*script_entity.ScriptGroupMember
	var count int64
	find := db.Ctx(ctx).Model(&script_entity.ScriptGroupMember{}).Where("status=?", consts.ACTIVE)
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := find.Order("createtime desc").Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}
