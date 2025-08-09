package script_repo

import (
	"context"

	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"

	"github.com/cago-frame/cago/database/db"
	"github.com/cago-frame/cago/pkg/consts"
	"github.com/cago-frame/cago/pkg/utils/httputils"
)

type ScriptInviteRepo interface {
	Find(ctx context.Context, scriptId int64, id int64) (*script_entity.ScriptInvite, error)
	FindByCode(ctx context.Context, code string) (*script_entity.ScriptInvite, error)
	FindAccessPage(ctx context.Context, scriptId int64, page httputils.PageRequest) ([]*script_entity.ScriptInvite, int64, error)
	FindGroupPage(ctx context.Context, scriptId int64, groupId int64, page httputils.PageRequest) ([]*script_entity.ScriptInvite, int64, error)
	Create(ctx context.Context, scriptInvite *script_entity.ScriptInvite) error
	Update(ctx context.Context, scriptInvite *script_entity.ScriptInvite) error
	Delete(ctx context.Context, scriptId int64, id int64) error
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

func (u *scriptInviteRepo) Find(ctx context.Context, scriptId int64, id int64) (*script_entity.ScriptInvite, error) {
	ret := &script_entity.ScriptInvite{}
	if err := db.Ctx(ctx).Where("id=? and script_id=? and status=?", id, scriptId, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *scriptInviteRepo) FindByCode(ctx context.Context, code string) (*script_entity.ScriptInvite, error) {
	ret := &script_entity.ScriptInvite{}
	if err := db.Ctx(ctx).Where("code=? and status=?", code, consts.ACTIVE).First(ret).Error; err != nil {
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

func (u *scriptInviteRepo) Delete(ctx context.Context, scriptId int64, id int64) error {
	return db.Ctx(ctx).Model(&script_entity.ScriptInvite{}).Where("id=? and script_id=?", id, scriptId).
		Update("status", consts.DELETE).Error
}

func (u *scriptInviteRepo) FindAccessPage(ctx context.Context, scriptId int64, page httputils.PageRequest) ([]*script_entity.ScriptInvite, int64, error) {
	var list []*script_entity.ScriptInvite
	var count int64
	find := db.Ctx(ctx).Model(&script_entity.ScriptInvite{}).
		Where("script_id=? and code_type=? and type=? and status=?", scriptId,
			script_entity.InviteCodeTypeCode, script_entity.InviteTypeAccess, consts.ACTIVE)
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	orderQuery := u.getInviteListOrder(page.GetSort(), page.GetOrder())
	if err := find.Order(orderQuery).Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}
func (u *scriptInviteRepo) getInviteListOrder(sort string, order string) string {
	var sortQuery string
	switch sort {
	case "invite_status", "expiretime":
		sortQuery = sort
	default:
		sortQuery = "createtime"
	}
	var orderQuery string
	switch order {
	case "ascend":
		orderQuery = "asc"
	case "descend":
		fallthrough
	default:
		orderQuery = "desc"
	}
	result := sortQuery + " " + orderQuery
	if result == "createtime desc" {
		result = "invite_status,expiretime"
	}
	return result
}
func (u *scriptInviteRepo) FindGroupPage(ctx context.Context, scriptId int64, groupId int64, page httputils.PageRequest) ([]*script_entity.ScriptInvite, int64, error) {
	var list []*script_entity.ScriptInvite
	var count int64
	find := db.Ctx(ctx).Model(&script_entity.ScriptInvite{}).
		Where("script_id=? and group_id=? and code_type=? and type=? and status=?", scriptId, groupId,
			script_entity.InviteCodeTypeCode, script_entity.InviteTypeGroup, consts.ACTIVE)
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	orderQuery := u.getInviteListOrder(page.GetSort(), page.GetOrder())
	if err := find.Order(orderQuery).Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}
