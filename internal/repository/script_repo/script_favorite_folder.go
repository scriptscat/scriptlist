package script_repo

import (
	"context"
	entity "github.com/scriptscat/scriptlist/internal/model/entity/script_entity"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/utils/httputils"
)

type ScriptFavoriteFolderRepo interface {
	Find(ctx context.Context, id int64) (*entity.ScriptFavoriteFolder, error)
	FindPage(ctx context.Context, userId int64, private bool, page httputils.PageRequest) ([]*entity.ScriptFavoriteFolder, int64, error)
	Create(ctx context.Context, scriptFavoriteFolder *entity.ScriptFavoriteFolder) error
	Update(ctx context.Context, scriptFavoriteFolder *entity.ScriptFavoriteFolder) error
	Delete(ctx context.Context, id, userId int64) error
}

var defaultScriptFavoriteFolder ScriptFavoriteFolderRepo

func ScriptFavoriteFolder() ScriptFavoriteFolderRepo {
	return defaultScriptFavoriteFolder
}

func RegisterScriptFavoriteFolder(i ScriptFavoriteFolderRepo) {
	defaultScriptFavoriteFolder = i
}

type scriptFavoriteFolderRepo struct {
}

func NewScriptFavoriteFolder() ScriptFavoriteFolderRepo {
	return &scriptFavoriteFolderRepo{}
}

func (u *scriptFavoriteFolderRepo) Find(ctx context.Context, id int64) (*entity.ScriptFavoriteFolder, error) {
	ret := &entity.ScriptFavoriteFolder{}
	if err := db.Ctx(ctx).Where("id=? and status=?", id, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *scriptFavoriteFolderRepo) Create(ctx context.Context, scriptFavoriteFolder *entity.ScriptFavoriteFolder) error {
	return db.Ctx(ctx).Create(scriptFavoriteFolder).Error
}

func (u *scriptFavoriteFolderRepo) Update(ctx context.Context, scriptFavoriteFolder *entity.ScriptFavoriteFolder) error {
	return db.Ctx(ctx).Updates(scriptFavoriteFolder).Error
}

func (u *scriptFavoriteFolderRepo) Delete(ctx context.Context, id, userId int64) error {
	return db.Ctx(ctx).Model(&entity.ScriptFavoriteFolder{}).Where("id=? and user_id=?", id, userId).Update("status", consts.DELETE).Error
}

func (u *scriptFavoriteFolderRepo) FindPage(ctx context.Context, userId int64, private bool, page httputils.PageRequest) ([]*entity.ScriptFavoriteFolder, int64, error) {
	var list []*entity.ScriptFavoriteFolder
	var count int64
	find := db.Ctx(ctx).Model(&entity.ScriptFavoriteFolder{}).Where("status=?", consts.ACTIVE)
	find = find.Where("user_id=?", userId)
	if !private {
		find = find.Where("private=?", 0)
	}
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := find.Order("createtime desc").Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}
