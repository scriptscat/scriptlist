package repository

import (
	"context"

	"github.com/scriptscat/scriptlist/internal/model/entity"
)

type IUser interface {
	Find(ctx context.Context, id int64) (*entity.User, error)
}

var defaultUser IUser

func User() IUser {
	return defaultUser
}

func RegisterUser(i IUser) {
	defaultUser = i
}
