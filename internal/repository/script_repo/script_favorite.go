package script_repo

import (
	"context"

	entity "github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"gorm.io/gorm"

	"github.com/cago-frame/cago/database/db"
	"github.com/cago-frame/cago/pkg/consts"
	"github.com/cago-frame/cago/pkg/utils/httputils"
)

type ScriptFavoriteRepo interface {
	Find(ctx context.Context, id int64) (*entity.ScriptFavorite, error)
	FindPage(ctx context.Context, page httputils.PageRequest) ([]*entity.ScriptFavorite, int64, error)
	Create(ctx context.Context, scriptFavorite *entity.ScriptFavorite) error
	Update(ctx context.Context, scriptFavorite *entity.ScriptFavorite) error
	Delete(ctx context.Context, id int64) error

	FindByFavoriteAndScriptID(ctx context.Context, userId, folderID, scriptID int64) (*entity.ScriptFavorite, error)
	FindByUserIDAndScriptID(ctx context.Context, userId, scriptID int64) ([]*entity.ScriptFavorite, error)
	CountUniqueUsersByScriptID(ctx context.Context, scriptId int64) (int64, error)
}

var defaultScriptFavorite ScriptFavoriteRepo

func ScriptFavorite() ScriptFavoriteRepo {
	return defaultScriptFavorite
}

func RegisterScriptFavorite(i ScriptFavoriteRepo) {
	defaultScriptFavorite = i
}

type scriptFavoriteRepo struct {
}

func NewScriptFavorite() ScriptFavoriteRepo {
	return &scriptFavoriteRepo{}
}

func (u *scriptFavoriteRepo) Find(ctx context.Context, id int64) (*entity.ScriptFavorite, error) {
	ret := &entity.ScriptFavorite{}
	if err := db.Ctx(ctx).Where("id=? and status=?", id, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *scriptFavoriteRepo) Create(ctx context.Context, scriptFavorite *entity.ScriptFavorite) error {
	if err := u.updateFolderCount(ctx, scriptFavorite); err != nil {
		return err
	}
	return db.Ctx(ctx).Create(scriptFavorite).Error
}

func (u *scriptFavoriteRepo) updateFolderCount(ctx context.Context, scriptFavorite *entity.ScriptFavorite) error {
	// 根据状态递增或递减收藏夹中的count
	update := db.Ctx(ctx).Model(&entity.ScriptFavoriteFolder{}).Where("id=?", scriptFavorite.FavoriteFolderID)
	if scriptFavorite.Status == consts.ACTIVE {
		update = update.UpdateColumn("count", gorm.Expr("count + ?", 1))
	} else {
		update = update.UpdateColumn("count", gorm.Expr("count - ?", 1))
	}
	return update.Error
}

func (u *scriptFavoriteRepo) Update(ctx context.Context, scriptFavorite *entity.ScriptFavorite) error {
	if err := u.updateFolderCount(ctx, scriptFavorite); err != nil {
		return err
	}
	return db.Ctx(ctx).Updates(scriptFavorite).Error
}

func (u *scriptFavoriteRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Model(&entity.ScriptFavorite{}).Where("id=?", id).Update("status", consts.DELETE).Error
}

func (u *scriptFavoriteRepo) FindPage(ctx context.Context, page httputils.PageRequest) ([]*entity.ScriptFavorite, int64, error) {
	var list []*entity.ScriptFavorite
	var count int64
	find := db.Ctx(ctx).Model(&entity.ScriptFavorite{}).Where("status=?", consts.ACTIVE)
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := find.Order("createtime desc").Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (u *scriptFavoriteRepo) FindByFavoriteAndScriptID(ctx context.Context, userId, folderID, scriptID int64) (*entity.ScriptFavorite, error) {
	ret := &entity.ScriptFavorite{}
	if err := db.Ctx(ctx).Where("user_id=? and favorite_folder_id=? and script_id=?", userId, folderID, scriptID).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *scriptFavoriteRepo) FindByUserIDAndScriptID(ctx context.Context, userId, scriptID int64) ([]*entity.ScriptFavorite, error) {
	ret := make([]*entity.ScriptFavorite, 0)
	if err := db.Ctx(ctx).Where("user_id=? and script_id=? and status=?", userId, scriptID, consts.ACTIVE).Find(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (u *scriptFavoriteRepo) CountUniqueUsersByScriptID(ctx context.Context, scriptId int64) (int64, error) {
	var count int64
	if err := db.Ctx(ctx).Model(&entity.ScriptFavorite{}).
		Select("COUNT(DISTINCT user_id)").
		Where("script_id=? and status=?", scriptId, consts.ACTIVE).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
