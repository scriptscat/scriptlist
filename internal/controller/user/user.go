package user

import (
	"context"

	api "github.com/scriptscat/scriptlist/internal/api/user"
	service "github.com/scriptscat/scriptlist/internal/service/user"
)

type User struct {
}

func NewUser() *User {
	return &User{}
}

// Info 用户信息
func (u *User) Info(ctx context.Context, req *api.InfoRequest) (*api.InfoResponse, error) {
	return service.User().Info(ctx, req)
}
