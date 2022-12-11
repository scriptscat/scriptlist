package script

import (
	"context"

	"github.com/scriptscat/scriptlist/internal/model/entity/script"
)

type ILibDefinition interface {
	Find(ctx context.Context, id int64) (*script.LibDefinition, error)
	Create(ctx context.Context, libDefinition *script.LibDefinition) error
	Update(ctx context.Context, libDefinition *script.LibDefinition) error
	Delete(ctx context.Context, id int64) error
}

var defaultLibDefinition ILibDefinition

func LibDefinition() ILibDefinition {
	return defaultLibDefinition
}

func RegisterLibDefinition(i ILibDefinition) {
	defaultLibDefinition = i
}
