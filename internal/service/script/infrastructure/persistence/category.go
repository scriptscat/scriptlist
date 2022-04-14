package persistence

import (
	"fmt"
	"time"

	"github.com/scriptscat/scriptlist/internal/pkg/cache"
	"github.com/scriptscat/scriptlist/internal/service/script/domain/entity"
	"github.com/scriptscat/scriptlist/internal/service/script/domain/repository"
	"gorm.io/gorm"
)

type category struct {
	db    *gorm.DB
	cache cache.Cache
}

func NewCategory(db *gorm.DB, cache cache.Cache) repository.Category {
	return &category{
		db:    db,
		cache: cache,
	}
}

func (c *category) List() ([]*entity.ScriptCategoryList, error) {
	ret := make([]*entity.ScriptCategoryList, 0)
	if err := c.db.Model(&entity.ScriptCategoryList{}).Order("sort").Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *category) LinkCategory(script, category int64) error {
	model := &entity.ScriptCategory{}
	defer c.cache.Del(fmt.Sprintf("script:category:list:%d", script))
	if err := c.db.Where("script_id=? and category_id=?", script, category).First(model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := c.db.Model(&entity.ScriptCategoryList{ID: category}).Update("num", gorm.Expr("num+1")).Error; err != nil {
				return err
			}
			return c.db.Save(&entity.ScriptCategory{
				CategoryId: category,
				ScriptId:   script,
				Createtime: time.Now().Unix(),
				Updatetime: 0,
			}).Error
		}
		return err
	}
	return nil
}

func (c *category) Save(category *entity.ScriptCategoryList) error {
	old := &entity.ScriptCategoryList{}
	if err := c.db.Where("name=?", category.Name).First(old).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.db.Save(category).Error
		}
		return err
	}
	category.ID = old.ID
	category.Num = old.Num
	category.Createtime = old.Createtime
	category.Updatetime = old.Updatetime
	return nil
}

func (c *category) GetCategoryByScriptId(scriptId int64) ([]*entity.ScriptCategoryList, error) {
	list := make([]*entity.ScriptCategoryList, 0)
	if err := c.cache.GetOrSet(fmt.Sprintf("script:category:list:%d", scriptId), &list, func() (interface{}, error) {
		tabname := (&entity.ScriptCategoryList{}).TableName()
		categoryTbName := (&entity.ScriptCategory{}).TableName()
		if err := c.db.Model(&entity.ScriptCategory{}).
			Select(tabname+".*").
			Joins(fmt.Sprintf("left join %s on %s.id=%s.category_id", tabname, tabname, categoryTbName)).
			Where("script_id=?", scriptId).Order("sort").Scan(&list).Error; err != nil {
			return nil, err
		}
		return list, nil
	}); err != nil {
		return nil, err
	}
	return list, nil
}
