package auth_svc

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/codfrm/cago/pkg/i18n"
	"github.com/scriptscat/scriptlist/internal/pkg/code"

	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/database/cache"
	cache2 "github.com/codfrm/cago/database/cache/cache"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/opentelemetry/trace"
	"github.com/codfrm/cago/pkg/utils"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
	api "github.com/scriptscat/scriptlist/internal/api/auth"
	"github.com/scriptscat/scriptlist/internal/model"
	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
	"github.com/scriptscat/scriptlist/internal/repository/user_repo"
	"github.com/scriptscat/scriptlist/pkg/oauth"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

const (
	TokenAuthMaxAge = 86400 * 7
	TokenAutoRegen  = 86400 * 3
)

//go:generate mockgen -source auth.go -destination mock/auth.go
type AuthSvc interface {
	// OAuthCallback 第三方登录
	OAuthCallback(ctx context.Context, req *api.OAuthCallbackRequest) (*api.OAuthCallbackResponse, error)
	// RequireLogin 处理鉴权中间件
	RequireLogin(force bool) gin.HandlerFunc
	// Get 获取用户鉴权信息
	Get(ctx context.Context) *model.AuthInfo
	// Login 登录获取token
	Login(ctx context.Context, uid int64) (*model.LoginToken, error)
	// Refresh 刷新token
	Refresh(ctx context.Context, uid int64, loginId, token string) (*model.LoginToken, error)
	// GetLoginToken 获取登录token信息
	GetLoginToken(ctx context.Context, uid int64, loginId, token string) (*model.LoginToken, error)
	// SetCtx 设置用户信息到上下文
	SetCtx(ctx context.Context, uid int64) (context.Context, error)
	Logout(ctx context.Context, uid int64, loginId, token string) (*model.LoginToken, error)
}

type authSvc struct {
}

var defaultAuth AuthSvc = &authSvc{}

func RegisterAuth(svc AuthSvc) {
	defaultAuth = svc
}

func Auth() AuthSvc {
	return defaultAuth
}

// OAuthCallback 第三方登录
func (a *authSvc) OAuthCallback(ctx context.Context, req *api.OAuthCallbackRequest) (*api.OAuthCallbackResponse, error) {
	// 请求论坛接口,进行登录
	config := &oauth.Config{}
	if err := configs.Default().Scan(ctx, "oauth.bbs", &config); err != nil {
		return nil, err
	}
	client := oauth.NewClient(config)
	resp, err := client.RequestAccessToken(ctx, req.Code)
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

// RequireLogin 鉴权中间件
func (a *authSvc) RequireLogin(force bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		loginId, _ := ctx.Cookie("login_id")
		token, _ := ctx.Cookie("token")
		if loginId == "" || token == "" {
			if force {
				httputils.HandleResp(ctx, httputils.NewError(http.StatusUnauthorized, -1, "未登录"))
				return
			}
			return
		}
		m := &model.LoginToken{}
		if err := cache.Ctx(ctx).Get("user:auth:login:" + loginId).Scan(m); err != nil {
			if errors.Is(err, cache2.ErrNil) {
				// 删除cookie
				ctx.SetCookie("login_id", "", -1, "/", "", false, true)
				ctx.SetCookie("token", "", -1, "/", "", false, true)
				if force {
					httputils.HandleResp(ctx, httputils.NewError(http.StatusUnauthorized, -1, "未登录"))
					return
				}
				return
			}
			httputils.HandleResp(ctx, err)
			return
		}
		if m.Token != token {
			// 删除cookie
			ctx.SetCookie("login_id", "", -1, "/", "", false, true)
			ctx.SetCookie("token", "", -1, "/", "", false, true)
			httputils.HandleResp(ctx, httputils.NewError(http.StatusUnauthorized, -1, "token error"))
			return
		}
		if m.Expired(TokenAuthMaxAge) {
			// 删除cookie
			ctx.SetCookie("login_id", "", -1, "/", "", false, true)
			ctx.SetCookie("token", "", -1, "/", "", false, true)
			httputils.HandleResp(ctx, httputils.NewError(http.StatusUnauthorized, -1, "token expired"))
			return
		}
		c, err := a.SetCtx(ctx.Request.Context(), m.UID)
		if err != nil {
			httputils.HandleResp(ctx, err)
			return
		}
		if ctx.Request.Method != http.MethodGet && !a.Get(c).EmailVerified {
			httputils.HandleResp(ctx, i18n.NewError(c, code.UserEmailNotVerified))
			return
		}
		ctx.Request = ctx.Request.WithContext(c)
	}
}

// SetCtx 设置用户信息到上下文
func (a *authSvc) SetCtx(ctx context.Context, uid int64) (context.Context, error) {
	// 获取用户信息
	user, err := user_repo.User().Find(ctx, uid)
	if err != nil {
		return nil, err
	}
	if err := user.IsBanned(ctx); err != nil {
		return nil, err
	}
	return a.SetCtxUser(ctx, user), nil
}

// SetCtxUser 设置用户信息到上下文
func (a *authSvc) SetCtxUser(ctx context.Context, user *user_entity.User) context.Context {
	// 设置用户信息,链路追踪和日志也添加上用户信息
	authInfo := &model.AuthInfo{
		UID:           user.ID,
		Username:      user.Username,
		Email:         user.Email,
		EmailVerified: !(user.Emailstatus == 0),
		AdminLevel:    model.AdminLevel(user.Adminid),
	}

	trace.SpanFromContext(ctx).SetAttributes(
		attribute.Int64("user_id", user.ID),
	)

	return context.WithValue(
		logger.WithContextLogger(ctx, logger.Ctx(ctx).
			With(zap.Int64("user_id", user.ID))),
		model.AuthInfo{}, authInfo)
}

// Get 获取用户鉴权信息
func (a *authSvc) Get(ctx context.Context) *model.AuthInfo {
	val := ctx.Value(model.AuthInfo{})
	if val == nil {
		return nil
	}
	return val.(*model.AuthInfo)
}

// Login 登录获取token
func (a *authSvc) Login(ctx context.Context, uid int64) (*model.LoginToken, error) {
	token := utils.RandString(16, utils.Mix)
	m := &model.LoginToken{
		ID:         utils.RandString(32, utils.Mix),
		UID:        uid,
		Token:      token,
		Createtime: time.Now().Unix(),
		Updatetime: time.Now().Unix(),
	}
	if err := cache.Ctx(ctx).Set("user:auth:login:"+m.ID, m, cache.Expiration(TokenAuthMaxAge*time.Second)).Err(); err != nil {
		return nil, err
	}
	return m, nil
}

// Refresh 刷新token
func (a *authSvc) Refresh(ctx context.Context, uid int64, loginId, token string) (*model.LoginToken, error) {
	// 判断token是否存在
	m := &model.LoginToken{}
	if err := cache.Ctx(ctx).Get("user:auth:login:" + loginId).Scan(m); err != nil {
		return nil, err
	}
	if m.UID != uid {
		return nil, httputils.NewError(http.StatusUnauthorized, -1, "token不匹配")
	}
	if m.Token != token {
		return nil, httputils.NewError(http.StatusUnauthorized, -1, "无效的token")
	}
	if m.Expired(TokenAuthMaxAge) {
		return nil, httputils.NewError(http.StatusUnauthorized, -1, "token已过期")
	}
	// 刷新token
	m.Token, m.LastToken = utils.RandString(16, utils.Mix), m.Token
	m.Updatetime = time.Now().Unix()
	if err := cache.Ctx(ctx).Set("user:auth:login:"+m.ID, m, cache.Expiration(TokenAuthMaxAge*time.Second)).Err(); err != nil {
		return nil, err
	}
	return m, nil
}

func (a *authSvc) GetLoginToken(ctx context.Context, uid int64, loginId, token string) (*model.LoginToken, error) {
	// 判断token是否存在
	m := &model.LoginToken{}
	if err := cache.Ctx(ctx).Get("user:auth:login:" + loginId).Scan(m); err != nil {
		return nil, err
	}
	if m.UID != uid {
		return nil, httputils.NewError(http.StatusUnauthorized, -1, "token不匹配")
	}
	if m.Token != token {
		return nil, httputils.NewError(http.StatusUnauthorized, -1, "无效的token")
	}
	if m.Expired(TokenAuthMaxAge) {
		return nil, httputils.NewError(http.StatusUnauthorized, -1, "token已过期")
	}
	return m, nil
}
func (a *authSvc) Logout(ctx context.Context, uid int64, loginId, token string) (*model.LoginToken, error) {
	m := &model.LoginToken{}
	if err := cache.Ctx(ctx).Get("user:auth:login:" + loginId).Scan(m); err != nil {
		return nil, err
	}
	if m.UID != uid {
		return nil, httputils.NewError(http.StatusUnauthorized, -1, "token不匹配")
	}
	if m.Token != token {
		return nil, httputils.NewError(http.StatusUnauthorized, -1, "无效的token")
	}
	err := cache.Ctx(ctx).Del("user:auth:login:" + m.ID)
	if err != nil {
		return nil, httputils.NewError(http.StatusUnauthorized, -1, "信息清除失败")
	}
	return m, nil
}
