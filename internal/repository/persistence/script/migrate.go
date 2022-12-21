package persistence

import (
	"context"

	"github.com/scriptscat/scriptlist/internal/model"
	entity "github.com/scriptscat/scriptlist/internal/model/entity/script"
	script3 "github.com/scriptscat/scriptlist/internal/repository/script"
)

type migrate struct {
}

func NewMigrate() script3.IScriptMigrate {
	return &migrate{}
}

func (m *migrate) SaveToEs(ctx context.Context, s *model.ScriptSearch) error {
	//TODO implement me
	panic("implement me")
}

func (m *migrate) List(ctx context.Context, start, size int) ([]*entity.Script, error) {
	//TODO implement me
	panic("implement me")
}

func (m *migrate) Convert(ctx context.Context, e *entity.Script) (*model.ScriptSearch, error) {
	//TODO implement me
	panic("implement me")
}
