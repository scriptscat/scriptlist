package repository

import (
	"github.com/scriptscat/scriptlist/internal/domain/user/entity"
	"gorm.io/datatypes"
)

type User interface {
	Find(id int64) (*entity.User, error)
	FindUserToken(id int64) (string, error)
	FindUserByToken(token string) (int64, error)
	SetUserToken(id int64, token string) error
	FindUserConfig(id int64) (*entity.UserConfig, error)
	SaveUserNotifyConfig(id int64, notify datatypes.JSONMap) error
}

type Follow interface {
}
