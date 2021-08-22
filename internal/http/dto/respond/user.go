package respond

import (
	"strconv"

	"github.com/scriptscat/scriptweb/internal/domain/user/entity"
)

type User struct {
	UID      int64  `json:"uid"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	IsAdmin  int64  `json:"is_admin"`
}

func ToUser(user *entity.User) *User {
	if user == nil {
		return &User{
			UID:      0,
			Username: "封禁用户",
		}
	}
	avatar := ""
	if user.Avatarstatus == 1 {
		avatar = "/api/v1/user/avatar/" + strconv.FormatInt(user.Uid, 10)
	}
	return &User{
		UID:      user.Uid,
		Username: user.Username,
		Avatar:   avatar,
		IsAdmin:  user.Adminid,
	}
}
