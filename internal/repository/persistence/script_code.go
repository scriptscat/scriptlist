package persistence

import (
	"context"

	"github.com/codfrm/cago/database/db"
	entity "github.com/scriptscat/scriptlist/internal/model/entity/script"
	"github.com/scriptscat/scriptlist/internal/repository"
)

type scriptCode struct {
}

func NewScriptCode() repository.IScriptCode {
	return &scriptCode{}
}

func (u *scriptCode) Find(ctx context.Context, id int64) (*entity.Code, error) {
	ret := &entity.Code{ID: id}
	if err := db.Ctx(ctx).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *scriptCode) Create(ctx context.Context, scriptCode *entity.Code) error {
	return db.Ctx(ctx).Create(scriptCode).Error
}

func (u *scriptCode) Update(ctx context.Context, scriptCode *entity.Code) error {
	return db.Ctx(ctx).Updates(scriptCode).Error
}

func (u *scriptCode) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Delete(&entity.Code{ID: id}).Error
}
