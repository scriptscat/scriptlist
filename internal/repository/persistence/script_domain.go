package persistence

import (
	"context"

	"github.com/codfrm/cago/database/db"
	script2 "github.com/scriptscat/scriptlist/internal/model/entity/script"
	script3 "github.com/scriptscat/scriptlist/internal/repository/script"
)

type scriptDomain struct {
}

func NewScriptDomain() script3.IScriptDomain {
	return &scriptDomain{}
}

func (u *scriptDomain) Find(ctx context.Context, id int64) (*script2.ScriptDomain, error) {
	ret := &script2.ScriptDomain{ID: id}
	if err := db.Ctx(ctx).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *scriptDomain) Create(ctx context.Context, scriptDomain *script2.ScriptDomain) error {
	return db.Ctx(ctx).Create(scriptDomain).Error
}

func (u *scriptDomain) Update(ctx context.Context, scriptDomain *script2.ScriptDomain) error {
	return db.Ctx(ctx).Updates(scriptDomain).Error
}

func (u *scriptDomain) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Delete(&script2.ScriptDomain{ID: id}).Error
}

func (u *scriptDomain) FindByDomain(ctx context.Context, id int64, domain string) (*script2.ScriptDomain, error) {
	ret := &script2.ScriptDomain{}
	if err := db.Ctx(ctx).First(ret, "script_id=? and domain=?", id, domain).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}
