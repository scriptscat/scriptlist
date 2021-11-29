package repository

import (
	"github.com/scriptscat/scriptlist/internal/domain/user/entity"
	"github.com/scriptscat/scriptlist/internal/http/dto/request"
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
	Find(uid, follow int64) (*entity.HomeFollow, error)
	List(uid int64, page request.Pages) ([]*entity.HomeFollow, error)
	FollowerList(uid int64, page request.Pages) ([]*entity.HomeFollow, error)
	Save(homeFollow *entity.HomeFollow) error
	Delete(uid, follow int64) error
}
