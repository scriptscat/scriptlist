package script

import (
	"context"

	"github.com/codfrm/cago/database/db"
	script2 "github.com/scriptscat/scriptlist/internal/model/entity/script"
	script3 "github.com/scriptscat/scriptlist/internal/repository/script_repo"
)

type libDefinition struct {
}

func NewLibDefinition() script3.ILibDefinition {
	return &libDefinition{}
}

func (u *libDefinition) Find(ctx context.Context, id int64) (*script2.LibDefinition, error) {
	ret := &script2.LibDefinition{ID: id}
	if err := db.Ctx(ctx).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *libDefinition) Create(ctx context.Context, libDefinition *script2.LibDefinition) error {
	return db.Ctx(ctx).Create(libDefinition).Error
}

func (u *libDefinition) Update(ctx context.Context, libDefinition *script2.LibDefinition) error {
	return db.Ctx(ctx).Updates(libDefinition).Error
}

func (u *libDefinition) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Delete(&script2.LibDefinition{ID: id}).Error
}
