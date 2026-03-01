package audit

import (
	"github.com/cago-frame/cago/pkg/utils/httputils"
	"github.com/cago-frame/cago/server/mux"
	"github.com/scriptscat/scriptlist/internal/model/entity/audit_entity"
)

// AuditLogItem 审计日志项
type AuditLogItem struct {
	ID         int64               `json:"id"`
	UserID     int64               `json:"user_id"`
	Username   string              `json:"username"`
	Action     audit_entity.Action `json:"action"`
	TargetType string              `json:"target_type"`
	TargetID   int64               `json:"target_id"`
	TargetName string              `json:"target_name"`
	IsAdmin    bool                `json:"is_admin"`
	Reason     string              `json:"reason"`
	Createtime int64               `json:"createtime"`
}

// ListRequest 全局管理日志（公开，仅返回管理员操作）
type ListRequest struct {
	mux.Meta              `path:"/audit-logs" method:"GET"`
	httputils.PageRequest `form:",inline"`
	Action                string `form:"action"`
}

type ListResponse struct {
	httputils.PageResponse[*AuditLogItem] `json:",inline"`
}

// ScriptListRequest 单脚本日志（需要脚本 manage 权限）
type ScriptListRequest struct {
	mux.Meta              `path:"/scripts/:id/audit-logs" method:"GET"`
	httputils.PageRequest `form:",inline"`
	ID                    int64 `uri:"id" binding:"required"`
}

type ScriptListResponse struct {
	httputils.PageResponse[*AuditLogItem] `json:",inline"`
}
