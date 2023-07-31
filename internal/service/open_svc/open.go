package open_svc

import (
	"context"

	api "github.com/scriptscat/scriptlist/internal/api/open"
)

type OpenSvc interface {
	// CrxDownload 谷歌crx下载服务
	CrxDownload(ctx context.Context, req *api.CrxDownloadRequest) (*api.CrxDownloadResponse, error)
}

type openSvc struct {
}

var defaultOpen = &openSvc{}

func Open() OpenSvc {
	return defaultOpen
}

// CrxDownload 谷歌crx下载服务
func (o *openSvc) CrxDownload(ctx context.Context, req *api.CrxDownloadRequest) (*api.CrxDownloadResponse, error) {
	return nil, nil
}
