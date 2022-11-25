package user

import (
	"context"
	"net/http"
	"strconv"

	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/middleware/sessions"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/trace"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
	api "github.com/scriptscat/scriptlist/internal/api/user"
	"github.com/scriptscat/scriptlist/internal/model"
	"github.com/scriptscat/scriptlist/pkg/oauth"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type IAuth interface {
	// OAuthCallback 第三方登录
	OAuthCallback(ctx context.Context, req *api.OAuthCallbackRequest) (*api.OAuthCallbackResponse, error)
	// Middleware 处理鉴权中间件
	Middleware(ctx *gin.Context)
	// Get 获取用户鉴权信息
	Get(ctx context.Context) *model.AuthInfo
}

type auth struct {
}

var defaultAuth = &auth{}

func Auth() IAuth {
	return defaultAuth
}

// OAuthCallback 第三方登录
func (a *auth) OAuthCallback(ctx context.Context, req *api.OAuthCallbackRequest) (*api.OAuthCallbackResponse, error) {
	// 请求论坛接口,进行登录
	config := &oauth.Config{}
	if err := configs.Default().Scan("oauth.bbs", &config); err != nil {
		return nil, err
	}
	client := oauth.NewClient(config)
	resp, err := client.RequestAccessToken(req.Code)
	if err != nil {
		return nil, err
	}
	user, err := client.RequestUser(resp.AccessToken)
	if err != nil {
		return nil, err
	}
	uid, _ := strconv.ParseInt(user.User.UID, 10, 64)
	return &api.OAuthCallbackResponse{
		RedirectUri: req.RedirectUri,
		UID:         uid,
	}, nil
}

// Middleware 鉴权中间件
func (a *auth) Middleware(ctx *gin.Context) {
	session := sessions.Ctx(ctx)
	uid, _ := session.Get("uid").(int64)
	if uid == 0 {
		httputils.HandleResp(ctx, httputils.NewError(
			http.StatusUnauthorized, -1, "未登录",
		))
		return
	}
	// 设置用户信息,链路追踪和日志也添加上用户信息
	authInfo := &model.AuthInfo{
		UID: uid,
	}
	trace.SpanFromContext(ctx.Request.Context()).SetAttributes(
		attribute.Int64("uid", uid),
	)
	ctx.Request = ctx.Request.WithContext(context.WithValue(
		logger.ContextWithLogger(ctx.Request.Context(), logger.Ctx(ctx.Request.Context()).
			With(zap.Int64("uid", uid))),
		model.AuthInfo{}, authInfo))
}

// Get 获取用户鉴权信息
func (a *auth) Get(ctx context.Context) *model.AuthInfo {
	val := ctx.Value(model.AuthInfo{})
	if val == nil {
		return nil
	}
	return val.(*model.AuthInfo)
}
