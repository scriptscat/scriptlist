package script_repo

import (
	"context"
	"strconv"
	"time"

	"github.com/cago-frame/cago/database/cache"
	cache2 "github.com/cago-frame/cago/database/cache/cache"
	"github.com/cago-frame/cago/database/db"
	entity "github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"gorm.io/gorm"
)

type ScriptCategoryRepo interface {
	Find(ctx context.Context, id int64) (*entity.ScriptCategory, error)
	Create(ctx context.Context, scriptCategory *entity.ScriptCategory) error
	Delete(ctx context.Context, scriptCategory *entity.ScriptCategory) error

	FindByScriptId(ctx context.Context, scriptId int64, categoryType entity.ScriptCategoryType) ([]*entity.ScriptCategory, error)
	DeleteByScriptId(ctx context.Context, id int64) error
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
	if err := db.Ctx(ctx).Create(scriptCategory).Error; err != nil {
		return err
	}
	if err := db.Ctx(ctx).Model(&entity.ScriptCategoryList{ID: scriptCategory.CategoryID}).
		Update("num", gorm.Expr("num+1")).Error; err != nil {
		return err
	}
	return s.keyDepend(ctx, scriptCategory.ScriptID).InvalidKey(ctx)
}

func (s *scriptCategoryRepo) Delete(ctx context.Context, scriptCategory *entity.ScriptCategory) error {
	if err := db.Ctx(ctx).Delete(&entity.ScriptCategory{ID: scriptCategory.ID}).Error; err != nil {
		return err
	}
	if err := db.Ctx(ctx).Model(&entity.ScriptCategoryList{ID: scriptCategory.CategoryID}).
		Update("num", gorm.Expr("num-1")).Error; err != nil {
		return err
	}
	return s.keyDepend(ctx, scriptCategory.ScriptID).InvalidKey(ctx)
}

func (s *scriptCategoryRepo) key(script int64) string {
	return "script:category:" + strconv.FormatInt(script, 10)
}

func (s *scriptCategoryRepo) FindByScriptId(ctx context.Context, script int64, scriptCategory entity.ScriptCategoryType) ([]*entity.ScriptCategory, error) {
	var ret []*entity.ScriptCategory
	if err := cache.Ctx(ctx).GetOrSet(s.key(script)+":"+strconv.Itoa(int(scriptCategory)), func() (interface{}, error) {
		scriptCategoryTable := db.Default().NamingStrategy.TableName("script_category")
		scriptCategoryListTable := db.Default().NamingStrategy.TableName("script_category_list")
		if err := db.Ctx(ctx).
			Joins("LEFT JOIN "+scriptCategoryListTable+" ON "+scriptCategoryTable+".category_id = "+scriptCategoryListTable+".id").
			Where(scriptCategoryTable+".script_id = ? AND "+scriptCategoryListTable+".type = ?", script, scriptCategory).
			Find(&ret).Error; err != nil {
			return nil, err
		}
		return ret, nil
	}, cache2.Expiration(time.Hour), cache2.WithKeyDepend(cache.Default(), s.key(script)+":dep")).Scan(&ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *scriptCategoryRepo) DeleteByScriptId(ctx context.Context, scriptId int64) error {
	if err := db.Ctx(ctx).Delete(&entity.ScriptCategory{}, "script_id=?", scriptId).Error; err != nil {
		return err
	}
	return s.keyDepend(ctx, scriptId).InvalidKey(ctx)
}

func (s *scriptCategoryRepo) keyDepend(ctx context.Context, script int64) *cache2.KeyDepend {
	return cache2.NewKeyDepend(cache.Default(), s.key(script)+":dep")
}
