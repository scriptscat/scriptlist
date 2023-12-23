package script_entity

import (
	"context"
	"time"

	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/i18n"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
)

type AccessType int32

const (
	AccessTypeUser  AccessType = 1
	AccessTypeGroup AccessType = 2
)

type AccessRole string

const (
	AccessRoleGuest   AccessRole = "guest"
	AccessRoleManager AccessRole = "manager"
	AccessRoleOwner   AccessRole = "owner"
)

type ScriptAccess struct {
	ID         int64      `gorm:"column:id;type:bigint(20);not null;primary_key"`
	ScriptID   int64      `gorm:"column:script_id;type:bigint(20);not null;index:script_id"`
	LinkID     int64      `gorm:"column:link_id;type:bigint(20);not null"`             // 用户id或者用户组id
	Type       AccessType `gorm:"column:type;type:tinyint(4);not null"`                // 1: 用户 2: 用户组
	Role       AccessRole `gorm:"column:access_permission;type:varchar(255);not null"` // 角色 访客: guest, 管理员: manager, 拥有者: owner
	Status     int32      `gorm:"column:status;type:int(11);not null"`                 // 0: 正常 1: 禁用
	Expiretime int64      `gorm:"column:expiretime;type:bigint(20)"`
	Createtime int64      `gorm:"column:createtime;type:bigint(20);not null"`
	Updatetime int64      `gorm:"column:updatetime;type:bigint(20)"`
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
