package script_entity

type ScriptWatchLevel int

const (
	ScriptWatchLevelNone ScriptWatchLevel = iota
	ScriptWatchLevelVersion
	ScriptWatchLevelIssue
	ScriptWatchLevelIssueComment
)

type ScriptWatch struct {
	ID         int64            `gorm:"column:id;type:bigint(20);not null;primary_key"`
	UserID     int64            `gorm:"column:user_id;type:bigint(20);not null;index:script_user,unique" json:"user_id"`
	ScriptID   int64            `gorm:"column:script_id;type:bigint(20);not null;index:script_user,unique;index:script" json:"script_id"`
	Level      ScriptWatchLevel `gorm:"column:level;type:int(11);not null" json:"level"`
	Createtime int64            `gorm:"column:createtime;type:bigint(20);not null"`
	Updatetime int64            `gorm:"column:updatetime;type:bigint(20)"`
}
