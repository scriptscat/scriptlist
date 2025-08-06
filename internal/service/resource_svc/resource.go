package resource_svc

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/utils"
	api "github.com/scriptscat/scriptlist/internal/api/resource"
	"github.com/scriptscat/scriptlist/internal/model/entity/resource_entity"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
	"github.com/scriptscat/scriptlist/internal/repository/resource_repo"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
)

type ResourceSvc interface {
	// UploadImage 上传图片
	UploadImage(ctx context.Context, image *multipart.FileHeader, req *api.UploadImageRequest) (*api.UploadImageResponse, error)
	// ViewImage 查看图片
	ViewImage(ctx context.Context, id string) (*resource_entity.Resource, []byte, error)
}

type resourceSvc struct {
	dir string
}

var defaultResource = &resourceSvc{
	dir: "./resource/",
}

func Resource() ResourceSvc {
	return defaultResource
}

// UploadImage 上传图片
func (r *resourceSvc) UploadImage(ctx context.Context, image *multipart.FileHeader, req *api.UploadImageRequest) (*api.UploadImageResponse, error) {
	base := path.Join(r.dir, "images", time.Now().Format("2006/0102"))
	if err := os.MkdirAll(base, 0750); err != nil {
		return nil, err
	}
	f, err := image.Open()
	if err != nil {
		return nil, err
	}
	defer func(f multipart.File) {
		_ = f.Close()
	}(f)
	bImage := make([]byte, image.Size)
	if _, err := f.Read(bImage); err != nil {
		return nil, err
	}
	ct := http.DetectContentType(bImage)
	if !strings.Contains(ct, "image") {
		return nil, i18n.NewError(ctx, code.ResourceNotImage)
	}
	resource := &resource_entity.Resource{
		ResourceID: utils.RandString(16, utils.Mix),
		UserID:     auth_svc.Auth().Get(ctx).UID,
		LinkID:     req.LinkID,
		Comment:    req.Comment,
		Name:       image.Filename,
		Path: path.Join(base,
			fmt.Sprintf("%d%s_%d%s", time.Now().Unix(), utils.RandString(3, utils.Number),
				auth_svc.Auth().Get(ctx).UID, path.Ext(image.Filename))),
		ContentType: ct,
		Status:      consts.ACTIVE,
		Createtime:  time.Now().Unix(),
	}
	if err := os.WriteFile(resource.Path, bImage, os.ModePerm); err != nil {
		return nil, err
	}
	if err := resource_repo.Resource().Create(ctx, resource); err != nil {
		return nil, err
	}
	return &api.UploadImageResponse{
		ID:          resource.ResourceID,
		LinkID:      resource.LinkID,
		Comment:     resource.Comment,
		Name:        resource.Name,
		ContentType: resource.ContentType,
		Createtime:  resource.Createtime,
	}, nil
}

// ViewImage 查看图片
func (r *resourceSvc) ViewImage(ctx context.Context, id string) (*resource_entity.Resource, []byte, error) {
	res, err := resource_repo.Resource().FindByResourceID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if res == nil {
		return nil, nil, i18n.NewErrorWithStatus(ctx, http.StatusNotFound, code.ResourceNotFound)
	}
	b, err := os.ReadFile(res.Path)
	if err != nil {
		return nil, nil, err
	}
	return res, b, nil
}
