package resource_ctr

import (
	"strconv"

	"github.com/codfrm/cago/database/redis"
	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/limit"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
	api "github.com/scriptscat/scriptlist/internal/api/resource"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	"github.com/scriptscat/scriptlist/internal/service/resource_svc"
)

type Resource struct {
	limit *limit.PeriodLimit
}

func NewResource() *Resource {
	return &Resource{
		limit: limit.NewPeriodLimit(
			300, 20, redis.Default(), "limit:resource",
		),
	}
}

// UploadImage 上传图片
func (r *Resource) UploadImage(gCtx *gin.Context, req *api.UploadImageRequest) (*api.UploadImageResponse, error) {
	ctx := gCtx
	img, err := gCtx.FormFile("image")
	if err != nil {
		return nil, err
	}
	// 1M限制
	if img.Size > 1024*1024*5 {
		return nil, i18n.NewError(ctx, code.ResourceImageTooLarge)
	}
	resp, err := r.limit.FuncTake(ctx, strconv.FormatInt(auth_svc.Auth().Get(ctx).UID, 10), func() (interface{}, error) {
		return resource_svc.Resource().UploadImage(gCtx, img, req)
	})
	if err != nil {
		return nil, err
	}
	return resp.(*api.UploadImageResponse), nil
}

func (r *Resource) ViewImage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		res, b, err := resource_svc.Resource().ViewImage(ctx, id)
		if err != nil {
			httputils.HandleResp(ctx, err)
			return
		}
		ctx.Writer.Header().Set("Content-Type", res.ContentType)
		_, _ = ctx.Writer.Write(b)
	}
}
