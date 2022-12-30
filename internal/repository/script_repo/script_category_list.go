package script_repo

import (
	"context"

	"github.com/codfrm/cago/database/db"
	entity "github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
)

type ScriptCategoryListRepo interface {
	Find(ctx context.Context, id int64) (*entity.ScriptCategoryList, error)
	Create(ctx context.Context, scriptCategoryList *entity.ScriptCategoryList) error
	Update(ctx context.Context, scriptCategoryList *entity.ScriptCategoryList) error
	Delete(ctx context.Context, id int64) error

	FindByName(ctx context.Context, name string) (*entity.ScriptCategoryList, error)
}

var defaultScriptCategoryList ScriptCategoryListRepo

func ScriptCategoryList() ScriptCategoryListRepo {
	return defaultScriptCategoryList
}

func RegisterScriptCategoryList(i ScriptCategoryListRepo) {
	defaultScriptCategoryList = i
}

type scriptCategoryListRepo struct {
}

func NewScriptCategoryListRepo() ScriptCategoryListRepo {
	return &scriptCategoryListRepo{}
}

func (s *scriptCategoryListRepo) Find(ctx context.Context, id int64) (*entity.ScriptCategoryList, error) {
	ret := &entity.ScriptCategoryList{ID: id}
	if err := db.Ctx(ctx).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (s *scriptCategoryListRepo) FindByName(ctx context.Context, name string) (*entity.ScriptCategoryList, error) {
	ret := &entity.ScriptCategoryList{}
	if err := db.Ctx(ctx).First(ret, "name=?", name).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (s *scriptCategoryListRepo) Create(ctx context.Context, scriptCategoryList *entity.ScriptCategoryList) error {
	return db.Ctx(ctx).Create(scriptCategoryList).Error
}

func (s *scriptCategoryListRepo) Update(ctx context.Context, scriptCategoryList *entity.ScriptCategoryList) error {
	return db.Ctx(ctx).Updates(scriptCategoryList).Error
}

func (s *scriptCategoryListRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Delete(&entity.ScriptCategoryList{ID: id}).Error
}
