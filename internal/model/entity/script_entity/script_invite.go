package script_entity

import (
	"context"
	"time"

	"github.com/codfrm/cago/pkg/i18n"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
)

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
	InviteStatusUnused  InviteStatus = 1 + iota // 未使用/等待接受
	InviteStatusUsed                            // 已使用
	InviteStatusExpired                         // 过期
	InviteStatusPending                         // 等待审核
	InviteStatusReject                          // 拒绝
)

type ScriptInvite struct {
	ID       int64          `gorm:"column:id;type:bigint(20);not null;primary_key"`
	ScriptID int64          `gorm:"column:script_id;type:bigint(20);not null;index:script_id;index:idx_script_id_invite_status_expiretime,priority:1;index:idx_script_id_expiretime,priority:1"`
	Code     string         `gorm:"column:code;type:varchar(128);not null;index:code,unique"`
	CodeType InviteCodeType `gorm:"column:code_type;type:tinyint(4);not null"` // 邀请码类型 1=邀请码 2=邀请链接
	GroupID  int64          `gorm:"column:group_id;type:bigint(20)"`           // 群组id
	Type     InviteType     `gorm:"column:type;type:tinyint(4);not null"`      // 邀请类型 1=权限邀请码 2=群组邀请码
	UserID   int64          `gorm:"column:user_id;type:bigint(20)"`            // 使用用户 当code_type=2 invite_type=1时 改字段为相关联的access/group member id
	IsAudit  int32          `gorm:"column:is_audit;type:tinyint(4);not null"`  // 是否需要审核 1=是 2=否
	// 等待审核->已使用 等待审核->拒绝
	// 未使用->已使用 未使用->等待审核
	InviteStatus InviteStatus `gorm:"column:invite_status;type:tinyint(4);not null;index:idx_script_id_invite_status_expiretime,priority:2;index:idx_invite_status_expiretime,priority:1"` // 邀请码状态 1=未使用 2=已使用 3=已过期 4=等待审核 5=拒绝
	Status       int32        `gorm:"column:status;type:tinyint(4);not null"`
	Expiretime   int64        `gorm:"column:expiretime;type:bigint(20);not null;index:idx_script_id_invite_status_expiretime,priority:3;index:idx_script_id_expiretime,priority:2;index:idx_invite_status_expiretime,priority:2"`
	Createtime   int64        `gorm:"column:createtime;type:bigint(20);not null"`
	Updatetime   int64        `gorm:"column:updatetime;type:bigint(20);not null"`
}

func (i *ScriptInvite) Check(ctx context.Context) error {
	if i == nil {
		return i18n.NewNotFoundError(ctx, code.AccessInviteNotFound)
	}
	if i.IsExpired() {
		return i18n.NewNotFoundError(ctx, code.AccessInviteExpired)
	}
	if !i.CanUse() {
		return i18n.NewNotFoundError(ctx, code.AccessInviteUsed)
	}
	return nil
}

func (i *ScriptInvite) GetInviteStatus() InviteStatus {
	if i.IsExpired() {
		return InviteStatusExpired
	}
	return i.InviteStatus
}

// IsExpired 是否过期
func (i *ScriptInvite) IsExpired() bool {
	return i.Expiretime > 0 && i.Expiretime < time.Now().Unix()
}

// CanUse 是否可以使用
func (i *ScriptInvite) CanUse() bool {
	return i.InviteStatus == InviteStatusUnused
}
