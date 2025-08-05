package script_repo

import (
	"context"
	"fmt"
	"time"

	"github.com/codfrm/cago/database/cache"
	"github.com/codfrm/cago/database/db"
	entity "github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
)

type ScriptCategoryListRepo interface {
	Find(ctx context.Context, id int64) (*entity.ScriptCategoryList, error)
	Create(ctx context.Context, scriptCategoryList *entity.ScriptCategoryList) error
	Update(ctx context.Context, scriptCategoryList *entity.ScriptCategoryList) error
	Delete(ctx context.Context, id int64) error

	FindByNameAndType(ctx context.Context, name string, categoryType entity.ScriptCategoryType) (*entity.ScriptCategoryList, error)
	FindByNamePrefixAndType(ctx context.Context, namePrefix string, categoryType entity.ScriptCategoryType) ([]*entity.ScriptCategoryList, error)
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

func (s *scriptCategoryListRepo) key(id int64) string {
	return fmt.Sprintf("script:category:list:%d", id)
}

func (s *scriptCategoryListRepo) Find(ctx context.Context, id int64) (*entity.ScriptCategoryList, error) {
	ret := &entity.ScriptCategoryList{ID: id}
	if err := cache.Ctx(ctx).GetOrSet(s.key(id), func() (interface{}, error) {
		if err := db.Ctx(ctx).First(ret).Error; err != nil {
			if db.RecordNotFound(err) {
				return nil, nil
			}
			return nil, err
		}
		return ret, nil
	}, cache.Expiration(time.Hour)).Scan(&ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *scriptCategoryListRepo) FindByNameAndType(ctx context.Context, name string, categoryType entity.ScriptCategoryType) (*entity.ScriptCategoryList, error) {
	ret := &entity.ScriptCategoryList{}
	if err := db.Ctx(ctx).First(ret, "name=? and type=?", name, categoryType).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (s *scriptCategoryListRepo) FindByNamePrefixAndType(ctx context.Context, namePrefix string, categoryType entity.ScriptCategoryType) ([]*entity.ScriptCategoryList, error) {
	var ret []*entity.ScriptCategoryList
	if err := db.Ctx(ctx).Where("name LIKE ? AND type = ?", namePrefix+"%", categoryType).Find(&ret).Error; err != nil {
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
