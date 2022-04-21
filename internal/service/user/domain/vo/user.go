package vo

import (
	"strconv"

	"github.com/scriptscat/scriptlist/internal/service/user/cnt"
	"github.com/scriptscat/scriptlist/internal/service/user/domain/entity"
)

type User struct {
	UID         int64          `json:"uid"`
	Username    string         `json:"username"`
	Avatar      string         `json:"avatar"`
	IsAdmin     cnt.AdminLevel `json:"is_admin"`
	Email       string         `json:"email,omitempty"`
	EmailStatus int64          `json:"email_status"`
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
		UID:         user.Uid,
		Username:    user.Username,
		Avatar:      avatar,
		IsAdmin:     cnt.AdminLevel(user.Adminid),
		EmailStatus: user.Emailstatus,
	}
}

func ToSelfUser(user *entity.User) *User {
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
		UID:         user.Uid,
		Username:    user.Username,
		Avatar:      avatar,
		IsAdmin:     cnt.AdminLevel(user.Adminid),
		Email:       user.Email,
		EmailStatus: user.Emailstatus,
	}
}
