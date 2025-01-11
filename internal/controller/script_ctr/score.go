package script_ctr

import (
	"context"

	"github.com/codfrm/cago/pkg/utils/muxutils"
	"github.com/codfrm/cago/server/mux"
	"github.com/gin-gonic/gin"
	api "github.com/scriptscat/scriptlist/internal/api/script"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	service "github.com/scriptscat/scriptlist/internal/service/script_svc"
)

type Score struct {
}

func NewScore() *Score {
	return &Score{}
}

func (s *Score) Router(r *mux.Router) {
	muxutils.BindTree(r, []*muxutils.RouterTree{{
		// 无需登录
		Handler: []interface{}{
			s.ScoreList,
		},
	}, {
		// 需要登录
		Middleware: []gin.HandlerFunc{
			auth_svc.Auth().RequireLogin(true),
			service.Script().RequireScript(),
		},
		Handler: []interface{}{
			muxutils.Use(service.Script().IsArchive()).Append(
				s.PutScore,
			),
			s.SelfScore,
			// 只有管理员才能删除评分
			&muxutils.RouterTree{
				Middleware: []gin.HandlerFunc{
					service.Access().CheckHandler("script", "delete:score"),
				},
				Handler: []interface{}{s.DelScore},
			},
			//管理员和作者才可以回复评分
			&muxutils.RouterTree{
				Middleware: []gin.HandlerFunc{
					service.Access().CheckHandler("script", "reply:score"),
				},
				Handler: []interface{}{s.ReplyScore},
			},
		},
	}})
}

// PutScore 脚本评分
func (s *Score) PutScore(ctx context.Context, req *api.PutScoreRequest) (*api.PutScoreResponse, error) {
	return service.Score().PutScore(ctx, req)
}

// ScoreList 获取脚本评分列表
func (s *Score) ScoreList(ctx context.Context, req *api.ScoreListRequest) (*api.ScoreListResponse, error) {
	return service.Score().ScoreList(ctx, req)
}

// SelfScore 用于获取自己对脚本的评价
func (s *Score) SelfScore(ctx context.Context, req *api.SelfScoreRequest) (*api.SelfScoreResponse, error) {
	return service.Score().SelfScore(ctx, req)
}

// DelScore 用于删除脚本的评价，注意，只有管理员才有权限删除评价
func (s *Score) DelScore(ctx context.Context, req *api.DelScoreRequest) (*api.DelScoreResponse, error) {
	return service.Score().DelScore(ctx, req)
}

// ReplyScore 对用户评分进行回复，只有管理员和作者才可以进行回复
func (s *Score) ReplyScore(ctx context.Context, req *api.ReplyScoreRequest) (*api.ReplyScoreResponse, error) {
	return service.Score().ReplyScore(ctx, req)
}
