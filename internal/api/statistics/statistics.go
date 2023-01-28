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
