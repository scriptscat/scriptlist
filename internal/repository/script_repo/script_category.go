package script_repo

import (
	"context"
	"time"

	"github.com/codfrm/cago/database/db"
	entity "github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"gorm.io/gorm"
)

type ScriptCategoryRepo interface {
	Find(ctx context.Context, id int64) (*entity.ScriptCategory, error)
	Create(ctx context.Context, scriptCategory *entity.ScriptCategory) error
	Update(ctx context.Context, scriptCategory *entity.ScriptCategory) error
	Delete(ctx context.Context, id int64) error

	LinkCategory(ctx context.Context, script, category int64) error
}

var defaultScriptCategory ScriptCategoryRepo

func ScriptCategory() ScriptCategoryRepo {
	return defaultScriptCategory
}

func RegisterScriptCategory(i ScriptCategoryRepo) {
	defaultScriptCategory = i
}

type scriptCategoryRepo struct {
}

func NewScriptCategoryRepo() ScriptCategoryRepo {
	return &scriptCategoryRepo{}
}

func (s *scriptCategoryRepo) Find(ctx context.Context, id int64) (*entity.ScriptCategory, error) {
	ret := &entity.ScriptCategory{ID: id}
	if err := db.Ctx(ctx).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (s *scriptCategoryRepo) Create(ctx context.Context, scriptCategory *entity.ScriptCategory) error {
	return db.Ctx(ctx).Create(scriptCategory).Error
}

func (s *scriptCategoryRepo) Update(ctx context.Context, scriptCategory *entity.ScriptCategory) error {
	return db.Ctx(ctx).Updates(scriptCategory).Error
}

func (s *scriptCategoryRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Delete(&entity.ScriptCategory{ID: id}).Error
}

func (s *scriptCategoryRepo) LinkCategory(ctx context.Context, script, category int64) error {
	model := &entity.ScriptCategory{}
	if err := db.Ctx(ctx).Where("script_id=? and category_id=?", script, category).First(model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := db.Ctx(ctx).Model(&entity.ScriptCategoryList{ID: category}).Update("num", gorm.Expr("num+1")).Error; err != nil {
				return err
			}
			return db.Ctx(ctx).Save(&entity.ScriptCategory{
				CategoryID: category,
				ScriptID:   script,
				Createtime: time.Now().Unix(),
				Updatetime: 0,
			}).Error
		}
		return err
	}
	return nil
}
