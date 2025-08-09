package script_repo

import (
	"context"

	"github.com/cago-frame/cago/database/db"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
)

type LibDefinitionRepo interface {
	Find(ctx context.Context, id int64) (*script_entity.LibDefinition, error)
	Create(ctx context.Context, libDefinition *script_entity.LibDefinition) error
	Update(ctx context.Context, libDefinition *script_entity.LibDefinition) error
	Delete(ctx context.Context, id int64) error
}

var defaultLibDefinition LibDefinitionRepo

func LibDefinition() LibDefinitionRepo {
	return defaultLibDefinition
}

func RegisterLibDefinition(i LibDefinitionRepo) {
	defaultLibDefinition = i
}

type libDefinitionRepo struct {
}

func NewLibDefinitionRepo() LibDefinitionRepo {
	return &libDefinitionRepo{}
}

func (u *libDefinitionRepo) Find(ctx context.Context, id int64) (*script_entity.LibDefinition, error) {
	ret := &script_entity.LibDefinition{ID: id}
	if err := db.Ctx(ctx).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *libDefinitionRepo) Create(ctx context.Context, libDefinition *script_entity.LibDefinition) error {
	return db.Ctx(ctx).Create(libDefinition).Error
}

func (u *libDefinitionRepo) Update(ctx context.Context, libDefinition *script_entity.LibDefinition) error {
	return db.Ctx(ctx).Updates(libDefinition).Error
}

func (u *libDefinitionRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Delete(&script_entity.LibDefinition{ID: id}).Error
}
