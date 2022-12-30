package script_repo

import (
	"context"

	"github.com/codfrm/cago/database/db"
	entity "github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
)

type ScriptCodeRepo interface {
	Find(ctx context.Context, id int64) (*entity.Code, error)
	Create(ctx context.Context, scriptCode *entity.Code) error
	Update(ctx context.Context, scriptCode *entity.Code) error
	Delete(ctx context.Context, id int64) error

	FindByVersion(ctx context.Context, scriptId int64, version string) (*entity.Code, error)
	FindLatest(ctx context.Context, scriptId int64) (*entity.Code, error)
}

var defaultScriptCode ScriptCodeRepo

func ScriptCode() ScriptCodeRepo {
	return defaultScriptCode
}

func RegisterScriptCode(i ScriptCodeRepo) {
	defaultScriptCode = i
}

type scriptCodeRepo struct {
}

func NewScriptCodeRepo() ScriptCodeRepo {
	return &scriptCodeRepo{}
}

func (u *scriptCodeRepo) Find(ctx context.Context, id int64) (*entity.Code, error) {
	ret := &entity.Code{ID: id}
	if err := db.Ctx(ctx).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *scriptCodeRepo) Create(ctx context.Context, scriptCode *entity.Code) error {
	return db.Ctx(ctx).Create(scriptCode).Error
}

func (u *scriptCodeRepo) Update(ctx context.Context, scriptCode *entity.Code) error {
	return db.Ctx(ctx).Updates(scriptCode).Error
}

func (u *scriptCodeRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Delete(&entity.Code{ID: id}).Error
}

func (u *scriptCodeRepo) FindByVersion(ctx context.Context, scriptId int64, version string) (*entity.Code, error) {
	ret := &entity.Code{}
	if err := db.Ctx(ctx).First(ret, "script_id=? and version=?", scriptId, version).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *scriptCodeRepo) FindLatest(ctx context.Context, scriptId int64) (*entity.Code, error) {
	//TODO implement me
	panic("implement me")
}