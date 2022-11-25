package api

import (
	cache2 "github.com/codfrm/cago/database/cache"
	"github.com/codfrm/cago/middleware/sessions"
	"github.com/codfrm/cago/middleware/sessions/cache"
	"github.com/codfrm/cago/server/mux"
	_ "github.com/scriptscat/scriptlist/docs"
	"github.com/scriptscat/scriptlist/internal/controller/user"
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
	{
		rg := r.Group("/")
		if err := rg.Bind(auth); err != nil {
			return err
		}
	}
	{
		rg := r.Group("/")
		rg.Use(auth.Middleware())
		controller := user.NewUser()
		if err := rg.Bind(controller); err != nil {
			return err
		}
	}

	return nil
}
