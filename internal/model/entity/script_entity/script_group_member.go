package script_entity

import (
	"context"
	"time"

	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/i18n"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
)

// ScriptGroupMember 脚本组成员
type ScriptGroupMember struct {
	ID         int64 `gorm:"column:id;type:bigint(20);not null;primary_key"`
	ScriptID   int64 `gorm:"column:script_id;type:bigint(20);not null;index:script_user"`
	GroupID    int64 `gorm:"column:group_id;type:bigint(20);not null;index:group"`
	UserID     int64 `gorm:"column:user_id;type:bigint(20);not null;index:script_user"`
	Status     int32 `gorm:"column:status;type:tinyint(4);not null"`
	Expiretime int64 `gorm:"column:expiretime;type:bigint(20);not null"`
	Createtime int64 `gorm:"column:createtime;type:bigint(20);not null"`
	Updatetime int64 `gorm:"column:updatetime;type:bigint(20);not null"`
}

func (m *ScriptGroupMember) Check(ctx context.Context) error {
	if m == nil {
		return i18n.NewNotFoundError(ctx, code.GroupMemberNotFound)
	}
	if m.Status != consts.ACTIVE {
		return i18n.NewNotFoundError(ctx, code.GroupMemberNotFound)
	}
	return nil
}

func (m *ScriptGroupMember) IsExpired() bool {
	return m.Expiretime < time.Now().Unix()
}
