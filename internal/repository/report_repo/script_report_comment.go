//go:generate mockgen -source=./script_report_comment.go -destination=./mock/script_report_comment.go -package=mock_report_repo

package report_repo

import (
	"context"

	"github.com/cago-frame/cago/database/db"
	"github.com/cago-frame/cago/pkg/consts"
	"github.com/scriptscat/scriptlist/internal/model/entity/report_entity"
)

type ScriptReportCommentRepo interface {
	Find(ctx context.Context, reportId int64, id int64) (*report_entity.ScriptReportComment, error)
	FindAll(ctx context.Context, reportId int64) ([]*report_entity.ScriptReportComment, error)
	Create(ctx context.Context, comment *report_entity.ScriptReportComment) error
	Delete(ctx context.Context, id int64) error
	CountByReport(ctx context.Context, reportId int64) (int64, error)
}

var defaultScriptReportComment ScriptReportCommentRepo

func Comment() ScriptReportCommentRepo {
	return defaultScriptReportComment
}

func RegisterScriptReportComment(i ScriptReportCommentRepo) {
	defaultScriptReportComment = i
}

type scriptReportCommentRepo struct{}

func NewScriptReportComment() ScriptReportCommentRepo {
	return &scriptReportCommentRepo{}
}

func (r *scriptReportCommentRepo) Find(ctx context.Context, reportId int64, id int64) (*report_entity.ScriptReportComment, error) {
	ret := &report_entity.ScriptReportComment{}
	if err := db.Ctx(ctx).Where("id=? and report_id=? and status!=?", id, reportId, consts.DELETE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (r *scriptReportCommentRepo) FindAll(ctx context.Context, reportId int64) ([]*report_entity.ScriptReportComment, error) {
	var list []*report_entity.ScriptReportComment
	if err := db.Ctx(ctx).Where("report_id=? and status!=?", reportId, consts.DELETE).
		Order("createtime asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *scriptReportCommentRepo) Create(ctx context.Context, comment *report_entity.ScriptReportComment) error {
	return db.Ctx(ctx).Create(comment).Error
}

func (r *scriptReportCommentRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Model(&report_entity.ScriptReportComment{}).
		Where("id=?", id).Update("status", consts.DELETE).Error
}

func (r *scriptReportCommentRepo) CountByReport(ctx context.Context, reportId int64) (int64, error) {
	var count int64
	if err := db.Ctx(ctx).Model(&report_entity.ScriptReportComment{}).
		Where("report_id=? and status!=?", reportId, consts.DELETE).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
