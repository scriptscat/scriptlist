package repository

import (
	"github.com/scriptscat/scriptlist/internal/pkg/db"
	"gorm.io/gorm"
)

type follow struct {
	db *gorm.DB
}

func NewFollow() Follow {
	return &follow{
		db: db.Db,
	}
}
