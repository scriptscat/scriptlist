package script

import (
	"context"

	entity "github.com/scriptscat/scriptlist/internal/model/entity/script"
)

type IScriptCategory interface {
	Find(ctx context.Context, id int64) (*entity.ScriptCategory, error)
	Create(ctx context.Context, scriptCategory *entity.ScriptCategory) error
	Update(ctx context.Context, scriptCategory *entity.ScriptCategory) error
	Delete(ctx context.Context, id int64) error

	LinkCategory(ctx context.Context, script, category int64) error
}

var defaultScriptCategory IScriptCategory

func ScriptCategory() IScriptCategory {
	return defaultScriptCategory
}

func RegisterScriptCategory(i IScriptCategory) {
	defaultScriptCategory = i
}
