package repository

import (
	"github.com/scriptscat/scriptlist/internal/service/resource/domain/entity"
)

type Resource interface {
	Save(r *entity.Resource) error
	Find(id string) (*entity.Resource, error)
}
