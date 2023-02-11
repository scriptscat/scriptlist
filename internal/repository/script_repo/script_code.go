package script_repo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/codfrm/cago/database/cache"
	cache2 "github.com/codfrm/cago/database/cache/cache"
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
	FindLatest(ctx context.Context, scriptId int64, offset int, withcode bool) (*entity.Code, error)
	FindPreLatest(ctx context.Context, scriptId int64, offset int, withcode bool) (*entity.Code, error)
	FindAllLatest(ctx context.Context, scriptId int64, offset int, withcode bool) (*entity.Code, error)
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

func (u *scriptCodeRepo) key(id int64) string {
	return "script:code:" + strconv.FormatInt(id, 10)
}

func (u *scriptCodeRepo) KeyDepend(id int64) *cache2.KeyDepend {
	return cache2.NewKeyDepend(cache.Default(), u.key(id)+":dep")
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
	if err := db.Ctx(ctx).Create(scriptCode).Error; err != nil {
		return err
	}
	return u.KeyDepend(scriptCode.ScriptID).InvalidKey(ctx)
}

func (u *scriptCodeRepo) Update(ctx context.Context, scriptCode *entity.Code) error {
	if err := db.Ctx(ctx).Updates(scriptCode).Error; err != nil {
		return err
	}
	return u.KeyDepend(scriptCode.ScriptID).InvalidKey(ctx)
}

func (u *scriptCodeRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Delete(&entity.Code{ID: id}).Error
}

func (u *scriptCodeRepo) FindByVersion(ctx context.Context, scriptId int64, version string, withcode bool) (*entity.Code, error) {
	ret := &entity.Code{}
	if err := cache.Ctx(ctx).GetOrSet(u.key(scriptId)+fmt.Sprintf(":%s:%v", version, withcode), func() (interface{}, error) {
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
	}, cache2.Expiration(time.Hour), cache2.WithDepend(u.KeyDepend(scriptId))).Scan(&ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (u *scriptCodeRepo) FindLatest(ctx context.Context, scriptId int64, offset int, withcode bool) (*entity.Code, error) {
	ret := &entity.Code{}
	if err := cache.Ctx(ctx).GetOrSet(u.key(scriptId)+fmt.Sprintf(":%d:%v", offset, withcode), func() (interface{}, error) {
		q := db.Ctx(ctx)
		if !withcode {
			q = q.Select(ret.Fields())
		}
		if err := q.Order("createtime desc").Offset(offset).
			First(ret, "script_id=? and is_pre_release=?",
				scriptId, entity.DisablePreReleaseScript).Error; err != nil {
			if db.RecordNotFound(err) {
				return nil, nil
			}
			return nil, err
		}
		return ret, nil
	}, cache2.Expiration(time.Hour), cache2.WithDepend(u.KeyDepend(scriptId))).Scan(&ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (u *scriptCodeRepo) FindPreLatest(ctx context.Context, scriptId int64, offset int, withcode bool) (*entity.Code, error) {
	ret := &entity.Code{}
	if err := cache.Ctx(ctx).GetOrSet(u.key(scriptId)+fmt.Sprintf(":pre:%d:%v", offset, withcode), func() (interface{}, error) {
		q := db.Ctx(ctx)
		if !withcode {
			q = q.Select(ret.Fields())
		}
		if err := q.Order("createtime desc").Offset(offset).
			First(ret, "script_id=? and is_pre_release=?",
				scriptId, entity.EnablePreReleaseScript).Error; err != nil {
			if db.RecordNotFound(err) {
				return nil, nil
			}
			return nil, err
		}
		return ret, nil
	}, cache2.Expiration(time.Hour), cache2.WithDepend(u.KeyDepend(scriptId))).Scan(&ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (u *scriptCodeRepo) FindAllLatest(ctx context.Context, scriptId int64, offset int, withcode bool) (*entity.Code, error) {
	ret := &entity.Code{}
	if err := cache.Ctx(ctx).GetOrSet(u.key(scriptId)+fmt.Sprintf(":all:%d:%v", offset, withcode), func() (interface{}, error) {
		q := db.Ctx(ctx)
		if !withcode {
			q = q.Select(ret.Fields())
		}
		if err := q.Order("createtime desc").Offset(offset).
			First(ret, "script_id=?", scriptId).Error; err != nil {
			if db.RecordNotFound(err) {
				return nil, nil
			}
			return nil, err
		}
		return ret, nil
	}, cache2.Expiration(time.Hour), cache2.WithDepend(u.KeyDepend(scriptId))).Scan(&ret); err != nil {
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
