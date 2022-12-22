package script_repo

import (
	"context"

	"github.com/scriptscat/scriptlist/internal/model"
	entity "github.com/scriptscat/scriptlist/internal/model/entity/script"
)

type IScriptMigrate interface {
	// SaveToEs 保存脚本数据到elasticsearch
	SaveToEs(ctx context.Context, s *model.ScriptSearch) error
	// List 列出脚本数据
	List(ctx context.Context, start, size int) ([]*entity.Script, error)
	// Convert 转换脚本数据
	Convert(ctx context.Context, e *entity.Script) (*model.ScriptSearch, error)
}

var defaultSearch IScriptMigrate

func Migrate() IScriptMigrate {
	return defaultSearch
}

func RegisterMigrate(i IScriptMigrate) {
	defaultSearch = i
}
