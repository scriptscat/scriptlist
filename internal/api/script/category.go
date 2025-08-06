package script

import (
	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
)

// CategoryListRequest 脚本分类列表
type CategoryListRequest struct {
	mux.Meta `path:"/scripts/category" method:"GET"`
	Prefix   string                           `form:"prefix"`                            // 前缀
	Type     script_entity.ScriptCategoryType `form:"type" binding:"required,oneof=1 2"` // 分类类型: 1: 脚本分类, 2: Tag
}

type CategoryListResponse struct {
	Categories []*CategoryListItem `json:"categories"` // 分类列表
}
