package issue_entity

import (
	"context"

	"github.com/cago-frame/cago/pkg/consts"
	"github.com/cago-frame/cago/pkg/i18n"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
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
