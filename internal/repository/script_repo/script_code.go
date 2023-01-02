package script_repo

import (
	"context"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/pkg/utils/httputils"
	entity "github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
)

type ScriptCodeRepo interface {
	Find(ctx context.Context, id int64) (*entity.Code, error)
	Create(ctx context.Context, scriptCode *entity.Code) error
	Update(ctx context.Context, scriptCode *entity.Code) error
	Delete(ctx context.Context, id int64) error

	FindByVersion(ctx context.Context, scriptId int64, version string, withcode bool) (*entity.Code, error)
	FindLatest(ctx context.Context, scriptId int64, withcode bool) (*entity.Code, error)
	List(ctx context.Context, id int64, request httputils.PageRequest) ([]*entity.Code, int64, error)
}

var defaultScriptCode ScriptCodeRepo

func ScriptCode() ScriptCodeRepo {
	return defaultScriptCode
}

func RegisterScriptCode(i ScriptCodeRepo) {
	defaultScriptCode = i
}

type scriptCodeRepo struct {
}

func NewScriptCodeRepo() ScriptCodeRepo {
	return &scriptCodeRepo{}
}

func (u *scriptCodeRepo) Find(ctx context.Context, id int64) (*entity.Code, error) {
	ret := &entity.Code{ID: id}
	if err := db.Ctx(ctx).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *scriptCodeRepo) Create(ctx context.Context, scriptCode *entity.Code) error {
	return db.Ctx(ctx).Create(scriptCode).Error
}

func (u *scriptCodeRepo) Update(ctx context.Context, scriptCode *entity.Code) error {
	return db.Ctx(ctx).Updates(scriptCode).Error
}

func (u *scriptCodeRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Delete(&entity.Code{ID: id}).Error
}

func (u *scriptCodeRepo) FindByVersion(ctx context.Context, scriptId int64, version string, withcode bool) (*entity.Code, error) {
	ret := &entity.Code{}
	q := db.Ctx(ctx)
	// 由于code过大,使用此方法不返回code
	if !withcode {
		q = q.Select(ret.Fields())
	}
	if err := q.First(ret, "script_id=? and version=?", scriptId, version).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *scriptCodeRepo) FindLatest(ctx context.Context, scriptId int64, withcode bool) (*entity.Code, error) {
	ret := &entity.Code{}
	q := db.Ctx(ctx)
	if !withcode {
		q = q.Select(ret.Fields())
	}
	if err := q.Order("createtime desc").First(ret, "script_id=?", scriptId).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *scriptCodeRepo) List(ctx context.Context, id int64, request httputils.PageRequest) ([]*entity.Code, int64, error) {
	list := make([]*entity.Code, 0)
	q := db.Ctx(ctx).Where("script_id=?", id)
	q = q.Select((&entity.Code{}).Fields())
	if err := q.Order("createtime desc").Offset(request.GetOffset()).Limit(request.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
