package script

import (
	"context"
	"strings"

	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/scriptlist/internal/model/entity/script"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
)

type Item struct {
	ID int64 `json:"id"`
}

// ListRequest 获取脚本列表
type ListRequest struct {
	mux.Meta                     `path:"/script" method:"GET"`
	httputils.PageRequest[*Item] `                            form:",inline"`
}

type ListResponse struct {
	httputils.PageResponse[*Item] `json:",inline"`
}

// CreateRequest 创建脚本
type CreateRequest struct {
	mux.Meta    `path:"/script" method:"POST"`
	Content     string               `                             form:"content"     binding:"required,max=102400"   label:"脚本详细描述"`
	Code        string               `                             form:"code"        binding:"required,max=10485760" label:"脚本代码"`
	Name        string               `                             form:"name"        binding:"max=128"               label:"库的名字"`
	Description string               `                             form:"description" binding:"max=10240"             label:"库的描述"`
	Definition  string               `                             form:"definition"  binding:"max=10240"             label:"库的定义文件"`
	Version     string               `                             form:"version"     binding:"max=32"                label:"库的版本"`
	Type        script.Type          `                             form:"type"        binding:"required"              label:"脚本类型"` // 脚本类型：1 用户脚本 2 脚本引用库 3 订阅脚本(不支持)
	Public      script.Public        `                             form:"public"      binding:"required"              label:"公开类型"` // 公开类型：1 公开 2 半公开
	Unwell      script.UnwellContent `                             form:"unwell"      binding:"required"              label:"不适内容"` // 不适内容: 1 不适 2 适用
	Changelog   string               `                             form:"changelog"   binding:"max=102400"            label:"更新日志"`
}

func (s *CreateRequest) Validate(ctx context.Context) error {
	if s.Type == script.LibraryType {
		s.Name, s.Description = strings.TrimSpace(s.Name), strings.TrimSpace(s.Description)
		if s.Name == "" {
			return i18n.NewError(ctx, code.ScriptNameIsEmpty)
		}
		if s.Description == "" {
			return i18n.NewError(ctx, code.ScriptDescIsEmpty)
		}
	}
	return nil
}

type CreateResponse struct {
	ID int64 `json:"id"`
}
