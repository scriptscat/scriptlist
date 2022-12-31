package script_repo

import (
	"context"

	"github.com/codfrm/cago/database/db"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
)

type ScriptDomainRepo interface {
	Find(ctx context.Context, id int64) (*script_entity.ScriptDomain, error)
	Create(ctx context.Context, scriptDomain *script_entity.ScriptDomain) error
	Update(ctx context.Context, scriptDomain *script_entity.ScriptDomain) error
	Delete(ctx context.Context, id int64) error

	FindByDomain(ctx context.Context, id int64, domain string) (*script_entity.ScriptDomain, error)
	List(ctx context.Context, scriptId int64) ([]*script_entity.ScriptDomain, error)
}

var defaultScriptDomain ScriptDomainRepo

func Domain() ScriptDomainRepo {
	return defaultScriptDomain
}

func RegisterScriptDomain(i ScriptDomainRepo) {
	defaultScriptDomain = i
}

type scriptDomainRepo struct {
}

func NewScriptDomainRepo() ScriptDomainRepo {
	return &scriptDomainRepo{}
}

func (u *scriptDomainRepo) Find(ctx context.Context, id int64) (*script_entity.ScriptDomain, error) {
	ret := &script_entity.ScriptDomain{ID: id}
	if err := db.Ctx(ctx).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *scriptDomainRepo) Create(ctx context.Context, scriptDomain *script_entity.ScriptDomain) error {
	return db.Ctx(ctx).Create(scriptDomain).Error
}

func (u *scriptDomainRepo) Update(ctx context.Context, scriptDomain *script_entity.ScriptDomain) error {
	return db.Ctx(ctx).Updates(scriptDomain).Error
}

func (u *scriptDomainRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Delete(&script_entity.ScriptDomain{ID: id}).Error
}

func (u *scriptDomainRepo) FindByDomain(ctx context.Context, id int64, domain string) (*script_entity.ScriptDomain, error) {
	ret := &script_entity.ScriptDomain{}
	if err := db.Ctx(ctx).First(ret, "script_id=? and domain=?", id, domain).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *scriptDomainRepo) List(ctx context.Context, scriptId int64) ([]*script_entity.ScriptDomain, error) {
	var ret []*script_entity.ScriptDomain
	if err := db.Ctx(ctx).Find(&ret, "script_id=?", scriptId).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}
