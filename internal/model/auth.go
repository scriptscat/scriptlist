package model

import "time"

type AdminLevel int64

const (
	Admin          AdminLevel = 1 + iota // 管理员
	SuperModerator                       // 超级版主
	Moderator                            // 版主
)

func (a AdminLevel) IsAdmin() bool {
	return a >= Admin
}

type AuthInfo struct {
	UID           int64
	Username      string
	Email         string
	EmailVerified bool
	AdminLevel    AdminLevel
	Expiretime    time.Time
}

type LoginToken struct {
	ID         string `json:"id"` // 登录id
	UID        int64  `json:"uid"`
	Token      string `json:"token"`
	LastToken  string `json:"last_token"` // 上一次的token
	Createtime int64  `json:"createtime"`
	Updatetime int64  `json:"updatetime"`
}

// Expired 是否过期
func (l *LoginToken) Expired(t int64) bool {
	return l.Updatetime+t < time.Now().Unix()
}
