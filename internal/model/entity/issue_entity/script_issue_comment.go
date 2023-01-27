package issue_entity

import (
	"context"
	"net/http"

	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/i18n"
	"github.com/scriptscat/scriptlist/internal/model"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
)

const (
	CommentTypeComment CommentType = iota + 1
	CommentTypeChangeTitle
	CommentTypeChangeLabel
	CommentTypeOpen
	CommentTypeClose
	CommentTypeDelete
)

type CommentType int

type ScriptIssueComment struct {
	ID         int64       `gorm:"column:id;type:bigint(20);not null;primary_key"`
	IssueID    int64       `gorm:"column:issue_id;type:bigint(20);not null;index:issue_id"`
	UserID     int64       `gorm:"column:user_id;type:bigint(20);not null"`
	Content    string      `gorm:"column:content;type:text;not null"`
	Type       CommentType `gorm:"column:type;type:tinyint(4);default:0;not null"`
	Status     int32       `gorm:"column:status;type:tinyint(4);default:0;not null"`
	Createtime int64       `gorm:"column:createtime;type:bigint(20)"`
	Updatetime int64       `gorm:"column:updatetime;type:bigint(20)"`
}

func (c *ScriptIssueComment) CheckOperate(ctx context.Context) error {
	if c == nil {
		return i18n.NewError(ctx, code.IssueCommentNotFound)
	}
	if c.Status != consts.ACTIVE {
		return i18n.NewError(ctx, code.IssueCommentNotFound)
	}
	return nil
}

// CheckPermission 检查是否有权限操作
func (c *ScriptIssueComment) CheckPermission(ctx context.Context, script *script_entity.Script, issue *ScriptIssue) error {
	if err := c.CheckOperate(ctx); err != nil {
		return err
	}
	if err := issue.CheckOperate(ctx, script); err != nil {
		return err
	}
	uid := auth_svc.Auth().Get(ctx).UID
	// 检查uid是否是反馈者或者脚本作者
	if c.UserID != uid && script.UserID != uid && auth_svc.Auth().Get(ctx).AdminLevel.IsAdmin(model.SuperModerator) {
		return i18n.NewErrorWithStatus(ctx, http.StatusForbidden, code.UserNotPermission)
	}
	return nil
}
