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
	SessionID     string `json:"session_id" binding:"required"` // 当前会话id,随机生成一串字符,可以是uuid
	ScriptID      int64  `json:"script_id" binding:"required"`  // 脚本id,在初始化脚本时设置(例如new ScriptStatistics(1))
	VisitorID     string `json:"visitor_id" binding:"required"` // 访客id,从浏览器指纹或者GM_value中获取
	OperationPage string `json:"operation_page"`                // 操作页面,当前页面的链接
	InstallPage   string `json:"install_page"`                  // 安装页面,从GM_info中提取
	Version       string `json:"version" binding:"required"`    // 版本,从GM_info中提取
	VisitTime     int64  `json:"visit_time" binding:"required"` // 访问时间,脚本执行的时间,unix时间戳(10位)
	Duration      int32  `json:"duration"`                      // 停留时长(秒)
	ExitTime      int64  `json:"exit_time"`                     // 退出时间
	Iframe        bool   `json:"iframe"`                        // 是否在iframe中运行
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
