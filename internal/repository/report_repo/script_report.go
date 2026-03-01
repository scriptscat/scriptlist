//go:generate mockgen -source=./script_report.go -destination=./mock/script_report.go -package=mock_report_repo

package report_repo

import (
	"context"

	"github.com/cago-frame/cago/database/db"
	"github.com/cago-frame/cago/pkg/consts"
	"github.com/cago-frame/cago/pkg/utils/httputils"
	"github.com/scriptscat/scriptlist/internal/model/entity/report_entity"
)

type ScriptReportRepo interface {
	Find(ctx context.Context, scriptId int64, id int64) (*report_entity.ScriptReport, error)
	FindPage(ctx context.Context, scriptId int64, status int32, page httputils.PageRequest) ([]*report_entity.ScriptReport, int64, error)
	Create(ctx context.Context, report *report_entity.ScriptReport) error
	Update(ctx context.Context, report *report_entity.ScriptReport) error
	Delete(ctx context.Context, scriptId int64, id int64) error
	CountByScript(ctx context.Context, scriptId int64, status int32) (int64, error)
}

var defaultScriptReport ScriptReportRepo

func Report() ScriptReportRepo {
	return defaultScriptReport
}

func RegisterScriptReport(i ScriptReportRepo) {
	defaultScriptReport = i
}

type scriptReportRepo struct{}

func NewScriptReport() ScriptReportRepo {
	return &scriptReportRepo{}
}

func (r *scriptReportRepo) Find(ctx context.Context, scriptId int64, id int64) (*report_entity.ScriptReport, error) {
	ret := &report_entity.ScriptReport{}
	if err := db.Ctx(ctx).Where("id=? and script_id=? and status!=?", id, scriptId, consts.DELETE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (r *scriptReportRepo) FindPage(ctx context.Context, scriptId int64, status int32, page httputils.PageRequest) ([]*report_entity.ScriptReport, int64, error) {
	var list []*report_entity.ScriptReport
	var count int64
	query := db.Ctx(ctx).Model(&report_entity.ScriptReport{}).
		Where("script_id=? and status!=?", scriptId, consts.DELETE)
	if status > 0 {
		query = query.Where("status=?", status)
	}
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("createtime desc").Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (r *scriptReportRepo) Create(ctx context.Context, report *report_entity.ScriptReport) error {
	return db.Ctx(ctx).Create(report).Error
}

func (r *scriptReportRepo) Update(ctx context.Context, report *report_entity.ScriptReport) error {
	return db.Ctx(ctx).Model(report).Select("*").Updates(report).Error
}

func (r *scriptReportRepo) Delete(ctx context.Context, scriptId int64, id int64) error {
	return db.Ctx(ctx).Model(&report_entity.ScriptReport{}).
		Where("id=? and script_id=?", id, scriptId).Update("status", consts.DELETE).Error
}

func (r *scriptReportRepo) CountByScript(ctx context.Context, scriptId int64, status int32) (int64, error) {
	var count int64
	query := db.Ctx(ctx).Model(&report_entity.ScriptReport{}).
		Where("script_id=? and status!=?", scriptId, consts.DELETE)
	if status > 0 {
		query = query.Where("status=?", status)
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
