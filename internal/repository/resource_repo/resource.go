package resource_repo

import (
	"context"

	"github.com/cago-frame/cago/database/db"
	"github.com/cago-frame/cago/pkg/consts"
	"github.com/cago-frame/cago/pkg/utils/httputils"
	"github.com/scriptscat/scriptlist/internal/model/entity/resource_entity"
)

type ResourceRepo interface {
	Find(ctx context.Context, id int64) (*resource_entity.Resource, error)
	FindPage(ctx context.Context, page httputils.PageRequest) ([]*resource_entity.Resource, int64, error)
	Create(ctx context.Context, resource *resource_entity.Resource) error
	Update(ctx context.Context, resource *resource_entity.Resource) error
	Delete(ctx context.Context, id int64) error
	// FindByResourceID 根据资源id获取资源
	FindByResourceID(ctx context.Context, id string) (*resource_entity.Resource, error)
}

var defaultResource ResourceRepo

func Resource() ResourceRepo {
	return defaultResource
}

func RegisterResource(i ResourceRepo) {
	defaultResource = i
}

type resourceRepo struct {
}

func NewResource() ResourceRepo {
	return &resourceRepo{}
}

func (u *resourceRepo) Find(ctx context.Context, id int64) (*resource_entity.Resource, error) {
	ret := &resource_entity.Resource{ID: id}
	if err := db.Ctx(ctx).Where("status=?", consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *resourceRepo) Create(ctx context.Context, resource *resource_entity.Resource) error {
	return db.Ctx(ctx).Create(resource).Error
}

func (u *resourceRepo) Update(ctx context.Context, resource *resource_entity.Resource) error {
	return db.Ctx(ctx).Updates(resource).Error
}

func (u *resourceRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Model(&resource_entity.Resource{ID: id}).Update("status", consts.DELETE).Error
}

func (u *resourceRepo) FindPage(ctx context.Context, page httputils.PageRequest) ([]*resource_entity.Resource, int64, error) {
	var list []*resource_entity.Resource
	var count int64
	if err := db.Ctx(ctx).Model(&resource_entity.Resource{}).Where("status=?", consts.ACTIVE).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Ctx(ctx).Where("status=?", consts.ACTIVE).Order("createtime desc").Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (u *resourceRepo) FindByResourceID(ctx context.Context, id string) (*resource_entity.Resource, error) {
	ret := &resource_entity.Resource{}
	if err := db.Ctx(ctx).Where("resource_id=? and status=?", id, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}
