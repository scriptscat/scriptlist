package script_repo

import (
	"context"

	entity "github.com/scriptscat/scriptlist/internal/model/entity/script"
)

type IScriptCategoryList interface {
	Find(ctx context.Context, id int64) (*entity.ScriptCategoryList, error)
	Create(ctx context.Context, scriptCategoryList *entity.ScriptCategoryList) error
	Update(ctx context.Context, scriptCategoryList *entity.ScriptCategoryList) error
	Delete(ctx context.Context, id int64) error

	FindByName(ctx context.Context, name string) (*entity.ScriptCategoryList, error)
}

var defaultScriptCategoryList IScriptCategoryList

func ScriptCategoryList() IScriptCategoryList {
	return defaultScriptCategoryList
}

func RegisterScriptCategoryList(i IScriptCategoryList) {
	defaultScriptCategoryList = i
}
