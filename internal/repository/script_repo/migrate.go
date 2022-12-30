package script_repo

import (
	"context"

	"github.com/scriptscat/scriptlist/internal/model"
	entity "github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
)

type ScriptMigrateRepo interface {
	// SaveToEs 保存脚本数据到elasticsearch
	SaveToEs(ctx context.Context, s *model.ScriptSearch) error
	// List 列出脚本数据
	List(ctx context.Context, start, size int) ([]*entity.Script, error)
	// Convert 转换脚本数据
	Convert(ctx context.Context, e *entity.Script) (*model.ScriptSearch, error)
}

var defaultSearch ScriptMigrateRepo

func Migrate() ScriptMigrateRepo {
	return defaultSearch
}

func RegisterMigrate(i ScriptMigrateRepo) {
	defaultSearch = i
}

type migrateRepo struct {
}

func NewMigrateRepo() ScriptMigrateRepo {
	return &migrateRepo{}
}

func (m *migrateRepo) SaveToEs(ctx context.Context, s *model.ScriptSearch) error {
	//TODO implement me
	panic("implement me")
}

func (m *migrateRepo) List(ctx context.Context, start, size int) ([]*entity.Script, error) {
	//TODO implement me
	panic("implement me")
}

func (m *migrateRepo) Convert(ctx context.Context, e *entity.Script) (*model.ScriptSearch, error) {
	//TODO implement me
	panic("implement me")
}
