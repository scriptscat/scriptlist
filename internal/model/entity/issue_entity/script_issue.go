package issue_entity

import (
	"context"
	"net/http"
	"strings"

	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/i18n"
	"github.com/scriptscat/scriptlist/internal/model"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
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
func (s *ScriptIssue) CheckOperate(ctx context.Context, script *script_entity.Script) error {
	if err := script.CheckOperate(ctx); err != nil {
		return err
	}
	if s == nil {
		return i18n.NewErrorWithStatus(ctx, http.StatusNotFound, code.IssueNotFound)
	}
	// 脚本id不相等
	if s.ScriptID != script.ID {
		return i18n.NewErrorWithStatus(ctx, http.StatusNotFound, code.IssueNotFound)
	}
	// 非激活状态
	if s.Status != consts.ACTIVE && s.Status != consts.AUDIT {
		return i18n.NewErrorWithStatus(ctx, http.StatusNotFound, code.IssueIsDelete)
	}
	return nil
}

// CheckPermission 检查是否有权限
func (s *ScriptIssue) CheckPermission(ctx context.Context, script *script_entity.Script) error {
	if err := s.CheckOperate(ctx, script); err != nil {
		return err
	}
	uid := auth_svc.Auth().Get(ctx).UID
	// 检查uid是否是反馈者或者脚本作者
	if s.UserID != uid && script.UserID != uid && !auth_svc.Auth().Get(ctx).AdminLevel.IsAdmin(model.SuperModerator) {
		return i18n.NewErrorWithStatus(ctx, http.StatusForbidden, code.UserNotPermission)
	}
	return nil
}
