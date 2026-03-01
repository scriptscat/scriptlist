package audit_ctr

import (
	"github.com/gin-gonic/gin"
	api "github.com/scriptscat/scriptlist/internal/api/audit"
	"github.com/scriptscat/scriptlist/internal/service/audit_svc"
)

type AuditLog struct{}

func NewAuditLog() *AuditLog {
	return &AuditLog{}
}

// List 全局管理日志
func (a *AuditLog) List(ctx *gin.Context, req *api.ListRequest) (*api.ListResponse, error) {
	return audit_svc.AuditLog().List(ctx, req)
}

// ScriptList 单脚本管理日志
func (a *AuditLog) ScriptList(ctx *gin.Context, req *api.ScriptListRequest) (*api.ScriptListResponse, error) {
	return audit_svc.AuditLog().ScriptList(ctx, req)
}
