package api

import (
	"github.com/codfrm/cago/server/http"
	_ "github.com/scriptscat/scriptlist/docs"
)

// Router 路由表
// @title    脚本站 API 文档
// @version  2.0.0
// @BasePath /api/v2
func Router(r *http.Router) error {
	r = r.Group("/api/v2")

	return nil
}
