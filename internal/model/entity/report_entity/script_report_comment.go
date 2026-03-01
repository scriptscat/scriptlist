package report_entity

import (
	"context"
	"net/http"

	"github.com/cago-frame/cago/pkg/consts"
	"github.com/cago-frame/cago/pkg/i18n"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
)

// CommentType 评论类型
type CommentType int32

const (
	CommentTypeComment CommentType = iota + 1
	CommentTypeResolve
	CommentTypeReopen
)

// ScriptReportComment 举报评论
type ScriptReportComment struct {
	ID         int64       `gorm:"column:id;type:bigint(20);not null;primary_key;autoIncrement"`
	ReportID   int64       `gorm:"column:report_id;type:bigint(20);not null;index:report_id"`
	UserID     int64       `gorm:"column:user_id;type:bigint(20);not null"`
	Content    string      `gorm:"column:content;type:text"`
	Type       CommentType `gorm:"column:type;type:tinyint(4);default:1;not null"`
	Status     int32       `gorm:"column:status;type:tinyint(4);default:1;not null"`
	Createtime int64       `gorm:"column:createtime;type:bigint(20)"`
	Updatetime int64       `gorm:"column:updatetime;type:bigint(20)"`
}

// CheckOperate 检查是否可以操作
func (c *ScriptReportComment) CheckOperate(ctx context.Context) error {
	if c == nil {
		return i18n.NewErrorWithStatus(ctx, http.StatusNotFound, code.ReportCommentNotFound)
	}
	if c.Status != consts.ACTIVE {
		return i18n.NewErrorWithStatus(ctx, http.StatusNotFound, code.ReportCommentNotFound)
	}
	return nil
}
