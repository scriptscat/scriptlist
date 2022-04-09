package service

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/scriptscat/scriptlist/internal/pkg/errs"
	"github.com/scriptscat/scriptlist/internal/service/resource/domain/entity"
	"github.com/scriptscat/scriptlist/internal/service/resource/domain/repository"
	"github.com/scriptscat/scriptlist/pkg/utils"
)

type Resource interface {
	UploadScriptImage(uid int64, comment string, name string, image []byte) (*entity.Resource, error)
	ReadResource(id string) (*entity.Resource, []byte, error)
}

type resource struct {
	dir  string
	repo repository.Resource
}

func NewResource(repo repository.Resource) Resource {
	return &resource{
		dir:  "./resource/",
		repo: repo,
	}
}

func (r *resource) UploadScriptImage(uid int64, comment string, name string, image []byte) (*entity.Resource, error) {
	base := path.Join(r.dir, time.Now().Format("2006/0102"))
	if err := os.MkdirAll(base, 0755); err != nil {
		return nil, err
	}
	ct := http.DetectContentType(image)
	if strings.Index(ct, "image") == -1 {
		return nil, errs.ErrResourceNotImage
	}
	resource := &entity.Resource{
		ID:          utils.RandString(16, 2),
		Uid:         uid,
		Comment:     comment,
		Name:        name,
		ContentType: ct,
		Path:        path.Join(base, fmt.Sprintf("%d_%d.%s", time.Now().Unix(), uid, path.Ext(name))),
		Createtime:  time.Now().Unix(),
	}
	if err := os.WriteFile(resource.Path, image, 0644); err != nil {
		return nil, err
	}
	if err := r.repo.Save(resource); err != nil {
		return nil, err
	}
	return resource, nil
}

func (r *resource) ReadResource(id string) (*entity.Resource, []byte, error) {
	res, err := r.repo.Find(id)
	if err != nil {
		return nil, nil, err
	}
	if res == nil {
		return nil, nil, errs.ErrResourceNotFound
	}
	b, err := os.ReadFile(res.Path)
	if err != nil {
		return nil, nil, err
	}
	return res, b, nil
}
