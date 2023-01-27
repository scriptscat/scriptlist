package issue_repo

import (
	"context"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/scriptscat/scriptlist/internal/model/entity/issue_entity"
)

type ScriptIssueWatchRepo interface {
	Find(ctx context.Context, id int64) (*issue_entity.ScriptIssueWatch, error)
	FindPage(ctx context.Context, page httputils.PageRequest) ([]*issue_entity.ScriptIssueWatch, int64, error)
	Create(ctx context.Context, scriptIssueWatch *issue_entity.ScriptIssueWatch) error
	Update(ctx context.Context, scriptIssueWatch *issue_entity.ScriptIssueWatch) error
	Delete(ctx context.Context, id int64) error

	FindAll(ctx context.Context, issueId int64) ([]*issue_entity.ScriptIssueWatch, error)
	FindByUser(ctx context.Context, id, userId int64) (*issue_entity.ScriptIssueWatch, error)
}

var defaultScriptIssueWatch ScriptIssueWatchRepo

func Watch() ScriptIssueWatchRepo {
	return defaultScriptIssueWatch
}

func RegisterScriptIssueWatch(i ScriptIssueWatchRepo) {
	defaultScriptIssueWatch = i
}

type scriptIssueWatchRepo struct {
}

func NewScriptIssueWatch() ScriptIssueWatchRepo {
	return &scriptIssueWatchRepo{}
}

func (u *scriptIssueWatchRepo) Find(ctx context.Context, id int64) (*issue_entity.ScriptIssueWatch, error) {
	ret := &issue_entity.ScriptIssueWatch{ID: id}
	if err := db.Ctx(ctx).Where("status=?", consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *scriptIssueWatchRepo) Create(ctx context.Context, scriptIssueWatch *issue_entity.ScriptIssueWatch) error {
	return db.Ctx(ctx).Create(scriptIssueWatch).Error
}

func (u *scriptIssueWatchRepo) Update(ctx context.Context, scriptIssueWatch *issue_entity.ScriptIssueWatch) error {
	return db.Ctx(ctx).Updates(scriptIssueWatch).Error
}

func (u *scriptIssueWatchRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Model(&issue_entity.ScriptIssueWatch{ID: id}).Update("status", consts.DELETE).Error
}

func (u *scriptIssueWatchRepo) FindPage(ctx context.Context, page httputils.PageRequest) ([]*issue_entity.ScriptIssueWatch, int64, error) {
	var list []*issue_entity.ScriptIssueWatch
	var count int64
	if err := db.Ctx(ctx).Model(&issue_entity.ScriptIssueWatch{}).Where("status=?", consts.ACTIVE).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Ctx(ctx).Where("status=?", consts.ACTIVE).Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (u *scriptIssueWatchRepo) FindAll(ctx context.Context, issueId int64) ([]*issue_entity.ScriptIssueWatch, error) {
	var list []*issue_entity.ScriptIssueWatch
	if err := db.Ctx(ctx).Where("status=? and issue_id=?", consts.ACTIVE, issueId).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (u *scriptIssueWatchRepo) FindByUser(ctx context.Context, id, userId int64) (*issue_entity.ScriptIssueWatch, error) {
	ret := &issue_entity.ScriptIssueWatch{}
	if err := db.Ctx(ctx).Where("issue_id=? and user_id=?", id, userId).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}
