package issue_repo

import (
	"context"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/scriptscat/scriptlist/internal/model/entity/issue_entity"
)

type ScriptIssueCommentRepo interface {
	Find(ctx context.Context, issueId int64, id int64) (*issue_entity.ScriptIssueComment, error)
	FindPage(ctx context.Context, issueId int64, page httputils.PageRequest) ([]*issue_entity.ScriptIssueComment, int64, error)
	Create(ctx context.Context, scriptIssueComment *issue_entity.ScriptIssueComment) error
	Update(ctx context.Context, scriptIssueComment *issue_entity.ScriptIssueComment) error
	Delete(ctx context.Context, id int64) error

	FindAll(ctx context.Context, issueId int64) ([]*issue_entity.ScriptIssueComment, error)
}

var defaultScriptIssueComment ScriptIssueCommentRepo

func Comment() ScriptIssueCommentRepo {
	return defaultScriptIssueComment
}

func RegisterScriptIssueComment(i ScriptIssueCommentRepo) {
	defaultScriptIssueComment = i
}

type scriptIssueCommentRepo struct {
}

func NewScriptIssueComment() ScriptIssueCommentRepo {
	return &scriptIssueCommentRepo{}
}

func (u *scriptIssueCommentRepo) Find(ctx context.Context, issueId int64, id int64) (*issue_entity.ScriptIssueComment, error) {
	ret := &issue_entity.ScriptIssueComment{}
	if err := db.Ctx(ctx).Where("id=? and issue_id=? and status=?", id, issueId, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *scriptIssueCommentRepo) Create(ctx context.Context, scriptIssueComment *issue_entity.ScriptIssueComment) error {
	return db.Ctx(ctx).Create(scriptIssueComment).Error
}

func (u *scriptIssueCommentRepo) Update(ctx context.Context, scriptIssueComment *issue_entity.ScriptIssueComment) error {
	return db.Ctx(ctx).Updates(scriptIssueComment).Error
}

func (u *scriptIssueCommentRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Model(&issue_entity.ScriptIssueComment{ID: id}).Update("status", consts.DELETE).Error
}

func (u *scriptIssueCommentRepo) FindPage(ctx context.Context, issueId int64, page httputils.PageRequest) ([]*issue_entity.ScriptIssueComment, int64, error) {
	var list []*issue_entity.ScriptIssueComment
	var count int64
	if err := db.Ctx(ctx).Model(&issue_entity.ScriptIssueComment{}).Where("status=?", consts.ACTIVE).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Ctx(ctx).Where("issue_id=? and status=?", issueId, consts.ACTIVE).Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (u *scriptIssueCommentRepo) FindAll(ctx context.Context, issueId int64) ([]*issue_entity.ScriptIssueComment, error) {
	var list []*issue_entity.ScriptIssueComment
	if err := db.Ctx(ctx).Where("issue_id=? and status=?", issueId, consts.ACTIVE).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
