package user_entity

import (
	"context"
	"net/http"
	"strconv"

	"github.com/codfrm/cago/pkg/i18n"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
)

type User struct {
	UID                int64  `gorm:"column:uid" json:"uid" form:"uid"`
	Email              string `gorm:"column:email" json:"email" form:"email"`
	Username           string `gorm:"column:username" json:"username" form:"username"`
	Password           string `gorm:"column:password" json:"password" form:"password"`
	Status             int64  `gorm:"column:status" json:"status" form:"status"`
	Emailstatus        int64  `gorm:"column:emailstatus" json:"emailstatus" form:"emailstatus"`
	Avatarstatus       int64  `gorm:"column:avatarstatus" json:"avatarstatus" form:"avatarstatus"`
	Videophotostatus   int64  `gorm:"column:videophotostatus" json:"videophotostatus" form:"videophotostatus"`
	Adminid            int64  `gorm:"column:adminid" json:"adminid" form:"adminid"`
	Groupid            int64  `gorm:"column:groupid" json:"groupid" form:"groupid"`
	Groupexpiry        int64  `gorm:"column:groupexpiry" json:"groupexpiry" form:"groupexpiry"`
	Extgroupids        string `gorm:"column:extgroupids" json:"extgroupids" form:"extgroupids"`
	Regdate            int64  `gorm:"column:regdate" json:"regdate" form:"regdate"`
	Credits            int64  `gorm:"column:credits" json:"credits" form:"credits"`
	Notifysound        int64  `gorm:"column:notifysound" json:"notifysound" form:"notifysound"`
	Timeoffset         string `gorm:"column:timeoffset" json:"timeoffset" form:"timeoffset"`
	Newpm              int64  `gorm:"column:newpm" json:"newpm" form:"newpm"`
	Newprompt          int64  `gorm:"column:newprompt" json:"newprompt" form:"newprompt"`
	Accessmasks        int64  `gorm:"column:accessmasks" json:"accessmasks" form:"accessmasks"`
	Allowadmincp       int64  `gorm:"column:allowadmincp" json:"allowadmincp" form:"allowadmincp"`
	Onlyacceptfriendpm int64  `gorm:"column:onlyacceptfriendpm" json:"onlyacceptfriendpm" form:"onlyacceptfriendpm"`
	Conisbind          int64  `gorm:"column:conisbind" json:"conisbind" form:"conisbind"`
	Freeze             int64  `gorm:"column:freeze" json:"freeze" form:"freeze"`
}

type UserArchive User

func (u *User) TableName() string {
	return "pre_common_member"
}

func (u *UserArchive) TableName() string {
	return "pre_common_member_archive"
}

func (u *User) IsBanned(ctx context.Context) error {
	if u == nil {
		return i18n.NewError(ctx, code.UserNotFound)
	}
	if (u.Groupid >= 4 && u.Groupid <= 9) || u.Groupid == 20 || u.Freeze == 1 {
		// 禁止访问 禁止发言 等待验证会员 封禁用户组
		return i18n.NewErrorWithStatus(ctx, http.StatusForbidden, code.UserIsBanned)
	}
	return nil
}

func (u *User) Avatar() string {
	return "/api/v2/users/" + strconv.FormatInt(u.UID, 10) + "/avatar"
}

func (u *User) UserInfo() UserInfo {
	if u == nil {
		return UserInfo{}
	}
	return UserInfo{
		UserID:   u.UID,
		Username: u.Username,
		Avatar:   u.Avatar(),
	}
}

type UserInfo struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	IsAdmin  int    `json:"is_admin"`
}
