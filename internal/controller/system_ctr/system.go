package system_ctr

import (
	"github.com/cago-frame/cago/database/redis"
	"github.com/cago-frame/cago/pkg/limit"
	"github.com/gin-gonic/gin"
	api "github.com/scriptscat/scriptlist/internal/api/system"
	"github.com/scriptscat/scriptlist/internal/service/system_svc"
)

type System struct {
	limit limit.Limit
}

func NewSystem() *System {
	return &System{
		limit: limit.NewPeriodLimit(60, 2, redis.Default(), "system:feedback"),
	}
}

// Feedback 用户反馈请求
func (s *System) Feedback(ctx *gin.Context, req *api.FeedbackRequest) (*api.FeedbackResponse, error) {
	// 根据ip限流
	resp, err := s.limit.FuncTake(ctx, ctx.ClientIP(), func() (interface{}, error) {
		req.SetClientIp(ctx.ClientIP())
		return system_svc.System().Feedback(ctx, req)
	})
	if err != nil {
		return nil, err
	}
	return resp.(*api.FeedbackResponse), nil
}
