package script_entity

import (
	"context"
	"github.com/codfrm/cago/pkg/consts"
	api "github.com/scriptscat/scriptlist/internal/api/script"
	"github.com/scriptscat/scriptlist/internal/repository/user_repo"
	"time"
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
	InviteStatusUnused  InviteStatus = 1 + iota // 未使用
	InviteStatusUsed                            // 已使用
	InviteStatusExpired                         // 过期
	InviteStatusPending                         // 等待审核
	InviteStatusReject                          // 拒绝
)

type ScriptInvite struct {
	ID       int64          `gorm:"column:id;type:bigint(20);not null;primary_key"`
	ScriptID int64          `gorm:"column:script_id;type:bigint(20);not null;index:script_id"`
	Code     string         `gorm:"column:code;type:varchar(128);not null;index:code,unique"`
	CodeType InviteCodeType `gorm:"column:code_type;type:tinyint(4);not null"` // 邀请码类型 1=邀请码 2=邀请链接
	GroupID  int64          `gorm:"column:group_id;type:bigint(20)"`           // 群组id
	Type     InviteType     `gorm:"column:type;type:tinyint(4);not null"`      // 邀请类型 1=权限邀请码 2=群组邀请码
	UserID   int64          `gorm:"column:user_id;type:bigint(20)"`            // 使用用户
	IsAudit  int32          `gorm:"column:is_audit;type:tinyint(4);not null"`  // 是否需要审核 1=是 2=否
	// 等待审核->已使用 等待审核->拒绝
	// 未使用->已使用 未使用->等待审核
	InviteStatus InviteStatus `gorm:"column:invite_status;type:tinyint(4);not null"` // 邀请码状态 1=未使用 2=已使用 3=已过期 4=等待审核 5=拒绝
	Status       int32        `gorm:"column:status;type:tinyint(4);not null"`
	Expiretime   int64        `gorm:"column:expiretime;type:bigint(20);not null"`
	Createtime   int64        `gorm:"column:createtime;type:bigint(20);not null"`
	Updatetime   int64        `gorm:"column:updatetime;type:bigint(20);not null"`
}

// IsExpired 是否过期
func (i *ScriptInvite) IsExpired() bool {
	return i.Expiretime > 0 && i.Expiretime < time.Now().Unix()
}

// CanUse 是否可以使用
func (i *ScriptInvite) CanUse() bool {
	return i.InviteStatus == InviteStatusUnused
}

func (i *ScriptInvite) ToInviteCode(ctx context.Context) (*api.InviteCode, error) {
	ret := &api.InviteCode{
		ID:           i.ID,
		Code:         i.Code,
		UserID:       0,
		Username:     "",
		IsAudit:      i.IsAudit == consts.YES,
		InviteStatus: i.InviteStatus,
		Expiretime:   i.Expiretime,
		Createtime:   i.Createtime,
	}
	if i.UserID > 0 {
		user, err := user_repo.User().Find(ctx, i.UserID)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, nil
		}
		ret.UserID = user.ID
		ret.Username = user.Username
	}
	if ret.InviteStatus == InviteStatusUnused && ret.Expiretime > 0 && ret.Expiretime < time.Now().Unix() {
		ret.InviteStatus = InviteStatusExpired
	}

	return ret, nil
}
