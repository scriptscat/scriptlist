package script_ctr

import (
	"context"

	api "github.com/scriptscat/scriptlist/internal/api/script"
	"github.com/scriptscat/scriptlist/internal/service/script_svc"
)

type Category struct {
}

func NewCategory() *Category {
	return &Category{}
}

// CategoryList 获取脚本分类列表
func (c *Category) CategoryList(ctx context.Context, req *api.CategoryListRequest) (*api.CategoryListResponse, error) {
	return script_svc.Category().CategoryList(ctx, req)
}
