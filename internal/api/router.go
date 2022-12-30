package api

import (
	"github.com/codfrm/cago/configs"
	cache2 "github.com/codfrm/cago/database/cache"
	"github.com/codfrm/cago/middleware/sessions"
	"github.com/codfrm/cago/middleware/sessions/cache"
	"github.com/codfrm/cago/server/mux"
	_ "github.com/scriptscat/scriptlist/docs"
	controller2 "github.com/scriptscat/scriptlist/internal/controller/script_ctr"
	"github.com/scriptscat/scriptlist/internal/controller/user_ctr"
	_ "github.com/scriptscat/scriptlist/internal/repository/persistence"
)

// Router 路由表
// @title    脚本站 API 文档
// @version  2.0.0
// @BasePath /api/v2
func Router(root *mux.Router) error {
	r := root.Group("/api/v2")
	r.Use(sessions.Middleware("SESSION",
		cache.NewCacheStore(cache2.Default(), "session"),
	))
	auth := user_ctr.NewAuth()
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
		)
		r.Group("/").Bind(
			controller.Info, // 获取用户信息
		)
	}
	// 脚本
	{
		controller := controller2.NewScript()
		// 需要用户登录的路由组
		r.Group("/", auth.Middleware(true)).Bind(
			controller.Create,
			controller.UpdateCode,
			controller.MigrateEs,
		)
		// 处理下载
		root.GET("/scripts/code/:id/*name", auth.Middleware(false), controller.Download())
		// 无需用户登录的路由组
		r.Group("/").Bind(
			controller.List,
		)
	}

	return nil
}
