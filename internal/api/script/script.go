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
	PostID               int64                          `json:"post_id"`
	Name                 string                         `json:"name"`
	Description          string                         `json:"description"`
	Category             []*CategoryList                `json:"category"`
	Status               int64                          `json:"status"`
	Score                int64                          `json:"score"`
	ScoreNum             int64                          `json:"score_num"`
	Type                 script_entity.Type             `json:"type"`
	Public               int                            `json:"public"`
	Unwell               int                            `json:"unwell"`
	Archive              int                            `json:"archive"`
	Danger               int                            `json:"danger"`
	EnablePreRelease     script_entity.EnablePreRelease `json:"enable_pre_release"`
	TodayInstall         int64                          `json:"today_install"`
	TotalInstall         int64                          `json:"total_install"`
	Createtime           int64                          `json:"createtime"`
	Updatetime           int64                          `json:"updatetime"`
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
	Meta                 string                         `json:"meta,omitempty"`
	MetaJson             interface{}                    `json:"meta_json"`
	ScriptID             int64                          `json:"script_id"`
	Version              string                         `json:"version"`
	Changelog            string                         `json:"changelog"`
	IsPreRelease         script_entity.EnablePreRelease `json:"is_pre_release"`
	Status               int64                          `json:"status"`
	Createtime           int64                          `json:"createtime"`
	Code                 string                         `json:"code,omitempty"`
	Definition           string                         `json:"definition,omitempty"`
}

// ListRequest 获取脚本列表
type ListRequest struct {
	mux.Meta              `path:"/scripts" method:"GET"`
	httputils.PageRequest `form:",inline"`
	Keyword               string `form:"keyword"`
	Domain                string `form:"domain"`
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
	Type        script_entity.Type          `form:"type" binding:"required,oneof=1 2 3" label:"脚本类型"` // 脚本类型：1 用户脚本 2 脚本引用库 3 订阅脚本(不支持)
	Public      script_entity.Public        `form:"public" binding:"required,oneof=1 2" label:"公开类型"` // 公开类型：1 公开 2 半公开
	Unwell      script_entity.UnwellContent `form:"unwell" binding:"required,oneof=1 2" label:"不适内容"` // 不适内容: 1 不适 2 适用
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
	Version      string                         `binding:"required,max=128" form:"version" label:"库的版本号"`
	Content      string                         `binding:"required,max=102400" form:"content" label:"脚本详细描述"`
	Code         string                         `binding:"required,max=10485760" form:"code" label:"脚本代码"`
	Definition   string                         `binding:"max=102400" form:"definition" label:"库的定义文件"`
	Changelog    string                         `binding:"max=102400" form:"changelog" label:"更新日志"`
	IsPreRelease script_entity.EnablePreRelease `form:"is_pre_release" json:"is_pre_release" binding:"omitempty,oneof=0 1 2" label:"是否预发布"`
	//Public       script_entity.Public           `form:"public" binding:"required,oneof=1 2" label:"公开类型"` // 公开类型：1 公开 2 半公开
	//Unwell       script_entity.UnwellContent    `form:"unwell" binding:"required,oneof=1 2" label:"不适内容"`
}

type UpdateCodeResponse struct {
}

// DeleteCodeRequest 删除脚本/库代码
type DeleteCodeRequest struct {
	mux.Meta `path:"/scripts/:id/code/:codeId" method:"DELETE"`
	ID       int64 `uri:"id" binding:"required"`
	CodeID   int64 `uri:"codeId" binding:"required"`
}

type DeleteCodeResponse struct {
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
	Content string `json:"content"`
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
	SyncUrl          string                         `json:"sync_url"`
	ContentUrl       string                         `json:"content_url"`
	DefinitionUrl    string                         `json:"definition_url"`
	SyncMode         script_entity.SyncMode         `json:"sync_mode"`
	EnablePreRelease script_entity.EnablePreRelease `json:"enable_pre_release"`
	GrayControls     []*script_entity.GrayControl   `json:"gray_controls"`
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
	Sync      bool
	SyncError string
}

// UpdateLibInfoRequest 更新库信息
type UpdateLibInfoRequest struct {
	mux.Meta    `path:"/scripts/:id/lib-info" method:"PUT"`
	Name        string `json:"name" binding:"max=128" label:"库的名字"`
	Description string `json:"description" binding:"max=102400" label:"库的描述"`
}

type UpdateLibInfoResponse struct {
}

// UpdateSyncSettingRequest 更新同步配置
type UpdateSyncSettingRequest struct {
	mux.Meta      `path:"/scripts/:id/sync" method:"PUT"`
	SyncUrl       string                 `json:"sync_url" binding:"omitempty,url,max=1024" label:"代码同步url"`
	ContentUrl    string                 `json:"content_url" binding:"omitempty,url,max=1024" label:"详细描述同步url"`
	DefinitionUrl string                 `json:"definition_url" binding:"omitempty,url,max=1024" label:"定义文件同步url"`
	SyncMode      script_entity.SyncMode `json:"sync_mode" binding:"number" label:"同步模式"`
}

type UpdateSyncSettingResponse struct {
	Sync      bool
	SyncError string
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

// UpdateCodeSettingRequest 更新脚本设置
type UpdateCodeSettingRequest struct {
	mux.Meta     `path:"/scripts/:id/code/:codeId" method:"PUT"`
	ID           int64                          `uri:"id" binding:"required"`
	CodeID       int64                          `uri:"codeId" binding:"required"`
	Changelog    string                         `json:"changelog" binding:"max=102400" label:"更新日志"`
	IsPreRelease script_entity.EnablePreRelease `json:"is_pre_release" binding:"oneof=1 2" label:"是否预发布"`
}

type UpdateCodeSettingResponse struct {
}

// UpdateScriptPublicRequest 更新脚本公开类型
type UpdateScriptPublicRequest struct {
	mux.Meta `path:"/scripts/:id/public" method:"PUT"`
	ID       int64                `uri:"id" binding:"required"`
	Public   script_entity.Public `json:"public" binding:"required,oneof=1 2" label:"公开类型"`
}

type UpdateScriptPublicResponse struct {
}

// UpdateScriptUnwellRequest 更新脚本不适内容
type UpdateScriptUnwellRequest struct {
	mux.Meta `path:"/scripts/:id/unwell" method:"PUT"`
	ID       int64                       `uri:"id" binding:"required"`
	Unwell   script_entity.UnwellContent `json:"unwell" binding:"required,oneof=1 2" label:"不适内容"`
}

type UpdateScriptUnwellResponse struct {
}

// UpdateScriptGrayRequest 更新脚本灰度策略
type UpdateScriptGrayRequest struct {
	mux.Meta         `path:"/scripts/:id/gray" method:"PUT"`
	EnablePreRelease script_entity.EnablePreRelease `json:"enable_pre_release" binding:"oneof=1 2" label:"是否开启预发布"`
	GrayControls     []*script_entity.GrayControl   `json:"gray_controls" binding:"required" label:"灰度策略"`
}

type UpdateScriptGrayResponse struct {
}

// WebhookRequest 处理webhook请求
type WebhookRequest struct {
	mux.Meta         `path:"/webhook/:uid" method:"POST"`
	UID              int64  `uri:"uid" binding:"required"`
	UA               string `header:"User-Agent" binding:"required"`
	XHubSignature256 string `header:"X-Hub-Signature-256" binding:"required"`
}

type WebhookResponse struct {
}

// LastScoreRequest 最新评分脚本
type LastScoreRequest struct {
	mux.Meta              `path:"/scripts/last-score" method:"GET"`
	httputils.PageRequest `form:",inline"`
}

type LastScoreResponse struct {
	httputils.PageResponse[*Script] `json:",inline"`
}
