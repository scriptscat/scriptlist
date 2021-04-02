package repository

import (
	"time"

	"github.com/scriptscat/scriptweb/internal/domain/script/entity"
	"github.com/scriptscat/scriptweb/internal/pkg/db"
	"gorm.io/gorm"
)

type category struct {
}

func NewCategory() Category {
	return &category{}
}

func (c *category) List() ([]*entity.ScriptCategoryList, error) {
	ret := make([]*entity.ScriptCategoryList, 0)
	if err := db.Db.Model(&entity.ScriptCategoryList{}).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *category) LinkCategory(script, category int64) error {
	model := &entity.ScriptCategory{}
	if err := db.Db.Where("script_id=? and category_id=?", script, category).First(model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := db.Db.Model(&entity.ScriptCategoryList{ID: category}).Update("num", gorm.Expr("num+1")).Error; err != nil {
				return err
			}
			return db.Db.Save(&entity.ScriptCategory{
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
	if err := db.Db.Where("name=?", category.Name).First(old).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return db.Db.Save(category).Error
		}
		return err
	}
	category.ID = old.ID
	category.Num = old.Num
	category.Createtime = old.Createtime
	category.Updatetime = old.Updatetime
	return nil
}
