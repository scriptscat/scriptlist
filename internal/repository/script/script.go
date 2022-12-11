package script

import (
	"context"

	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/scriptscat/scriptlist/internal/model"
	entity "github.com/scriptscat/scriptlist/internal/model/entity/script"
)

type IScript interface {
	Find(ctx context.Context, id int64) (*entity.Script, error)
	Create(ctx context.Context, script *entity.Script) error
	Update(ctx context.Context, script *entity.Script) error
	Delete(ctx context.Context, id int64) error

	Search(ctx context.Context, keyword, sort string, scriptType int, page httputils.PageRequest) ([]*model.ScriptSearch, int64, error)
}

var defaultScript IScript

func Script() IScript {
	return defaultScript
}

func RegisterScript(i IScript) {
	defaultScript = i
}
