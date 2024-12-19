package api

import (
	"context"

	"github.com/scriptscat/scriptlist/internal/service/auth_svc"

	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/server/mux"
	_ "github.com/scriptscat/scriptlist/docs"
	"github.com/scriptscat/scriptlist/internal/controller/auth_ctr"
	"github.com/scriptscat/scriptlist/internal/controller/issue_ctr"
	"github.com/scriptscat/scriptlist/internal/controller/open_ctr"
	"github.com/scriptscat/scriptlist/internal/controller/resource_ctr"
	"github.com/scriptscat/scriptlist/internal/controller/script_ctr"
	"github.com/scriptscat/scriptlist/internal/controller/statistics_ctr"
	"github.com/scriptscat/scriptlist/internal/controller/user_ctr"
)

// Router 路由表
// @title    脚本站 API 文档
// @version  2.0.0
// @BasePath /api/v2
func Router(ctx context.Context, root *mux.Router) error {
	r := root.Group("/api/v2")
	auth := auth_ctr.NewAuth()
	// 用户-auth
	{
		// 调试环境开启简化登录
		if configs.Default().Debug && configs.Default().Env == configs.DEV {
			r.GET("/login/debug", auth.Debug())
		}
		r.GET("/login/oauth", auth.OAuthCallback())
	}
	// 用户
	{
		controller := user_ctr.NewUser()
		r.Group("/", auth_svc.Auth().RequireLogin(true)).Bind(
			controller.CurrentUser, // 获取当前用户信息
			controller.Follow,
			controller.GetWebhook,
			controller.RefreshWebhook,
			controller.GetConfig,
			controller.UpdateConfig,
			controller.Search,
			controller.LogOut,
		)
		r.GET("/users/:uid/avatar", controller.Avatar())
		r.Group("/").Bind(
			controller.Info, // 获取用户信息
		)
		r.Group("/", auth_svc.Auth().RequireLogin(false)).Bind(
			controller.Script,
			controller.GetFollow,
		)
	}
	// 脚本
	scriptCtr := script_ctr.NewScript()
	scriptCtr.Router(root, r)
	// 脚本评分
	scoreCtr := script_ctr.NewScore()
	scoreCtr.Router(r)
	// 群组管理
	scriptGroupCtr := script_ctr.NewGroup()
	scriptGroupCtr.Router(r)
	// 邀请码
	scriptInvCtr := script_ctr.NewAccessInvite()
	scriptInvCtr.Router(r)
	// 脚本权限
	scriptAccessCtr := script_ctr.NewAccess()
	scriptAccessCtr.Router(r)
	// 脚本反馈
	issueCtr := issue_ctr.NewIssue()
	issueCtr.Router(r)
	// 脚本反馈评论
	issueCommentCtr := issue_ctr.NewComment()
	issueCommentCtr.Router(r)
	// 脚本统计
	statisticsCtr := statistics_ctr.NewStatistics()
	statisticsCtr.Router(r)
	// 资源
	{
		controller := resource_ctr.NewResource()
		// 需要登录的路由组
		r.Group("/", auth_svc.Auth().RequireLogin(true)).Bind(
			controller.UploadImage,
		)
		// 不需要登录
		r.GET("/resource/image/:id", controller.ViewImage())
	}
	// 开放接口
	{
		controller := open_ctr.NewOpen()
		rg := r.Group("/")
		rg.GET("/open/crx-download/:id", controller.CrxDownload())
	}
	return nil
}
