package script_entity

import (
	"context"
	"time"

	"github.com/cago-frame/cago/pkg/consts"
	"github.com/cago-frame/cago/pkg/i18n"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
)

type AccessType int32

const (
	AccessTypeUser AccessType = 1 + iota
	AccessTypeGroup
)

type AccessRole string

const (
	AccessRoleGuest   AccessRole = "guest"
	AccessRoleManager AccessRole = "manager"
	AccessRoleOwner   AccessRole = "owner"
)

var AccessRoleMap = map[AccessRole]int{
	AccessRoleGuest:   1,
	AccessRoleManager: 2,
	AccessRoleOwner:   3,
}

// Compare 比较权限优先级
func (a AccessRole) Compare(b AccessRole) int {
	if AccessRoleMap[a] > AccessRoleMap[b] {
		return 1
	}
	return 0
}

type AccessInviteStatus int32

const (
	AccessInviteStatusAccept AccessInviteStatus = 1 + iota
	AccessInviteStatusReject
	AccessInviteStatusPending
)

type ScriptAccess struct {
	ID           int64              `gorm:"column:id;type:bigint(20);not null;primary_key"`
	ScriptID     int64              `gorm:"column:script_id;type:bigint(20);not null;index:script_id"`
	LinkID       int64              `gorm:"column:link_id;type:bigint(20);not null"`             // 用户id或者用户组id
	Type         AccessType         `gorm:"column:type;type:tinyint(4);not null"`                // 1: 用户 2: 用户组
	Role         AccessRole         `gorm:"column:access_permission;type:varchar(255);not null"` // 角色 访客: guest, 管理员: manager, 拥有者: owner
	InviteStatus AccessInviteStatus `gorm:"column:invite_status;type:int(11);not null"`          // 1: 已接受 2: 已拒绝 3: 待接受
	Status       int32              `gorm:"column:status;type:int(11);not null"`                 // 1: 正常 2: 禁用
	Expiretime   int64              `gorm:"column:expiretime;type:bigint(20)"`
	Createtime   int64              `gorm:"column:createtime;type:bigint(20);not null"`
	Updatetime   int64              `gorm:"column:updatetime;type:bigint(20)"`
}

func (a *ScriptAccess) IsExpired() bool {
	if a.Status != consts.ACTIVE {
		return true
	}
	if a.Expiretime == 0 {
		return false
	}
	return a.Expiretime < time.Now().Unix()
}

func (a *ScriptAccess) Check(ctx context.Context) error {
	if a == nil {
		return i18n.NewNotFoundError(ctx, code.AccessNotFound)
	}
	if a.Status != consts.ACTIVE {
		return i18n.NewNotFoundError(ctx, code.AccessNotFound)
	}
	return nil
}

// IsValid 是否有效
func (m *ScriptAccess) IsValid(ctx context.Context) bool {
	if err := m.Check(ctx); err != nil {
		return false
	} else if m.InviteStatus != AccessInviteStatusAccept {
		return false
	}
	return !m.IsExpired()
}
