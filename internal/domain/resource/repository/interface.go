package repository

import "github.com/scriptscat/scriptweb/internal/domain/resource/entity"

type Resource interface {
	Save(r *entity.Resource) error
	Find(id string) (*entity.Resource, error)
}
