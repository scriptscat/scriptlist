package script_entity

type InviteCodeType int32

const (
	InviteCodeTypeCode InviteCodeType = 1 + iota
	InviteCodeTypeLink
)

type InviteType int32

const (
	InviteTypeAccess InviteType = 1 + iota
	InviteTypeGroup
)

type InviteStatus int32

const (
	InviteStatusUnused InviteStatus = 1 + iota
	InviteStatusUsed
	InviteStatusExpired
	InviteStatusPending
	InviteStatusReject
)

type ScriptInvite struct {
	ID           int64          `gorm:"column:id;type:bigint(20);not null;primary_key"`
	ScriptID     int64          `gorm:"column:script_id;type:bigint(20);not null;index:script_id"`
	Code         string         `gorm:"column:code;type:varchar(128);not null;index:code,unique"`
	CodeType     InviteCodeType `gorm:"column:code_type;type:tinyint(4);not null"`     // 邀请码类型 1=邀请码 2=邀请链接
	GroupID      int64          `gorm:"column:group_id;type:bigint(20)"`               // 群组id
	Type         InviteType     `gorm:"column:type;type:tinyint(4);not null"`          // 邀请类型 1=权限邀请码 2=群组邀请码
	UserID       int64          `gorm:"column:user_id;type:bigint(20)"`                // 使用用户
	IsAudit      int32          `gorm:"column:is_audit;type:tinyint(4);not null"`      // 是否需要审核
	InviteStatus InviteStatus   `gorm:"column:invite_status;type:tinyint(4);not null"` // 邀请码状态 1=未使用 2=已使用 3=已过期 4=等待 5=拒绝
	Status       int32          `gorm:"column:status;type:tinyint(4);not null"`
	Expiretime   int64          `gorm:"column:expiretime;type:bigint(20);not null"`
	Createtime   int64          `gorm:"column:createtime;type:bigint(20);not null"`
	Updatetime   int64          `gorm:"column:updatetime;type:bigint(20);not null"`
}
