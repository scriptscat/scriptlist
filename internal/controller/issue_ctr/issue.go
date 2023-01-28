package issue_ctr

import (
	"context"
	"strconv"

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
	limit *limit.PeriodLimit
}

func NewIssue() *Issue {
	return &Issue{
		limit: limit.NewPeriodLimit(
			300, 10, redis.Default(), "limit:create:issue",
		),
	}
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
