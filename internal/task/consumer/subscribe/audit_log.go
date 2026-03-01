package subscribe

import (
	"context"
	"time"

	"github.com/cago-frame/cago/pkg/broker/broker"
	"github.com/cago-frame/cago/pkg/logger"
	"github.com/scriptscat/scriptlist/internal/model/entity/audit_entity"
	"github.com/scriptscat/scriptlist/internal/repository/audit_repo"
	"github.com/scriptscat/scriptlist/internal/task/producer"
	"go.uber.org/zap"
)

// AuditLog 审计日志消费者
type AuditLog struct{}

func (a *AuditLog) Subscribe(ctx context.Context) error {
	if err := producer.SubscribeScriptDelete(ctx, a.scriptDelete, broker.Group("audit")); err != nil {
		return err
	}
	if err := producer.SubscribeScriptCodeUpdate(ctx, a.scriptCodeUpdate, broker.Group("audit")); err != nil {
		return err
	}
	if err := producer.SubscribeScriptCreate(ctx, a.scriptCreate, broker.Group("audit")); err != nil {
		return err
	}
	return nil
}

func (a *AuditLog) scriptDelete(ctx context.Context, msg *producer.ScriptDeleteMsg) error {
	log := &audit_entity.AuditLog{
		UserID:     msg.OperatorUID,
		Username:   msg.OperatorUsername,
		Action:     audit_entity.ActionScriptDelete,
		TargetType: "script",
		TargetID:   msg.Script.ID,
		TargetName: msg.Script.Name,
		IsAdmin:    msg.IsAdmin,
		Reason:     msg.Reason,
		Createtime: time.Now().Unix(),
	}
	if err := audit_repo.AuditLog().Create(ctx, log); err != nil {
		logger.Ctx(ctx).Error("审计日志写入失败", zap.Error(err), zap.Int64("script_id", msg.Script.ID))
		return err
	}
	return nil
}

func (a *AuditLog) scriptCodeUpdate(ctx context.Context, msg *producer.ScriptCodeUpdateMsg) error {
	log := &audit_entity.AuditLog{
		UserID:     msg.OperatorUID,
		Username:   msg.OperatorUsername,
		Action:     audit_entity.ActionScriptUpdate,
		TargetType: "script",
		TargetID:   msg.Script.ID,
		TargetName: msg.Script.Name,
		IsAdmin:    msg.IsAdmin,
		Createtime: time.Now().Unix(),
	}
	if err := audit_repo.AuditLog().Create(ctx, log); err != nil {
		logger.Ctx(ctx).Error("审计日志写入失败", zap.Error(err), zap.Int64("script_id", msg.Script.ID))
		return err
	}
	return nil
}

func (a *AuditLog) scriptCreate(ctx context.Context, msg *producer.ScriptCreateMsg) error {
	log := &audit_entity.AuditLog{
		UserID:     msg.OperatorUID,
		Username:   msg.OperatorUsername,
		Action:     audit_entity.ActionScriptCreate,
		TargetType: "script",
		TargetID:   msg.Script.ID,
		TargetName: msg.Script.Name,
		IsAdmin:    msg.IsAdmin,
		Createtime: time.Now().Unix(),
	}
	if err := audit_repo.AuditLog().Create(ctx, log); err != nil {
		logger.Ctx(ctx).Error("审计日志写入失败", zap.Error(err), zap.Int64("script_id", msg.Script.ID))
		return err
	}
	return nil
}
