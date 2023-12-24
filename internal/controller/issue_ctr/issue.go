package issue_ctr

import (
	"context"
	"strconv"

	"github.com/codfrm/cago/pkg/utils/muxutils"
	"github.com/codfrm/cago/server/mux"
	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptlist/internal/service/script_svc"

	"github.com/codfrm/cago/database/redis"
	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/limit"
	api "github.com/scriptscat/scriptlist/internal/api/issue"
	"github.com/scriptscat/scriptlist/internal/model/entity/issue_entity"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	"github.com/scriptscat/scriptlist/internal/service/issue_svc"
)

type Issue struct {
	limit limit.Limit
}

func NewIssue() *Issue {
	return &Issue{
		limit: limit.NewPeriodLimit(
			300, 10, redis.Default(), "limit:create:issue",
		),
	}
}

func (i *Issue) skipSelf() script_svc.CheckOption {
	return script_svc.WithCheckSkip(func(ctx *gin.Context) (bool, error) {
		return issue_svc.Issue().CtxIssue(ctx).UserID == auth_svc.Auth().Get(ctx).UID, nil
	})
}

func (i *Issue) Router(r *mux.Router) {
	muxutils.BindTree(r, []*muxutils.RouterTree{{
		Middleware: []gin.HandlerFunc{script_svc.Script().RequireScript()},
		Handler: []interface{}{
			// 无需登录
			i.List,
			muxutils.Use(issue_svc.Issue().RequireIssue()).Append(i.GetIssue),
			// 需要登录
			muxutils.Use(auth_svc.Auth().RequireLogin(true)).Append(i.CreateIssue),
			// 需要登录且issue存在
			&muxutils.RouterTree{
				Middleware: []gin.HandlerFunc{
					auth_svc.Auth().RequireLogin(true),
					issue_svc.Issue().RequireIssue(),
				},
				Handler: []interface{}{
					i.Watch,
					// 归档了不允许操作
					muxutils.Use(script_svc.Script().IsArchive()).Append(
						i.GetWatch,
						muxutils.Use(script_svc.Access().
							CheckHandler("issue", "manage", i.skipSelf())).Append(
							i.Open,
							i.Close,
							i.UpdateLabels,
						),
						muxutils.Use(script_svc.Access().
							CheckHandler("issue", "delete")).Append(
							i.Delete,
						),
					),
				},
			},
		},
	}})
}

// List 获取脚本反馈列表
func (i *Issue) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	return issue_svc.Issue().List(ctx, req)
}

// CreateIssue 创建脚本反馈
func (i *Issue) CreateIssue(ctx context.Context, req *api.CreateIssueRequest) (*api.CreateIssueResponse, error) {
	for _, v := range req.Labels {
		_, ok := issue_entity.Label[v]
		if !ok {
			return nil, i18n.NewError(ctx, code.IssueLabelNotExist)
		}
	}
	resp, err := i.limit.FuncTake(ctx, strconv.FormatInt(auth_svc.Auth().Get(ctx).UID, 10), func() (interface{}, error) {
		return issue_svc.Issue().CreateIssue(ctx, req)
	})
	if err != nil {
		return nil, err
	}
	return resp.(*api.CreateIssueResponse), nil
}

// GetIssue 获取issue信息
func (i *Issue) GetIssue(ctx context.Context, req *api.GetIssueRequest) (*api.GetIssueResponse, error) {
	return issue_svc.Issue().GetIssue(ctx, req)
}

// GetWatch 获取issue关注状态
func (i *Issue) GetWatch(ctx context.Context, req *api.GetWatchRequest) (*api.GetWatchResponse, error) {
	return issue_svc.Issue().GetWatch(ctx, req)
}

// Watch 关注issue
func (i *Issue) Watch(ctx context.Context, req *api.WatchRequest) (*api.WatchResponse, error) {
	resp, err := i.limit.FuncTake(ctx, strconv.FormatInt(auth_svc.Auth().Get(ctx).UID, 10), func() (interface{}, error) {
		return issue_svc.Issue().Watch(ctx, auth_svc.Auth().Get(ctx).UID, req)
	})
	if err != nil {
		return nil, err
	}
	return resp.(*api.WatchResponse), nil
}

// Close 关闭issue
func (i *Issue) Close(ctx context.Context, req *api.CloseRequest) (*api.CloseResponse, error) {
	resp, err := i.limit.FuncTake(ctx, strconv.FormatInt(auth_svc.Auth().Get(ctx).UID, 10), func() (interface{}, error) {
		return issue_svc.Issue().Close(ctx, req)
	})
	if err != nil {
		return nil, err
	}
	return resp.(*api.CloseResponse), nil
}

// Open 打开issue
func (i *Issue) Open(ctx context.Context, req *api.OpenRequest) (*api.OpenResponse, error) {
	resp, err := i.limit.FuncTake(ctx, strconv.FormatInt(auth_svc.Auth().Get(ctx).UID, 10), func() (interface{}, error) {
		return issue_svc.Issue().Open(ctx, req)
	})
	if err != nil {
		return nil, err
	}
	return resp.(*api.OpenResponse), nil
}

// Delete 删除issue
func (i *Issue) Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error) {
	return issue_svc.Issue().Delete(ctx, req)
}

// UpdateLabels 更新issue标签
func (i *Issue) UpdateLabels(ctx context.Context, req *api.UpdateLabelsRequest) (*api.UpdateLabelsResponse, error) {
	for _, v := range req.Labels {
		_, ok := issue_entity.Label[v]
		if !ok {
			return nil, i18n.NewError(ctx, code.IssueLabelNotExist)
		}
	}
	return issue_svc.Issue().UpdateLabels(ctx, req)
}
