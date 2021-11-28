package service

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptlist/internal/domain/issue/broker"
	"github.com/scriptscat/scriptlist/internal/domain/issue/entity"
	"github.com/scriptscat/scriptlist/internal/domain/issue/repository"
	"github.com/scriptscat/scriptlist/internal/http/dto/request"
	"github.com/scriptscat/scriptlist/internal/pkg/cnt"
	"github.com/scriptscat/scriptlist/internal/pkg/errs"
	"github.com/scriptscat/scriptlist/pkg/utils"
)

const (
	CommentTypeComment = iota + 1
	CommentTypeChangeTitle
	CommentTypeChangeLabel
	CommentTypeOpen
	CommentTypeClose
	CommentTypeDelete
)

var Label = map[string]*entity.IssueLabel{
	"bug":      {Label: "bug", Name: "BUG", Description: "反馈一个bug", Color: "#ff0"},
	"feature":  {Label: "feature", Name: "新功能", Description: "请求增加新功能", Color: "#a2eeef"},
	"question": {Label: "question", Name: "问题", Description: "对脚本的使用存在问题", Color: "#d876e3"},
}

type Issue interface {
	List(script int64, keyword string, labels []string, status int, page request.Pages) ([]*entity.ScriptIssue, error)
	Issue(script, user int64, title, content string, label []string) (*entity.ScriptIssue, error)
	UpdateIssue(issue, user int64, title, content string) error
	GetIssue(issue int64) (*entity.ScriptIssue, error)
	DelIssue(issue, operator int64) error
	// 对issue的操作

	Open(issue, operator int64) error
	Close(issue, operator int64) error
	Label(issue, operator int64, label []string) error

	CommentList(issue int64, page request.Pages) ([]*entity.ScriptIssueComment, error)
	GetComment(commentId int64) (*entity.ScriptIssueComment, error)
	Comment(issue, user int64, content string) (*entity.ScriptIssueComment, error)
	UpdateComment(comment, user int64, content string) error
	DelComment(commentId int64) error
}

type issue struct {
	issueRepo   repository.Issue
	commentRepo repository.IssueComment
}

func NewIssue(issueRepo repository.Issue, commentRepo repository.IssueComment) Issue {
	return &issue{
		issueRepo:   issueRepo,
		commentRepo: commentRepo,
	}
}

func (i *issue) List(script int64, keyword string, label []string, status int, page request.Pages) ([]*entity.ScriptIssue, error) {
	return i.issueRepo.List(script, keyword, label, status, page)
}

func (i *issue) Issue(script, user int64, title, content string, labels []string) (*entity.ScriptIssue, error) {
	for _, v := range labels {
		if _, ok := Label[v]; !ok {
			return nil, errs.NewBadRequestError(1000, "错误的标签")
		}
	}
	issue := &entity.ScriptIssue{
		ScriptID:   script,
		UserID:     user,
		Title:      title,
		Content:    content,
		Labels:     strings.Join(labels, ","),
		Status:     cnt.ACTIVE,
		Createtime: time.Now().Unix(),
	}
	if err := i.issueRepo.Save(issue); err != nil {
		return nil, err
	}
	_ = broker.PublishScriptIssueCreate(issue.ID, script)
	return issue, nil
}

func (i *issue) UpdateIssue(issueId, user int64, title, content string) error {
	issue, err := i.GetIssue(issueId)
	if err != nil {
		return err
	}
	if issue.UserID != user {
		return errs.NewError(http.StatusForbidden, 1000, "没有权限修改反馈内容")
	}
	oldTitle := issue.Title
	issue.Title = title
	issue.Content = content
	issue.Updatetime = time.Now().Unix()
	if err := i.issueRepo.Save(issue); err != nil {
		return err
	}
	if oldTitle == title {
		return nil
	}
	return i.commentRepo.Save(&entity.ScriptIssueComment{
		IssueID:    issue.ID,
		UserID:     user,
		Content:    utils.MarshalJson(gin.H{"newtitle": title, "oldtitle": oldTitle}),
		Type:       CommentTypeChangeTitle,
		Status:     cnt.ACTIVE,
		Createtime: time.Now().Unix(),
	})
}

func (i *issue) GetIssue(issueId int64) (*entity.ScriptIssue, error) {
	issue, err := i.issueRepo.FindById(issueId)
	if err != nil {
		return nil, err
	}
	if issue == nil || issue.Status == cnt.DELETE {
		return nil, errs.NewError(http.StatusNotFound, 1000, "反馈不存在")
	}
	return issue, nil
}

func (i *issue) DelIssue(issueId int64, operator int64) error {
	if err := i.changeStatus(issueId, cnt.DELETE); err != nil {
		return err
	}
	return i.commentRepo.Save(&entity.ScriptIssueComment{
		IssueID:    issueId,
		UserID:     operator,
		Content:    "删除反馈",
		Type:       CommentTypeDelete,
		Status:     cnt.ACTIVE,
		Createtime: time.Now().Unix(),
	})
}

func (i *issue) Open(issueId, operator int64) error {
	if err := i.changeStatus(issueId, cnt.ACTIVE); err != nil {
		return err
	}
	return i.commentSave(&entity.ScriptIssueComment{
		IssueID:    issueId,
		UserID:     operator,
		Content:    "打开反馈",
		Type:       CommentTypeOpen,
		Status:     cnt.ACTIVE,
		Createtime: time.Now().Unix(),
	})
}

func (i *issue) Close(issueId, operator int64) error {
	if err := i.changeStatus(issueId, cnt.BAN); err != nil {
		return err
	}
	return i.commentSave(&entity.ScriptIssueComment{
		IssueID:    issueId,
		UserID:     operator,
		Content:    "关闭反馈",
		Type:       CommentTypeClose,
		Status:     cnt.ACTIVE,
		Createtime: time.Now().Unix(),
	})
}

func (i *issue) commentSave(comment *entity.ScriptIssueComment) error {
	if err := i.commentRepo.Save(comment); err != nil {
		return err
	}
	_ = broker.PublishScriptIssueCommentCreate(comment.IssueID, comment.ID)
	return nil
}

func (i *issue) changeStatus(issueId int64, status int) error {
	issue, err := i.GetIssue(issueId)
	if err != nil {
		return err
	}
	if issue.Status == status {
		return errs.NewBadRequestError(1000, "状态未发生改变")
	}
	issue.Status = status
	return i.issueRepo.Save(issue)
}

func (i *issue) Label(issueId, operator int64, label []string) error {
	issue, err := i.GetIssue(issueId)
	if err != nil {
		return err
	}
	oldLabel := strings.Split(issue.Labels, ",")
	oldLabelMap := make(map[string]struct{})
	for _, v := range oldLabel {
		oldLabelMap[v] = struct{}{}
	}
	labelMap := make(map[string]struct{})
	for _, v := range label {
		labelMap[v] = struct{}{}
	}
	update := make([]string, 0)
	add := make([]string, 0)
	for k := range labelMap {
		if _, ok := oldLabelMap[k]; !ok {
			add = append(add, k)
		}
		update = append(update, k)
	}
	del := make([]string, 0)
	for k := range oldLabelMap {
		if _, ok := labelMap[k]; !ok {
			del = append(del, k)
		}
	}
	if len(add) == 0 && len(del) == 0 {
		return errs.NewBadRequestError(1000, "标签没有发生改变")
	}
	for _, v := range update {
		if _, ok := Label[v]; !ok {
			return errs.NewBadRequestError(1000, "错误的标签")
		}
	}
	issue.Labels = strings.Join(update, ",")
	if err := i.issueRepo.Save(issue); err != nil {
		return err
	}
	return i.commentSave(&entity.ScriptIssueComment{
		IssueID:    issueId,
		UserID:     operator,
		Content:    utils.MarshalJson(gin.H{"add": add, "del": del}),
		Type:       CommentTypeChangeLabel,
		Status:     cnt.ACTIVE,
		Createtime: time.Now().Unix(),
	})
}

func (i *issue) CommentList(issue int64, page request.Pages) ([]*entity.ScriptIssueComment, error) {
	return i.commentRepo.List(issue, cnt.ACTIVE, page)
}

func (i *issue) Comment(issueId, user int64, content string) (*entity.ScriptIssueComment, error) {
	if content == "" {
		return nil, errs.NewBadRequestError(1000, "评论内容不能为空")
	}
	_, err := i.GetIssue(issueId)
	if err != nil {
		return nil, err
	}
	ret := &entity.ScriptIssueComment{
		IssueID:    issueId,
		UserID:     user,
		Content:    content,
		Type:       CommentTypeComment,
		Status:     cnt.ACTIVE,
		Createtime: time.Now().Unix(),
	}
	if err := i.commentSave(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (i *issue) GetComment(commentId int64) (*entity.ScriptIssueComment, error) {
	comment, err := i.commentRepo.FindById(commentId)
	if err != nil {
		return nil, err
	}
	if comment == nil || comment.Status != cnt.ACTIVE {
		return nil, errs.NewError(http.StatusNotFound, 1000, "评论不存在")
	}
	return comment, nil
}

func (i *issue) UpdateComment(commentId, user int64, content string) error {
	if content == "" {
		return errs.NewBadRequestError(1000, "评论内容不能为空")
	}
	comment, err := i.GetComment(commentId)
	if err != nil {
		return err
	}
	if comment.UserID != user || comment.Type != CommentTypeComment {
		return errs.NewError(http.StatusForbidden, 1001, "没有权限进行修改")
	}
	_, err = i.GetIssue(comment.IssueID)
	if err != nil {
		return err
	}
	comment.Content = content
	comment.Updatetime = time.Now().Unix()
	return i.commentRepo.Save(comment)
}

func (i *issue) DelComment(commentId int64) error {
	comment, err := i.GetComment(commentId)
	if err != nil {
		return err
	}
	comment.Status = cnt.DELETE
	comment.Updatetime = time.Now().Unix()
	return i.commentRepo.Save(comment)
}
