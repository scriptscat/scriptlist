package audit_repo

import (
	"context"

	"github.com/cago-frame/cago/database/db"
	"github.com/scriptscat/scriptlist/internal/model/entity/audit_entity"
)

//go:generate mockgen -source=./audit_log.go -destination=./mock/audit_log.go
type AuditLogRepo interface {
	Create(ctx context.Context, auditLog *audit_entity.AuditLog) error
	FindPage(ctx context.Context, opts *ListOptions) ([]*audit_entity.AuditLog, int64, error)
}

// ListOptions 查询条件
type ListOptions struct {
	Action     string
	IsAdmin    *bool
	TargetType string
	TargetID   int64
	Offset     int
	Limit      int
}

var defaultAuditLog AuditLogRepo

func AuditLog() AuditLogRepo {
	return defaultAuditLog
}

func RegisterAuditLog(a AuditLogRepo) {
	defaultAuditLog = a
}

type auditLogRepo struct{}

func NewAuditLogRepo() AuditLogRepo {
	return &auditLogRepo{}
}

func (r *auditLogRepo) Create(ctx context.Context, auditLog *audit_entity.AuditLog) error {
	return db.Ctx(ctx).Create(auditLog).Error
}

func (r *auditLogRepo) FindPage(ctx context.Context, opts *ListOptions) ([]*audit_entity.AuditLog, int64, error) {
	var list []*audit_entity.AuditLog
	var count int64
	find := db.Ctx(ctx).Model(&audit_entity.AuditLog{})

	if opts.Action != "" {
		find = find.Where("action=?", opts.Action)
	}
	if opts.IsAdmin != nil {
		find = find.Where("is_admin=?", *opts.IsAdmin)
	}
	if opts.TargetType != "" {
		find = find.Where("target_type=?", opts.TargetType)
	}
	if opts.TargetID != 0 {
		find = find.Where("target_id=?", opts.TargetID)
	}

	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	limit := opts.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	if err := find.Order("createtime desc").Offset(opts.Offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}
