package script

import (
	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
)

// CategoryListRequest 脚本分类列表
type CategoryListRequest struct {
	mux.Meta `path:"/script/category" method:"GET"`
	Prefix   string                           `query:"prefix"` // 前缀
	Type     script_entity.ScriptCategoryType `query:"type"`   // 分类类型
}

type CategoryListResponse struct {
	Categories []*CategoryListItem `json:"categories"` // 分类列表
}
