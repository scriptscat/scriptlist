package open

import "github.com/cago-frame/cago/server/mux"

// CrxDownloadRequest 谷歌crx下载服务
type CrxDownloadRequest struct {
	mux.Meta `path:"/open/crx-download/:id" method:"GET"`
	ID       int64 `uri:"id"`
}

type CrxDownloadResponse struct {
}

// FaviconRequest 获取网站favicon图标
type FaviconRequest struct {
	mux.Meta `path:"/open/favicons" method:"GET"`
	Domain   string `form:"domain" binding:"required"`
	Sz       int    `form:"sz"`
}
