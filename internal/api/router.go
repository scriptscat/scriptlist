package api

import (
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/server/mux"
	_ "github.com/scriptscat/scriptlist/docs"
	"github.com/scriptscat/scriptlist/internal/controller/auth_ctr"
	"github.com/scriptscat/scriptlist/internal/controller/script_ctr"
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
		r.GET("/user/avatar/:uid", controller.Avatar())
		r.Group("/").Bind(
			controller.Info, // 获取用户信息
		)
		r.Group("/", auth.Middleware(false)).Bind(
			controller.Script,
			controller.GetFollow,
		)
	}
	// 脚本
	{
		controller := script_ctr.NewScript()
		// 需要用户登录的路由组
		r.Group("/", auth.Middleware(true)).Bind(
			controller.Create,
			controller.UpdateCode,
			controller.MigrateEs,
			controller.Watch,
			controller.GetSetting,
			controller.UpdateSetting,
			controller.Archive,
			controller.Delete,
		)
		// 处理下载
		root.GET("/scripts/code/:id/*name", auth.Middleware(false), controller.Download())
		// 无需用户登录的路由组
		r.Group("/").Bind(
			controller.List,
			controller.Info,
			controller.Code,
			controller.VersionList,
			controller.VersionCode,
		)
		// 半登录
		r.Group("/", auth.Middleware(false)).Bind(
			controller.State,
		)
	}
	return nil
}
