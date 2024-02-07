package open

import "github.com/codfrm/cago/server/mux"

// CrxDownloadRequest 谷歌crx下载服务
type CrxDownloadRequest struct {
	mux.Meta `path:"/open/crx-download/:id" method:"GET"`
	ID       int64 `uri:"id"`
}

type CrxDownloadResponse struct {
}
