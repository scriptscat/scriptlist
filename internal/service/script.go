package service

import (
	"gorm.io/gorm"
)

type Script interface {
	GetScript(id int64)
	GetScriptMeta(id int64)
}

type script struct {
	db *gorm.DB
}

func NewScript(db *gorm.DB) Script {
	return &script{
		db: db,
	}
}

func (s *script) GetScript(id int64) {

}

func (s *script) GetScriptMeta(id int64) {

}
