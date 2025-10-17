package issue_repo

import (
	"context"

	api "github.com/scriptscat/scriptlist/internal/api/issue"

	"github.com/cago-frame/cago/database/db"
	"github.com/cago-frame/cago/pkg/consts"
	"github.com/cago-frame/cago/pkg/utils/httputils"
	"github.com/scriptscat/scriptlist/internal/model/entity/issue_entity"
)

//go:generate mockgen -source=./script_issue.go -destination=./mock/script_issue.go
type ScriptIssueRepo interface {
	Find(ctx context.Context, scriptId int64, id int64) (*issue_entity.ScriptIssue, error)
	FindPage(ctx context.Context, req *api.ListRequest, page httputils.PageRequest) ([]*issue_entity.ScriptIssue, int64, error)
	Create(ctx context.Context, scriptIssue *issue_entity.ScriptIssue) error
	Update(ctx context.Context, scriptIssue *issue_entity.ScriptIssue) error
	Delete(ctx context.Context, scriptId int64, id int64) error
	CountByScript(ctx context.Context, scriptId int64, status int32) (int64, error)
}

var defaultScriptIssue ScriptIssueRepo

func Issue() ScriptIssueRepo {
	return defaultScriptIssue
}

func RegisterScriptIssue(i ScriptIssueRepo) {
	defaultScriptIssue = i
}

type scriptIssueRepo struct {
}

func NewScriptIssue() ScriptIssueRepo {
	return &scriptIssueRepo{}
}

func (u *scriptIssueRepo) Find(ctx context.Context, scriptId int64, id int64) (*issue_entity.ScriptIssue, error) {
	ret := &issue_entity.ScriptIssue{ID: id}
	if err := db.Ctx(ctx).Where("id=? and script_id=? and status!=?", id, scriptId, consts.DELETE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *scriptIssueRepo) Create(ctx context.Context, scriptIssue *issue_entity.ScriptIssue) error {
	return db.Ctx(ctx).Create(scriptIssue).Error
}

func (u *scriptIssueRepo) Update(ctx context.Context, scriptIssue *issue_entity.ScriptIssue) error {
	return db.Ctx(ctx).Model(&scriptIssue).
		Select("*").Updates(scriptIssue).Error
}

func (u *scriptIssueRepo) Delete(ctx context.Context, scriptId int64, id int64) error {
	return db.Ctx(ctx).Model(&issue_entity.ScriptIssue{}).
		Where("id=? and script_id=?", id, scriptId).Update("status", consts.DELETE).Error
}

func (u *scriptIssueRepo) FindPage(ctx context.Context, req *api.ListRequest, page httputils.PageRequest) ([]*issue_entity.ScriptIssue, int64, error) {
	var list []*issue_entity.ScriptIssue
	var count int64
	db := db.Ctx(ctx).Model(&issue_entity.ScriptIssue{}).
		Where("script_id=? and status!=?", req.ScriptID, consts.DELETE)
	if req.Status > 0 {
		db = db.Where("status=?", req.Status)
	}
	if req.Keyword != "" {
		db = db.Where("title like ?", "%"+req.Keyword+"%")
	}
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("createtime desc").Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

// CountByScript 统计脚本反馈数量
func (u *scriptIssueRepo) CountByScript(ctx context.Context, scriptId int64, status int32) (int64, error) {
	var count int64
	db := db.Ctx(ctx).Model(&issue_entity.ScriptIssue{}).
		Where("script_id=? and status!=?", scriptId, consts.DELETE)
	if status > 0 {
		db = db.Where("status=?", status)
	}
	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
