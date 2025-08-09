package issue_entity

import (
	"context"
	"net/http"
	"strings"

	"github.com/cago-frame/cago/pkg/consts"
	"github.com/cago-frame/cago/pkg/i18n"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
)

type IssueLabel struct {
	Label       string `json:"label"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

var Label = map[string]*IssueLabel{
	"bug":      {Label: "bug", Name: "BUG", Description: "反馈一个bug", Color: "#ff0000"},
	"feature":  {Label: "feature", Name: "新功能", Description: "请求增加新功能", Color: "#a2eeef"},
	"question": {Label: "question", Name: "问题", Description: "对脚本的使用存在问题", Color: "#d876e3"},
}

type ScriptIssue struct {
	ID         int64  `gorm:"column:id;type:bigint(20);not null;primary_key"`
	ScriptID   int64  `gorm:"column:script_id;type:bigint(20);not null;index:script_id"`
	UserID     int64  `gorm:"column:user_id;type:bigint(20);not null"`
	Title      string `gorm:"column:title;type:varchar(255);not null"`
	Content    string `gorm:"column:content;type:text"`
	Labels     string `gorm:"column:labels;type:varchar(255);default:''"`
	Status     int32  `gorm:"column:status;type:tinyint(4);default:0;not null"`
	Createtime int64  `gorm:"column:createtime;type:bigint(20)"`
	Updatetime int64  `gorm:"column:updatetime;type:bigint(20)"`
}

func (s *ScriptIssue) GetLabels() []string {
	if s.Labels == "" {
		return []string{}
	}
	return strings.Split(s.Labels, ",")
}

// CheckOperate 检查是否可以操作
func (s *ScriptIssue) CheckOperate(ctx context.Context) error {
	if s == nil {
		return i18n.NewErrorWithStatus(ctx, http.StatusNotFound, code.IssueNotFound)
	}
	// 非激活状态
	if s.Status != consts.ACTIVE && s.Status != consts.AUDIT {
		return i18n.NewErrorWithStatus(ctx, http.StatusNotFound, code.IssueIsDelete)
	}
	return nil
}
