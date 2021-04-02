package repository

import (
	"fmt"
	"time"

	"github.com/scriptscat/scriptweb/internal/domain/script/entity"
	"github.com/scriptscat/scriptweb/internal/pkg/cache"
	"github.com/scriptscat/scriptweb/internal/pkg/db"
)

type code struct {
}

func NewCode() ScriptCode {
	return &code{}
}

func (c *code) Find(id int64) (*entity.ScriptCode, error) {
	ret := &entity.ScriptCode{}
	if err := db.Db.Find(ret, "id=?", id).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *code) Save(code *entity.ScriptCode) error {
	if err := cache.NewKeyDepend(db.Cache, c.dependkey(code.ScriptId)).InvalidKey(); err != nil {
		return err
	}
	return db.Db.Save(code).Error
}

func (c *code) dependkey(scriptId int64) string {
	return fmt.Sprintf("script:code:list:%d", scriptId)
}

func (c *code) List(script, status int64) ([]*entity.ScriptCode, error) {
	list := make([]*entity.ScriptCode, 0)
	if err := db.Cache.GetOrSet(fmt.Sprintf("script:code:list:%d:%d", script, status), &list, func() (interface{}, error) {
		find := db.Db.Model(&entity.ScriptCode{}).Where("script_id=? and status=?", script, status).Order("createtime desc")
		if err := find.Scan(&list).Error; err != nil {
			return nil, err
		}
		return list, nil
	}, cache.WithTTL(time.Minute), cache.WithKeyDepend(db.Cache, c.dependkey(script))); err != nil {
		return nil, err
	}
	return list, nil
}
