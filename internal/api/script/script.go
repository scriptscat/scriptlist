package script

import (
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/cago/server/mux"
)

type Item struct {
	ID int64 `json:"id"`
}

// ListRequest 获取脚本列表
type ListRequest struct {
	mux.Meta                     `path:"/script" method:"GET"`
	httputils.PageRequest[*Item] `form:",inline"`
}

type ListResponse struct {
	httputils.PageResponse[*Item] `json:",inline"`
}

// CreateRequest 创建脚本
type CreateRequest struct {
	mux.Meta `path:"/script" method:"POST"`
}

type CreateResponse struct {
	ID int64 `json:"id"`
}
