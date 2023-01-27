package script_repo

import (
	"context"
	"time"

	"github.com/codfrm/cago/database/db"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
)

type ScriptWatchRepo interface {
	Find(ctx context.Context, id int64) (*script_entity.ScriptWatch, error)
	Create(ctx context.Context, scriptWatch *script_entity.ScriptWatch) error
	Update(ctx context.Context, scriptWatch *script_entity.ScriptWatch) error
	Delete(ctx context.Context, id int64) error

	// FindAll 查询出所有符合条件的记录
	FindAll(ctx context.Context, script int64, level script_entity.ScriptWatchLevel) ([]*script_entity.ScriptWatch, error)
	FindByUser(ctx context.Context, script, user int64) (*script_entity.ScriptWatch, error)
	Watch(ctx context.Context, script, user int64, level script_entity.ScriptWatchLevel) error
}

var defaultScriptWatch ScriptWatchRepo

func ScriptWatch() ScriptWatchRepo {
	return defaultScriptWatch
}

func RegisterScriptWatch(i ScriptWatchRepo) {
	defaultScriptWatch = i
}

type scriptWatchRepo struct {
}

func NewScriptWatchRepo() ScriptWatchRepo {
	return &scriptWatchRepo{}
}

func (s *scriptWatchRepo) Find(ctx context.Context, id int64) (*script_entity.ScriptWatch, error) {
	ret := &script_entity.ScriptWatch{ID: id}
	if err := db.Ctx(ctx).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (s *scriptWatchRepo) Create(ctx context.Context, scriptWatch *script_entity.ScriptWatch) error {
	return db.Ctx(ctx).Create(scriptWatch).Error
}

func (s *scriptWatchRepo) Update(ctx context.Context, scriptWatch *script_entity.ScriptWatch) error {
	return db.Ctx(ctx).Updates(scriptWatch).Error
}

func (s *scriptWatchRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Delete(&script_entity.ScriptWatch{ID: id}).Error
}

func (s *scriptWatchRepo) FindAll(ctx context.Context, scriptId int64, level script_entity.ScriptWatchLevel) ([]*script_entity.ScriptWatch, error) {
	var list []*script_entity.ScriptWatch
	if err := db.Ctx(ctx).Where("script_id=? and level>=?", scriptId, level).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *scriptWatchRepo) Watch(ctx context.Context, script, user int64, level script_entity.ScriptWatchLevel) error {
	watch, err := s.FindByUser(ctx, script, user)
	if err != nil {
		return err
	}
	if watch == nil {
		// 新建
		watch = &script_entity.ScriptWatch{
			UserID:     user,
			ScriptID:   script,
			Level:      level,
			Createtime: time.Now().Unix(),
		}
	} else {
		// 更新
		watch.Level = level
		watch.Updatetime = time.Now().Unix()
	}
	return db.Ctx(ctx).Save(watch).Error
}

func (s *scriptWatchRepo) FindByUser(ctx context.Context, script, user int64) (*script_entity.ScriptWatch, error) {
	ret := &script_entity.ScriptWatch{}
	if err := db.Ctx(ctx).Where("script_id=? and user_id=?", script, user).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}
