package api

import (
	cache2 "github.com/codfrm/cago/database/cache"
	"github.com/codfrm/cago/middleware/sessions"
	"github.com/codfrm/cago/middleware/sessions/cache"
	"github.com/codfrm/cago/server/mux"
	_ "github.com/scriptscat/scriptlist/docs"
	"github.com/scriptscat/scriptlist/internal/controller/script"
	"github.com/scriptscat/scriptlist/internal/controller/user"
	_ "github.com/scriptscat/scriptlist/internal/repository/persistence"
)

// Router 路由表
// @title    脚本站 API 文档
// @version  2.0.0
// @BasePath /api/v2
func Router(r *mux.Router) error {
	r = r.Group("/api/v2")
	r.Use(sessions.Middleware("SESSION",
		cache.NewCacheStore(cache2.Default(), "session"),
	))
	auth := user.NewAuth()
	// 用户-auth
	rg := r.Group("/")
	{
		if err := rg.Bind(auth); err != nil {
			return err
		}
	}
	// 用户
	{
		controller := user.NewUser()
		authRg := r.Group("/")
		{
			authRg.Use(auth.Middleware(true))
			if err := authRg.Bind(
				controller.CurrentUser,
			); err != nil {
				return err
			}
		}
		rg := r.Group("/")
		{
			if err := rg.Bind(
				controller,
			); err != nil {
				return err
			}
		}
	}
	// 脚本
	{
		controller := script.NewScript()
		authRg := r.Group("/")
		{
			authRg.Use(auth.Middleware(true))
			if err := authRg.Bind(
				controller.Create,
			); err != nil {
				return err
			}
		}
		rg := r.Group("/")
		{
			if err := rg.Bind(
				controller,
			); err != nil {
				return err
			}
		}
	}

	return nil
}
