package api

import (
	"context"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptlist/internal/infrastructure/middleware/token"
	"github.com/scriptscat/scriptlist/internal/pkg/errs"
	"github.com/scriptscat/scriptlist/internal/service/resource/service"
	repository2 "github.com/scriptscat/scriptlist/internal/service/safe/domain/repository"
	service2 "github.com/scriptscat/scriptlist/internal/service/safe/service"
)

type Resource struct {
	svc  service.Resource
	rate service2.Rate
}

func NewResource(svc service.Resource, rate service2.Rate) *Resource {
	return &Resource{
		svc:  svc,
		rate: rate,
	}
}

func (s *Resource) Registry(ctx context.Context, r *gin.Engine) {
	rg := r.Group("/api/v1/resource/image")
	rg.GET("/:id", s.getImg)
	rg.POST("", token.UserAuth(true), s.uploadImg)
}

func (s *Resource) getImg(c *gin.Context) {
	id := c.Param("id")
	res, b, err := s.svc.ReadResource(id)
	if err != nil {
		handelResp(c, err)
		return
	}
	c.Writer.Header().Set("Content-Type", res.ContentType)
	c.Writer.Write(b)
}

func (s *Resource) uploadImg(ctx *gin.Context) {
	handle(ctx, func() interface{} {
		uid, ok := token.UserId(ctx)
		if !ok {
			return errs.ErrNotLogin
		}
		comment := ctx.PostForm("comment")
		if comment != "script" {
			return errs.NewBadRequestError(1000, "图片注释不能为空或'script'其它内容")
		}
		img, err := ctx.FormFile("image")
		if err != nil {
			return err
		}
		f, err := img.Open()
		if err != nil {
			return err
		}
		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		if len(b) > 1048576 {
			return errs.NewBadRequestError(1001, "上传图片不能超过1M")
		}
		var ret gin.H
		if err := s.rate.Rate(&repository2.RateUserInfo{Uid: uid}, &repository2.RateRule{
			Name:     "upload-img",
			Interval: 0,
			DayMax:   40,
		}, func() error {
			res, err := s.svc.UploadScriptImage(uid, comment, img.Filename, b)
			if err != nil {
				return err
			}
			ret = gin.H{
				"id": res.ID,
			}
			return nil
		}); err != nil {
			return err
		}
		return ret
	})
}
