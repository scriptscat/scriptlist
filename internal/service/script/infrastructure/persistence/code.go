package persistence

import (
	"fmt"
	"time"

	request2 "github.com/scriptscat/scriptlist/internal/interfaces/api/dto/request"
	"github.com/scriptscat/scriptlist/internal/pkg/cache"
	"github.com/scriptscat/scriptlist/internal/pkg/cnt"
	"github.com/scriptscat/scriptlist/internal/pkg/errs"
	"github.com/scriptscat/scriptlist/internal/service/script/domain/entity"
	"github.com/scriptscat/scriptlist/internal/service/script/domain/repository"
	"gorm.io/gorm"
)

type code struct {
	db    *gorm.DB
	cache cache.Cache
}

func NewCode(db *gorm.DB, cache cache.Cache) repository.ScriptCode {
	return &code{db: db, cache: cache}
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
	if err := cache.NewKeyDepend(c.cache, c.dependkey(code.ScriptId)).InvalidKey(); err != nil {
		return err
	}
	return c.db.Save(code).Error
}

func (c *code) SaveDefinition(definition *entity.LibDefinition) error {
	return c.db.Save(definition).Error
}

func (c *code) FindScriptDomain(scriptId int64, domain string) (*entity.ScriptDomain, error) {
	ret := &entity.ScriptDomain{}
	if err := c.cache.GetOrSet(fmt.Sprintf("script:code:list:domain:%d:%s", scriptId, domain), ret, func() (interface{}, error) {
		if err := c.db.Where("script_id=? and domain=?", scriptId, domain).First(&ret).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return ret, nil
			}
			return nil, err
		}
		return ret, nil
	}, cache.WithTTL(time.Hour*30), cache.WithKeyDepend(c.cache, c.dependkey(scriptId))); err != nil {
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

type cacheList struct {
	List  []*entity.ScriptCode
	Total int64
}

func (c *code) List(script, status int64, page *request2.Pages) ([]*entity.ScriptCode, int64, error) {
	cacheList := &cacheList{
		List:  make([]*entity.ScriptCode, 0),
		Total: 0,
	}
	if err := c.cache.GetOrSet(fmt.Sprintf("script:code:list:%d:%d:%d:%d", script, status, page.Page(), page.Size()), cacheList, func() (interface{}, error) {
		find := c.db.Model(&entity.ScriptCode{}).Where("script_id=? and status=?", script, status).Order("createtime desc")
		if err := find.Count(&cacheList.Total).Error; err != nil {
			return nil, err
		}
		if page != request2.AllPage {
			find = find.Limit(page.Size()).Offset((page.Page() - 1) * page.Size())
		}
		if err := find.Scan(&cacheList.List).Error; err != nil {
			return nil, err
		}
		return cacheList, nil
	}, cache.WithTTL(time.Hour), cache.WithKeyDepend(c.cache, c.dependkey(script))); err != nil {
		return nil, 0, err
	}
	return cacheList.List, cacheList.Total, nil
}

func (c *code) GetLatestVersion(scriptId int64) (*entity.ScriptCode, error) {
	ret := &entity.ScriptCode{}
	if err := c.cache.GetOrSet(fmt.Sprintf("script:code:latest:%d", scriptId), ret, func() (interface{}, error) {
		codes, _, err := c.List(scriptId, cnt.ACTIVE, &request2.Pages{})
		if err != nil {
			return nil, err
		}
		if len(codes) == 0 {
			return nil, errs.ErrScriptAudit
		}
		return codes[0], nil
	}, cache.WithTTL(time.Hour), cache.WithKeyDepend(c.cache, c.dependkey(scriptId))); err != nil {
		return nil, err
	}
	return ret, nil
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
