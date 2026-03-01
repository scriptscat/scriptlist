package audit_svc

import (
	"context"

	"github.com/cago-frame/cago/pkg/utils/httputils"
	api "github.com/scriptscat/scriptlist/internal/api/audit"
	"github.com/scriptscat/scriptlist/internal/model/entity/audit_entity"
	"github.com/scriptscat/scriptlist/internal/repository/audit_repo"
)

type AuditLogSvc interface {
	List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error)
	ScriptList(ctx context.Context, req *api.ScriptListRequest) (*api.ScriptListResponse, error)
}

var defaultAuditLog AuditLogSvc = &auditLogSvc{}

func AuditLog() AuditLogSvc {
	return defaultAuditLog
}

type auditLogSvc struct{}

// List 查全局管理日志（仅返回管理员删除的脚本）
func (s *auditLogSvc) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	isAdmin := true
	opts := &audit_repo.ListOptions{
		Action:  string(audit_entity.ActionScriptDelete),
		IsAdmin: &isAdmin,
		Offset:  req.GetOffset(),
		Limit:   req.GetLimit(),
	}
	list, total, err := audit_repo.AuditLog().FindPage(ctx, opts)
	if err != nil {
		return nil, err
	}
	items := toAPIItems(list)
	return &api.ListResponse{
		PageResponse: httputils.PageResponse[*api.AuditLogItem]{
			Total: total,
			List:  items,
		},
	}, nil
}

// ScriptList 查单脚本日志
func (s *auditLogSvc) ScriptList(ctx context.Context, req *api.ScriptListRequest) (*api.ScriptListResponse, error) {
	opts := &audit_repo.ListOptions{
		TargetType: "script",
		TargetID:   req.ID,
		Offset:     req.GetOffset(),
		Limit:      req.GetLimit(),
	}
	list, total, err := audit_repo.AuditLog().FindPage(ctx, opts)
	if err != nil {
		return nil, err
	}
	items := toAPIItems(list)
	return &api.ScriptListResponse{
		PageResponse: httputils.PageResponse[*api.AuditLogItem]{
			Total: total,
			List:  items,
		},
	}, nil
}

func toAPIItems(list []*audit_entity.AuditLog) []*api.AuditLogItem {
	items := make([]*api.AuditLogItem, len(list))
	for i, v := range list {
		items[i] = &api.AuditLogItem{
			ID:         v.ID,
			UserID:     v.UserID,
			Username:   v.Username,
			Action:     v.Action,
			TargetType: v.TargetType,
			TargetID:   v.TargetID,
			TargetName: v.TargetName,
			IsAdmin:    v.IsAdmin,
			Reason:     v.Reason,
			Createtime: v.Createtime,
		}
	}
	return items
}
