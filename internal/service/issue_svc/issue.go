package issue_svc

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/scriptscat/scriptlist/internal/service/script_svc"

	"github.com/cago-frame/cago/pkg/consts"
	"github.com/cago-frame/cago/pkg/i18n"
	"github.com/cago-frame/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
	api "github.com/scriptscat/scriptlist/internal/api/issue"
	"github.com/scriptscat/scriptlist/internal/model/entity/issue_entity"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
	"github.com/scriptscat/scriptlist/internal/repository/issue_repo"
	"github.com/scriptscat/scriptlist/internal/repository/user_repo"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	"github.com/scriptscat/scriptlist/internal/task/producer"
)

type contextKey int

const (
	issueCtxKey contextKey = iota
	issueCommentCtxKey
)

type IssueSvc interface {
	// List 获取脚本反馈列表
	List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error)
	// CreateIssue 创建脚本反馈
	CreateIssue(ctx context.Context, req *api.CreateIssueRequest) (*api.CreateIssueResponse, error)
	// GetIssue 获取issue信息
	GetIssue(ctx context.Context, req *api.GetIssueRequest) (*api.GetIssueResponse, error)
	// GetWatch 获取issue关注状态
	GetWatch(ctx context.Context, req *api.GetWatchRequest) (*api.GetWatchResponse, error)
	// Watch 关注issue
	Watch(ctx context.Context, userId int64, req *api.WatchRequest) (*api.WatchResponse, error)
	// Close 关闭issue
	Close(ctx context.Context, req *api.CloseRequest) (*api.CloseResponse, error)
	// Open 打开issue
	Open(ctx context.Context, req *api.OpenRequest) (*api.OpenResponse, error)
	// Delete 删除issue
	Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error)
	// UpdateLabels 更新issue标签
	UpdateLabels(ctx context.Context, req *api.UpdateLabelsRequest) (*api.UpdateLabelsResponse, error)
	// RequireIssue 需要issue存在
	RequireIssue() gin.HandlerFunc
	// CtxIssue 获取issue
	CtxIssue(ctx context.Context) *issue_entity.ScriptIssue
}

type issueSvc struct {
}

var defaultIssue = &issueSvc{}

func Issue() IssueSvc {
	return defaultIssue
}

// List 获取脚本反馈列表
func (i *issueSvc) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	list, total, err := issue_repo.Issue().FindPage(ctx, req, req.PageRequest)
	if err != nil {
		return nil, err
	}
	resp := &api.ListResponse{
		PageResponse: httputils.PageResponse[*api.Issue]{
			Total: total,
			List:  make([]*api.Issue, len(list)),
		},
	}
	for n, v := range list {
		resp.List[n], _ = i.ToIssue(ctx, v)
	}
	return resp, nil
}

func (i *issueSvc) ToIssue(ctx context.Context, issue *issue_entity.ScriptIssue) (*api.Issue, error) {
	ret := &api.Issue{
		ID:         issue.ID,
		ScriptID:   issue.ScriptID,
		Title:      issue.Title,
		Labels:     issue.GetLabels(),
		Status:     issue.Status,
		Createtime: issue.Createtime,
		Updatetime: issue.Updatetime,
	}
	user, err := user_repo.User().Find(ctx, issue.UserID)
	if err != nil {
		return nil, err
	}
	ret.UserInfo = user.UserInfo()
	commentCount, err := issue_repo.Comment().CountByIssue(ctx, issue.ID)
	if err != nil {
		return nil, err
	}
	ret.CommentCount = commentCount
	return ret, nil
}

// CreateIssue 创建脚本反馈
func (i *issueSvc) CreateIssue(ctx context.Context, req *api.CreateIssueRequest) (*api.CreateIssueResponse, error) {
	issue := &issue_entity.ScriptIssue{
		ScriptID:   req.ScriptID,
		UserID:     auth_svc.Auth().Get(ctx).UID,
		Title:      req.Title,
		Content:    req.Content,
		Labels:     strings.Join(req.Labels, ","),
		Status:     consts.ACTIVE,
		Createtime: time.Now().Unix(),
	}
	if err := issue_repo.Issue().Create(ctx, issue); err != nil {
		return nil, err
	}
	// 发布消息
	return &api.CreateIssueResponse{ID: issue.ID}, producer.PublishIssueCreate(ctx, script_svc.Script().CtxScript(ctx), issue)
}

// GetIssue 获取issue信息
func (i *issueSvc) GetIssue(ctx context.Context, req *api.GetIssueRequest) (*api.GetIssueResponse, error) {
	issue := i.CtxIssue(ctx)
	ret, _ := i.ToIssue(ctx, issue)
	return &api.GetIssueResponse{
		Issue:   ret,
		Content: issue.Content,
	}, nil
}

// GetWatch 获取issue关注状态
func (i *issueSvc) GetWatch(ctx context.Context, req *api.GetWatchRequest) (*api.GetWatchResponse, error) {
	m, err := issue_repo.Watch().FindByUser(ctx, req.IssueID, auth_svc.Auth().Get(ctx).UID)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return &api.GetWatchResponse{
			Watch: false,
		}, nil
	}
	return &api.GetWatchResponse{
		Watch: m.Status == consts.ACTIVE,
	}, nil
}

// Watch 关注issue
func (i *issueSvc) Watch(ctx context.Context, userId int64, req *api.WatchRequest) (*api.WatchResponse, error) {
	m, err := issue_repo.Watch().FindByUser(ctx, req.IssueID, userId)
	if err != nil {
		return nil, err
	}
	var watch int32 = consts.DELETE
	if req.Watch {
		watch = consts.ACTIVE
	}
	if m == nil {
		m = &issue_entity.ScriptIssueWatch{
			UserID:     userId,
			IssueID:    req.IssueID,
			Status:     watch,
			Createtime: time.Now().Unix(),
		}
		if err := issue_repo.Watch().Create(ctx, m); err != nil {
			return nil, err
		}
	} else {
		m.Status = watch
		m.Updatetime = time.Now().Unix()
		if err := issue_repo.Watch().Update(ctx, m); err != nil {
			return nil, err
		}
	}
	return &api.WatchResponse{}, nil
}

// Close 关闭issue
func (i *issueSvc) Close(ctx context.Context, req *api.CloseRequest) (*api.CloseResponse, error) {
	// 检查是否有权限操作
	issue := i.CtxIssue(ctx)
	issue.Status = consts.AUDIT
	issue.Updatetime = time.Now().Unix()
	comment := &issue_entity.ScriptIssueComment{
		IssueID:    req.IssueID,
		UserID:     auth_svc.Auth().Get(ctx).UID,
		Content:    "关闭反馈",
		Type:       issue_entity.CommentTypeClose,
		Status:     consts.ACTIVE,
		Createtime: time.Now().Unix(),
	}
	if err := issue_repo.Issue().Update(ctx, issue); err != nil {
		return nil, err
	}
	if err := issue_repo.Comment().Create(ctx, comment); err != nil {
		return nil, err
	}
	// 发布消息
	resp, _ := Comment().ToComment(ctx, comment)
	return &api.CloseResponse{
		Comment: resp,
	}, producer.PublishCommentCreate(ctx, script_svc.Script().CtxScript(ctx), issue, comment)
}

// Open 打开issue
func (i *issueSvc) Open(ctx context.Context, req *api.OpenRequest) (*api.OpenResponse, error) {
	// 检查是否有权限操作
	issue := i.CtxIssue(ctx)
	issue.Status = consts.ACTIVE
	issue.Updatetime = time.Now().Unix()
	comment := &issue_entity.ScriptIssueComment{
		IssueID:    req.IssueID,
		UserID:     auth_svc.Auth().Get(ctx).UID,
		Content:    "打开反馈",
		Type:       issue_entity.CommentTypeOpen,
		Status:     consts.ACTIVE,
		Createtime: time.Now().Unix(),
	}
	if err := issue_repo.Issue().Update(ctx, issue); err != nil {
		return nil, err
	}
	if err := issue_repo.Comment().Create(ctx, comment); err != nil {
		return nil, err
	}
	// 发布消息
	resp, _ := Comment().ToComment(ctx, comment)
	return &api.OpenResponse{
		Comment: resp,
	}, producer.PublishCommentCreate(ctx, script_svc.Script().CtxScript(ctx), issue, comment)
}

// Delete 删除issue
func (i *issueSvc) Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error) {
	issue := i.CtxIssue(ctx)
	if err := issue_repo.Issue().Delete(ctx, script_svc.Script().CtxScript(ctx).ID, issue.ID); err != nil {
		return nil, err
	}
	return nil, nil
}

// UpdateLabels 更新issue标签
func (i *issueSvc) UpdateLabels(ctx context.Context, req *api.UpdateLabelsRequest) (*api.UpdateLabelsResponse, error) {
	issue := i.CtxIssue(ctx)
	// 对比标签变更
	var oldLabel []string
	if issue.Labels != "" {
		oldLabel = strings.Split(issue.Labels, ",")
	}
	oldLabelMap := make(map[string]struct{})
	for _, v := range oldLabel {
		oldLabelMap[v] = struct{}{}
	}
	labelMap := make(map[string]struct{})
	for _, v := range req.Labels {
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
		return nil, i18n.NewError(ctx, code.IssueLabelNotChange)
	}
	issue.Labels = strings.Join(update, ",")
	issue.Updatetime = time.Now().Unix()
	if err := issue_repo.Issue().Update(ctx, issue); err != nil {
		return nil, err
	}
	content, err := json.Marshal(gin.H{"add": add, "del": del})
	if err != nil {
		return nil, err
	}
	comment := &issue_entity.ScriptIssueComment{
		IssueID:    issue.ID,
		UserID:     auth_svc.Auth().Get(ctx).UID,
		Content:    string(content),
		Type:       issue_entity.CommentTypeChangeLabel,
		Status:     consts.ACTIVE,
		Createtime: time.Now().Unix(),
	}
	if err := issue_repo.Comment().Create(ctx, comment); err != nil {
		return nil, err
	}
	// 标签变更就不发布消息了
	resp, _ := Comment().ToComment(ctx, comment)
	return &api.UpdateLabelsResponse{
		Comment: resp,
	}, nil
}

func (i *issueSvc) RequireIssue() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sIssueId := ctx.Param("issueId")
		if sIssueId == "" {
			httputils.HandleResp(ctx, httputils.NewError(http.StatusNotFound, -1, "反馈ID不能为空"))
			return
		}
		issueId, err := strconv.ParseInt(sIssueId, 10, 64)
		if err != nil {
			httputils.HandleResp(ctx, err)
			return
		}
		script := script_svc.Script().CtxScript(ctx)
		issue, err := issue_repo.Issue().Find(ctx, script.ID, issueId)
		if err != nil {
			httputils.HandleResp(ctx, err)
			return
		}
		if err := issue.CheckOperate(ctx); err != nil {
			httputils.HandleResp(ctx, err)
			return
		}

		ctx.Request = ctx.Request.WithContext(context.WithValue(
			ctx.Request.Context(), issueCtxKey, issue,
		))

	}
}

func (i *issueSvc) CtxIssue(ctx context.Context) *issue_entity.ScriptIssue {
	return ctx.Value(issueCtxKey).(*issue_entity.ScriptIssue)
}
