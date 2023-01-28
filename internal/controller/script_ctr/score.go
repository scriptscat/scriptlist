package script_ctr

import (
	"context"

	api "github.com/scriptscat/scriptlist/internal/api/script"
	service "github.com/scriptscat/scriptlist/internal/service/script_svc"
)

type Score struct {
}

func NewScore() *Score {
	return &Score{}
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
