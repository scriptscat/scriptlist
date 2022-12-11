package persistence

import (
	"context"

	"github.com/codfrm/cago/database/db"
	entity "github.com/scriptscat/scriptlist/internal/model/entity/script"
	script2 "github.com/scriptscat/scriptlist/internal/repository/script"
)

type scriptCategoryList struct {
}

func NewScriptCategoryList() script2.IScriptCategoryList {
	return &scriptCategoryList{}
}

func (s *scriptCategoryList) Find(ctx context.Context, id int64) (*entity.ScriptCategoryList, error) {
	ret := &entity.ScriptCategoryList{ID: id}
	if err := db.Ctx(ctx).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (s *scriptCategoryList) FindByName(ctx context.Context, name string) (*entity.ScriptCategoryList, error) {
	ret := &entity.ScriptCategoryList{}
	if err := db.Ctx(ctx).First(ret, "name=?", name).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (s *scriptCategoryList) Create(ctx context.Context, scriptCategoryList *entity.ScriptCategoryList) error {
	return db.Ctx(ctx).Create(scriptCategoryList).Error
}

func (s *scriptCategoryList) Update(ctx context.Context, scriptCategoryList *entity.ScriptCategoryList) error {
	return db.Ctx(ctx).Updates(scriptCategoryList).Error
}

func (s *scriptCategoryList) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Delete(&entity.ScriptCategoryList{ID: id}).Error
}
