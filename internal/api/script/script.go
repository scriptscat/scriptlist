package script

import (
	"context"
	"strings"

	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
)

type Script struct {
	Script               *Code `json:"script"`
	ID                   int64 `json:"id"`
	user_entity.UserInfo `json:",inline"`
	PostID               int64           `json:"post_id"`
	Name                 string          `json:"name"`
	Description          string          `json:"description"`
	Category             []*CategoryList `json:"category"`
	Status               int64           `json:"status"`
	Score                int64           `json:"score"`
	ScoreNum             int64           `json:"score_num"`
	Type                 int             `json:"type"`
	Public               int             `json:"public"`
	Unwell               int             `json:"unwell"`
	Archive              int             `json:"archive"`
	TodayInstall         int64           `json:"today_install"`
	TotalInstall         int64           `json:"total_install"`
	Createtime           int64           `json:"createtime"`
	Updatetime           int64           `json:"updatetime"`
}

// CategoryList 拥有的分类列表
type CategoryList struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	// 本分类下脚本数量
	Num        int64 `json:"num"`
	Sort       int32 `json:"sort"`
	Createtime int64 `json:"createtime"`
	Updatetime int64 `json:"updatetime"`
}

type Code struct {
	ID                   int64 `json:"id" form:"id"`
	user_entity.UserInfo `json:",inline"`
	Meta                 string      `json:"meta,omitempty"`
	MetaJson             interface{} `json:"meta_json"`
	ScriptID             int64       `json:"script_id"`
	Version              string      `json:"version"`
	Changelog            string      `json:"changelog"`
	Status               int64       `json:"status"`
	Createtime           int64       `json:"createtime"`
	Code                 string      `json:"code,omitempty"`
	Definition           string      `json:"definition,omitempty"`
}

// ListRequest 获取脚本列表
type ListRequest struct {
	mux.Meta              `path:"/scripts" method:"GET"`
	httputils.PageRequest `form:",inline"`
	Keyword               string `form:"keyword"`
	ScriptType            int    `form:"script_type,default=0" binding:"oneof=0 1 2 3 4"` // 0:全部 1: 脚本 2: 库 3: 后台脚本 4: 定时脚本
	Sort                  string `form:"sort,default=today_download" binding:"oneof=today_download total_download score createtime updatetime"`
}

type ListResponse struct {
	httputils.PageResponse[*Script] `json:",inline"`
}

// CreateRequest 创建脚本
type CreateRequest struct {
	mux.Meta    `path:"/scripts" method:"POST"`
	Content     string                      `form:"content" binding:"required,max=102400" label:"脚本详细描述"`
	Code        string                      `form:"code" binding:"required,max=10485760" label:"脚本代码"`
	Name        string                      `form:"name" binding:"max=128" label:"库的名字"`
	Description string                      `form:"description" binding:"max=10240" label:"库的描述"`
	Definition  string                      `form:"definition" binding:"max=10240" label:"库的定义文件"`
	Version     string                      `form:"version" binding:"max=32" label:"库的版本"`
	Type        script_entity.Type          `form:"type" binding:"required" label:"脚本类型"`   // 脚本类型：1 用户脚本 2 脚本引用库 3 订阅脚本(不支持)
	Public      script_entity.Public        `form:"public" binding:"required" label:"公开类型"` // 公开类型：1 公开 2 半公开
	Unwell      script_entity.UnwellContent `form:"unwell" binding:"required" label:"不适内容"` // 不适内容: 1 不适 2 适用
	Changelog   string                      `form:"changelog" binding:"max=102400" label:"更新日志"`
}

func (s *CreateRequest) Validate(ctx context.Context) error {
	if s.Type == script_entity.LibraryType {
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

// UpdateCodeRequest 更新脚本/库代码
type UpdateCodeRequest struct {
	mux.Meta `path:"/scripts/:id/code" method:"PUT"`
	ID       int64 `uri:"id" binding:"required"`
	//Name string `form:"name" binding:"max=128" label:"库的名字"`
	//Description string `form:"description" binding:"max=102400" label:"库的描述"`
	Version    string                      `binding:"required,max=128" form:"version" label:"库的版本号"`
	Content    string                      `binding:"required,max=102400" form:"content" label:"脚本详细描述"`
	Code       string                      `binding:"required,max=10485760" form:"code" label:"脚本代码"`
	Definition string                      `binding:"max=102400" form:"definition" label:"库的定义文件"`
	Changelog  string                      `binding:"max=102400" form:"changelog" label:"更新日志"`
	Public     script_entity.Public        `form:"public" binding:"required,number" label:"公开类型"` // 公开类型：1 公开 2 半公开
	Unwell     script_entity.UnwellContent `form:"unwell" binding:"required,number" label:"不适内容"`
}

type UpdateCodeResponse struct {
}

// MigrateEsRequest 全量迁移数据到es
type MigrateEsRequest struct {
	mux.Meta `path:"/scripts/migrate/es" method:"POST"`
}

type MigrateEsResponse struct {
}

// InfoRequest 获取脚本信息
type InfoRequest struct {
	mux.Meta `path:"/scripts/:id" method:"GET"`
	ID       int64 `uri:"id" binding:"required"`
}

type InfoResponse struct {
	*Script `json:",inline"`
	Content string `json:"content"`
}

// CodeRequest 获取脚本代码信息
type CodeRequest struct {
	mux.Meta `path:"/scripts/:id/code" method:"GET"`
	ID       int64 `uri:"id" binding:"required"`
}

type CodeResponse struct {
	*Script `json:",inline"`
}

// VersionListRequest 获取版本列表
type VersionListRequest struct {
	mux.Meta              `path:"/scripts/:id/versions" method:"GET"`
	httputils.PageRequest `form:",inline"`
	ID                    int64 `uri:"id" binding:"required"`
}

type VersionListResponse struct {
	httputils.PageResponse[*Code] `json:",inline"`
}

// VersionCodeRequest 获取指定版本代码
type VersionCodeRequest struct {
	mux.Meta `path:"/scripts/:id/versions/:version/code" method:"GET"`
	ID       int64  `uri:"id" binding:"required"`
	Version  string `uri:"version" binding:"required"`
}

type VersionCodeResponse struct {
	*Script `json:",inline"`
}

// StateRequest 获取脚本状态,脚本关注等
type StateRequest struct {
	mux.Meta `path:"/scripts/:id/state" method:"GET"`
	ID       int64 `uri:"id" binding:"required"`
}

type StateResponse struct {
	Watch script_entity.ScriptWatchLevel `json:"watch"`
}

// WatchRequest 关注脚本
type WatchRequest struct {
	mux.Meta `path:"/scripts/:id/watch" method:"POST"`
	ID       int64                          `uri:"id" binding:"required"`
	Watch    script_entity.ScriptWatchLevel `json:"watch" binding:"oneof=0 1 2 3"`
}

type WatchResponse struct {
}

// GetSettingRequest 获取脚本设置
type GetSettingRequest struct {
	mux.Meta `path:"/scripts/:id/setting" method:"GET"`
	ID       int64 `uri:"id" binding:"required"`
}

type GetSettingResponse struct {
	SyncUrl       string                 `json:"sync_url"`
	ContentUrl    string                 `json:"content_url"`
	DefinitionUrl string                 `json:"definition_url"`
	SyncMode      script_entity.SyncMode `json:"sync_mode"`
}

// UpdateSettingRequest 更新脚本设置
type UpdateSettingRequest struct {
	mux.Meta      `path:"/scripts/:id/setting" method:"PUT"`
	ID            int64                  `uri:"id" binding:"required"`
	Name          string                 `json:"name" binding:"max=128" label:"库的名字"`
	Description   string                 `json:"description" binding:"max=102400" label:"库的描述"`
	SyncUrl       string                 `json:"sync_url" binding:"omitempty,url,max=1024" label:"代码同步url"`
	ContentUrl    string                 `json:"content_url" binding:"omitempty,url,max=1024" label:"详细描述同步url"`
	DefinitionUrl string                 `json:"definition_url" binding:"omitempty,url,max=1024" label:"定义文件同步url"`
	SyncMode      script_entity.SyncMode `json:"sync_mode" binding:"number" label:"同步模式"`
}

type UpdateSettingResponse struct {
}

// ArchiveRequest 归档脚本
type ArchiveRequest struct {
	mux.Meta `path:"/scripts/:id/archive" method:"PUT"`
	ID       int64 `uri:"id" binding:"required"`
	Archive  bool  `json:"archive" binding:"omitempty,required"`
}

type ArchiveResponse struct {
}

// DeleteRequest 删除脚本
type DeleteRequest struct {
	mux.Meta `path:"/scripts/:id" method:"DELETE"`
	ID       int64 `uri:"id" binding:"required"`
}

type DeleteResponse struct {
}
