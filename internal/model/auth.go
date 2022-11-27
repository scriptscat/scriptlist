package model

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
}
