package statistics

import "github.com/codfrm/cago/server/mux"

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
	SessionID     string `json:"session_id" binding:"required"`     // 会话id,随机生成
	ScriptID      int64  `json:"script_id" binding:"required"`      // 脚本id
	VisitorID     string `json:"visitor_id" binding:"required"`     // 访客id
	OperationPage string `json:"operation_page" binding:"required"` // 操作页面
	InstallPage   string `json:"install_page" binding:"required"`   // 安装页面
	Version       string `json:"version" binding:"required"`        // 版本
	VisitTime     int64  `json:"visit_time" binding:"required"`     // 访问时间
	Duration      int32  `json:"duration"`                          // 停留时长(秒)
	ExitTime      int64  `json:"exit_time"`                         // 退出时间
	UA            string
	IP            string
}

type CollectResponse struct {
}

// RealtimeChartRequest 实时统计数据图表
type RealtimeChartRequest struct {
	mux.Meta `path:"/statistics/:id/realtime/chart" method:"GET"`
	ID       int64 `uri:"id" binding:"required"`
}

type RealtimeChartResponse struct {
	Chart *Chart `json:"chart"`
}

// RealtimeRequest 实时统计数据
type RealtimeRequest struct {
	mux.Meta `path:"/statistics/:id/realtime" method:"GET"`
	ID       int64 `uri:"id" binding:"required"`
}

type RealtimeResponse struct {
}

// BasicInfoRequest 基本统计信息
type BasicInfoRequest struct {
	mux.Meta `path:"/statistics/:id/basic" method:"GET"`
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

type BasicInfoResponse struct {
	Limit           *Limit    `json:"limit"`
	PV              *Overview `json:"pv"`
	UV              *Overview `json:"uv"`
	UseTime         *Overview `json:"use_time"`
	NewUser         *Overview `json:"new_user"`
	OldUser         *Overview `json:"old_user"`
	Origin          *PieChart `json:"origin"`
	Version         *PieChart `json:"version"`
	OperationDomain *PieChart `json:"operation_domain"`
	System          *PieChart `json:"system"`
	Browser         *PieChart `json:"browser"`
}

// UserOriginRequest 用户来源统计
type UserOriginRequest struct {
	mux.Meta `path:"/statistics/:id/user/origin" method:"GET"`
	ID       int64 `uri:"id" binding:"required"`
}

type UserOriginResponse struct {
}
