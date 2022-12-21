package script

import (
	"context"

	"github.com/scriptscat/scriptlist/internal/model/entity/script"
)

type IScriptDomain interface {
	Find(ctx context.Context, id int64) (*script.ScriptDomain, error)
	Create(ctx context.Context, scriptDomain *script.ScriptDomain) error
	Update(ctx context.Context, scriptDomain *script.ScriptDomain) error
	Delete(ctx context.Context, id int64) error

	FindByDomain(ctx context.Context, id int64, domain string) (*script.ScriptDomain, error)
}

var defaultScriptDomain IScriptDomain

func Domain() IScriptDomain {
	return defaultScriptDomain
}

func RegisterScriptDomain(i IScriptDomain) {
	defaultScriptDomain = i
}
