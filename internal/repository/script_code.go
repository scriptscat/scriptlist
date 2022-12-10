package repository

import (
	"context"

	entity "github.com/scriptscat/scriptlist/internal/model/entity/script"
)

type IScriptCode interface {
	Find(ctx context.Context, id int64) (*entity.Code, error)
	Create(ctx context.Context, scriptCode *entity.Code) error
	Update(ctx context.Context, scriptCode *entity.Code) error
	Delete(ctx context.Context, id int64) error
}

var defaultScriptCode IScriptCode

func ScriptCode() IScriptCode {
	return defaultScriptCode
}

func RegisterScriptCode(i IScriptCode) {
	defaultScriptCode = i
}
