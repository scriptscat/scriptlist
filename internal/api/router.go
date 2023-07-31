package api

import (
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/server/mux"
	"github.com/gin-contrib/cors"
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
func Router(root *mux.Router) error {
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
		r.Group("/", auth.Middleware(true)).Bind(
			controller.CurrentUser, // 获取当前用户信息
			controller.Follow,
			controller.GetWebhook,
			controller.RefreshWebhook,
			controller.GetConfig,
			controller.UpdateConfig,
		)
		r.GET("/users/:uid/avatar", controller.Avatar())
		r.Group("/").Bind(
			controller.Info, // 获取用户信息
		)
		r.Group("/", auth.Middleware(false)).Bind(
			controller.Script,
			controller.GetFollow,
		)
	}
	// 脚本
	scriptCtr := script_ctr.NewScript()
	{
		// 需要用户登录的路由组
		r.Group("/", auth.Middleware(true)).Bind(
			scriptCtr.Create,
			scriptCtr.UpdateCode,
			scriptCtr.MigrateEs,
			scriptCtr.Watch,
			scriptCtr.GetSetting,
			scriptCtr.UpdateSetting,
		)
		// 慢慢迁移至中间件处理
		r.Group("/", auth.Middleware(true), scriptCtr.Middleware()).Bind(
			scriptCtr.Delete,
			scriptCtr.UpdateCodeSetting,
			scriptCtr.UpdateScriptPublic,
			scriptCtr.UpdateScriptUnwell,
			scriptCtr.UpdateScriptGray,
			scriptCtr.Archive,
			scriptCtr.DeleteCode,
		)
		// 处理下载
		root.GET("/scripts/code/:id/*name", auth.Middleware(false), scriptCtr.Download(false))
		root.GET("/scripts/pre/:id/*name", auth.Middleware(false), scriptCtr.Download(true))
		root.GET("/lib/:id/:version/*name", auth.Middleware(false), scriptCtr.DownloadLib())
		// 无需用户登录的路由组
		r.Group("/").Bind(
			scriptCtr.List,
			scriptCtr.LastScore,
			scriptCtr.Info,
			scriptCtr.Code,
			scriptCtr.VersionList,
			scriptCtr.VersionCode,
		)
		r.Any("/webhook/:uid", scriptCtr.Webhook)
		// 半登录
		r.Group("/", auth.Middleware(false)).Bind(
			scriptCtr.State,
		)
	}
	{
		issueComment := issue_ctr.NewComment()
		// 脚本反馈
		{
			controller := issue_ctr.NewIssue()
			// 需要登录的路由组
			r.Group("/", auth.Middleware(true), issueComment.Middleware()).Bind(
				controller.CreateIssue,
				controller.Open,
				controller.Close,
				controller.Watch,
				controller.GetWatch,
				controller.Delete,
				controller.UpdateLabels,
			)
			// 不需要登录
			r.Group("/", issueComment.Middleware()).Bind(
				controller.List,
				controller.GetIssue,
			)
		}
		// 脚本反馈评论
		{
			// 需要登录的路由组
			r.Group("/", auth.Middleware(true), issueComment.Middleware()).Bind(
				issueComment.CreateComment,
				issueComment.DeleteComment,
			)
			// 不需要登录
			r.Group("/", issueComment.Middleware()).Bind(
				issueComment.ListComment,
			)
		}
	}
	//脚本评分
	{
		controller := script_ctr.NewScore()
		//需要用户登录才能评分的路由组
		r.Group("/", auth.Middleware(true)).Bind(
			controller.PutScore,
			controller.SelfScore,
			controller.DelScore,
		)
		//无需用户登录的路由组
		r.Group("/").Bind(
			controller.ScoreList,
		)
	}
	// 统计
	{
		controller := statistics_ctr.NewStatistics()
		r.Group("/", auth.Middleware(true), controller.Middleware()).Bind(
			controller.Script,
			controller.ScriptRealtime,

			controller.AdvancedInfo,
			controller.UserOrigin,
			controller.RealtimeChart,
			controller.VisitList,
			controller.VisitDomain,
			controller.UpdateWhitelist,
		)
		rg := r.Group("/", cors.Default())
		rg.OPTIONS("/statistics/collect")
		rg.OPTIONS("/statistics/collect/whitelist")
		rg.Bind(
			controller.Collect,
			controller.CollectWhitelist,
		)
	}
	// 资源
	{
		controller := resource_ctr.NewResource()
		// 需要登录的路由组
		r.Group("/", auth.Middleware(true)).Bind(
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
