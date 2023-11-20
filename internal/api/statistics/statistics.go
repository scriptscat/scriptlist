package statistics

import (
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/cago/server/mux"
)

// ScriptRequest 脚本统计数据
type ScriptRequest struct {
	mux.Meta `path:"/script/:id/statistics" method:"GET"`
	ID       int64 `uri:"id" binding:"required"`
}

type Overview struct {
	Today     int64 `json:"today"`
	Yesterday int64 `json:"yesterday"`
	Week      int64 `json:"week"`
}

type Chart struct {
	X []string `json:"x"`
	Y []int64  `json:"y"`
}

type DUChart struct {
	Download *Chart `json:"download"`
	Update   *Chart `json:"update"`
}

type ScriptResponse struct {
	PagePv     *Overview `json:"page_pv"`
	PageUv     *Overview `json:"page_uv"`
	DownloadUv *Overview `json:"download_uv"`
	UpdateUv   *Overview `json:"update_uv"`
	UvChart    *DUChart  `json:"uv_chart"`
	PvChart    *DUChart  `json:"pv_chart"`
}

// ScriptRealtimeRequest 脚本实时统计数据
type ScriptRealtimeRequest struct {
	mux.Meta `path:"/script/:id/statistics/realtime" method:"GET"`
	ID       int64 `uri:"id" binding:"required"`
}

type ScriptRealtimeResponse struct {
	Download *Chart `json:"download"`
	Update   *Chart `json:"update"`
}

// CollectRequest 统计数据收集
type CollectRequest struct {
	mux.Meta      `path:"/statistics/collect" method:"POST"`
	SessionID     string `form:"session_id" json:"session_id" binding:"required"` // 当前会话id,随机生成一串字符,可以是uuid
	ScriptID      int64  `form:"script_id" json:"script_id"`                      // 脚本id,在初始化脚本时设置(例如new ScriptStatistics(1))
	StatisticsKey string `form:"statistics_key" json:"statistics_key"`            // 统计key
	VisitorID     string `form:"visitor_id" json:"visitor_id" binding:"required"` // 访客id,从浏览器指纹或者GM_value中获取
	OperationPage string `form:"operation_page" json:"operation_page"`            // 操作页面,当前页面的链接
	InstallPage   string `form:"install_page" json:"install_page"`                // 安装页面,从GM_info中提取
	Version       string `form:"version" json:"version" binding:"required"`       // 版本,从GM_info中提取
	VisitTime     int64  `form:"visit_time" json:"visit_time" binding:"required"` // 访问时间,脚本执行的时间,unix时间戳(10位)
	Duration      int32  `form:"duration" json:"duration"`                        // 停留时长(秒)
	ExitTime      int64  `form:"exit_time" json:"exit_time"`                      // 退出时间
	Iframe        bool   `form:"iframe" json:"iframe"`                            // 是否在iframe中运行
	UA            string
	IP            string
}

type CollectResponse struct {
}

type CollectWhitelistRequest struct {
	mux.Meta      `path:"/statistics/collect/whitelist" method:"POST"`
	StatisticsKey string `form:"statistics_key" json:"statistics_key"` // 统计key
}

type CollectWhitelistResponse struct {
	Whitelist []string `json:"whitelist"`
}

// RealtimeChartRequest 实时统计数据图表
type RealtimeChartRequest struct {
	mux.Meta `path:"/statistics/:id/realtime/chart" method:"GET"`
	ID       int64 `uri:"id" binding:"required"`
}

type RealtimeChartResponse struct {
	Chart *Chart `json:"chart"`
}

// VisitListRequest 访问列表
type VisitListRequest struct {
	mux.Meta              `path:"/statistics/:id/visit" method:"GET"`
	httputils.PageRequest `form:",inline"`
	ID                    int64 `uri:"id" binding:"required"`
}

type VisitItem struct {
	//SessionID     string `json:"session_id"`
	VisitorID string `json:"visitor_id"`
	//OperationHost string `json:"operation_host"`
	OperationPage string `json:"operation_page"`
	Duration      int32  `json:"duration"`
	VisitTime     int64  `json:"visit_time"`
	ExitTime      int64  `json:"exit_time"`
}

type VisitListResponse struct {
	httputils.PageResponse[*VisitItem] `json:",inline"`
}

// AdvancedInfoRequest 高级统计信息
type AdvancedInfoRequest struct {
	mux.Meta `path:"/statistics/:id/advanced" method:"GET"`
	ID       int64 `uri:"id" binding:"required"`
}

type Limit struct {
	// 限额
	Quota int64 `json:"quota"`
	// 用额
	Usage int64 `json:"usage"`
}

// PieChart 饼图
type PieChart struct {
	Key   string `json:"key"`
	Value int64  `json:"value"`
}

type AdvancedInfoResponse struct {
	StatisticsKey string      `json:"statistics_key"`
	Whitelist     []string    `json:"whitelist"`
	Limit         *Limit      `json:"limit"`
	PV            *Overview   `json:"pv"`
	UV            *Overview   `json:"uv"`
	IP            *Overview   `json:"ip"`
	UseTime       *Overview   `json:"use_time"`
	NewOldUser    []*PieChart `json:"new_old_user"`
	Version       []*PieChart `json:"version"`
	System        []*PieChart `json:"system"`
	Browser       []*PieChart `json:"browser"`
}

// UserOriginRequest 用户来源统计
type UserOriginRequest struct {
	mux.Meta              `path:"/statistics/:id/user-origin" method:"GET"`
	httputils.PageRequest `form:",inline"`
	ID                    int64 `uri:"id" binding:"required"`
}

type UserOriginResponse struct {
	httputils.PageResponse[*PieChart] `json:",inline"`
}

// VisitDomainRequest 访问域名统计
type VisitDomainRequest struct {
	mux.Meta              `path:"/statistics/:id/visit-domain" method:"GET"`
	httputils.PageRequest `form:",inline"`
	ID                    int64 `uri:"id" binding:"required"`
}

type VisitDomainResponse struct {
	httputils.PageResponse[*PieChart] `json:",inline"`
}

type UpdateWhitelistRequest struct {
	mux.Meta  `path:"/statistics/:id/whitelist" method:"PUT"`
	ID        int64    `uri:"id" binding:"required"`
	Whitelist []string `json:"whitelist"`
}

type UpdateWhitelistResponse struct {
	Whitelist []string `json:"whitelist"`
}
