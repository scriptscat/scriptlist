package cnt

type AdminLevel int64

const (
	Admin          AdminLevel = 1 + iota // 管理员
	SuperModerator                       // 超级版主
	Moderator                            // 版主
)

func (a AdminLevel) IsAdmin() bool {
	return a >= Admin
}
