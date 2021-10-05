package repository

import (
	"fmt"
	"time"

	"github.com/scriptscat/scriptweb/internal/domain/script/entity"
	"github.com/scriptscat/scriptweb/internal/pkg/cache"
	"github.com/scriptscat/scriptweb/internal/pkg/db"
	"gorm.io/gorm"
)

type code struct {
	db *gorm.DB
}

func NewCode() ScriptCode {
	return &code{db: db.Db}
}

func NewTxCode(db *gorm.DB) ScriptCode {
	return &code{db: db}
}

func (c *code) Find(id int64) (*entity.ScriptCode, error) {
	ret := &entity.ScriptCode{}
	if err := c.db.Find(ret, "id=?", id).Error; err != nil {
		return nil, err
	}
	if ret.ID == 0 {
		return nil, nil
	}
	return ret, nil
}

func (c *code) Save(code *entity.ScriptCode) error {
	if err := cache.NewKeyDepend(db.Cache, c.dependkey(code.ScriptId)).InvalidKey(); err != nil {
		return err
	}
	return c.db.Save(code).Error
}

func (c *code) SaveDefinition(definition *entity.LibDefinition) error {
	return c.db.Save(definition).Error
}

func (c *code) FindScriptDomain(scriptId int64, domain string) (*entity.ScriptDomain, error) {
	ret := &entity.ScriptDomain{}
	if err := db.Cache.GetOrSet(fmt.Sprintf("script:code:list:domain:%d:%s", scriptId, domain), ret, func() (interface{}, error) {
		if err := c.db.Where("script_id=? and domain=?", scriptId, domain).First(&ret).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return ret, nil
			}
			return nil, err
		}
		return ret, nil
	}, cache.WithTTL(time.Hour*30), cache.WithKeyDepend(db.Cache, c.dependkey(scriptId))); err != nil {
		return nil, err
	}
	if ret.ID == 0 {
		return nil, nil
	}
	return ret, nil
}

func (c *code) SaveScriptDomain(domain *entity.ScriptDomain) error {
	return c.db.Save(domain).Error
}

func (c *code) dependkey(scriptId int64) string {
	return fmt.Sprintf("script:code:list:%d", scriptId)
}

func (c *code) List(script, status int64) ([]*entity.ScriptCode, error) {
	list := make([]*entity.ScriptCode, 0)
	if err := db.Cache.GetOrSet(fmt.Sprintf("script:code:list:%d:%d", script, status), &list, func() (interface{}, error) {
		find := c.db.Model(&entity.ScriptCode{}).Where("script_id=? and status=?", script, status).Order("createtime desc")
		if err := find.Scan(&list).Error; err != nil {
			return nil, err
		}
		return list, nil
	}, cache.WithTTL(time.Minute), cache.WithKeyDepend(db.Cache, c.dependkey(script))); err != nil {
		return nil, err
	}
	return list, nil
}

func (c *code) FindByVersion(scriptId int64, version string) (*entity.ScriptCode, error) {
	ret := &entity.ScriptCode{}
	if err := c.db.Find(ret, "script_id=? and version=?", scriptId, version).Error; err != nil {
		return nil, err
	}
	if ret.ID == 0 {
		return nil, nil
	}
	return ret, nil
}

func (c *code) FindDefinitionByCodeId(codeid int64) (*entity.LibDefinition, error) {
	ret := &entity.LibDefinition{}
	if err := c.db.Where("code_id=?", codeid).First(ret).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}
