package system_svc

import (
	"context"
	"time"

	api "github.com/scriptscat/scriptlist/internal/api/system"
	"github.com/scriptscat/scriptlist/internal/model/entity/feedback_entity"
	"github.com/scriptscat/scriptlist/internal/repository/feedback_repo"
)

type SystemSvc interface {
	// Feedback 用户反馈请求
	Feedback(ctx context.Context, req *api.FeedbackRequest) (*api.FeedbackResponse, error)
}

type systemSvc struct {
}

var defaultSystem = &systemSvc{}

func System() SystemSvc {
	return defaultSystem
}

// Feedback 用户反馈请求
func (s *systemSvc) Feedback(ctx context.Context, req *api.FeedbackRequest) (*api.FeedbackResponse, error) {
	if err := feedback_repo.Feedback().Create(ctx, &feedback_entity.Feedback{
		Reason:     req.Reason,
		Content:    req.Content,
		ClientIp:   req.ClientIp(),
		Createtime: time.Now().Unix(),
	}); err != nil {
		return nil, err
	}
	return &api.FeedbackResponse{}, nil
}
