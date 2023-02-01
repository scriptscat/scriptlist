package script_repo

import (
	"context"
	"strconv"
	"time"

	"github.com/codfrm/cago/database/cache"
	cache2 "github.com/codfrm/cago/database/cache/cache"
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
	List(ctx context.Context, script int64) ([]*entity.ScriptCategory, error)
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

func (s *scriptCategoryRepo) key(script int64) string {
	return "script:category:" + strconv.FormatInt(script, 10)
}

func (s *scriptCategoryRepo) LinkCategory(ctx context.Context, script, category int64) error {
	model := &entity.ScriptCategory{}
	if err := db.Ctx(ctx).Where("script_id=? and category_id=?", script, category).First(model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := db.Ctx(ctx).Model(&entity.ScriptCategoryList{ID: category}).Update("num", gorm.Expr("num+1")).Error; err != nil {
				return err
			}
			if err := db.Ctx(ctx).Save(&entity.ScriptCategory{
				CategoryID: category,
				ScriptID:   script,
				Createtime: time.Now().Unix(),
				Updatetime: 0,
			}).Error; err != nil {
				return err
			}
			return cache2.NewKeyDepend(cache.Default(), s.key(script)+":dep").InvalidKey(ctx)
		}
		return err
	}
	return nil
}

func (s *scriptCategoryRepo) List(ctx context.Context, script int64) ([]*entity.ScriptCategory, error) {
	var ret []*entity.ScriptCategory
	if err := cache.Ctx(ctx).GetOrSet(s.key(script), func() (interface{}, error) {
		if err := db.Ctx(ctx).Find(&ret, "script_id=?", script).Error; err != nil {
			return nil, err
		}
		return ret, nil
	}, cache2.Expiration(time.Hour), cache2.WithKeyDepend(cache.Default(), s.key(script))).Scan(&ret); err != nil {
		return nil, err
	}
	return ret, nil
}
