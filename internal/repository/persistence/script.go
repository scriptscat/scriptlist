package persistence

import (
	"context"

	"github.com/codfrm/cago/database/db"
	entity "github.com/scriptscat/scriptlist/internal/model/entity/script"
	"github.com/scriptscat/scriptlist/internal/repository"
)

type script struct {
}

func NewScript() repository.IScript {
	return &script{}
}

func (u *script) Find(ctx context.Context, id int64) (*entity.Script, error) {
	ret := &entity.Script{ID: id}
	if err := db.Ctx(ctx).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *script) Create(ctx context.Context, script *entity.Script) error {
	return db.Ctx(ctx).Create(script).Error
}

func (u *script) Update(ctx context.Context, script *entity.Script) error {
	return db.Ctx(ctx).Updates(script).Error
}

func (u *script) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Delete(&entity.Script{ID: id}).Error
}
