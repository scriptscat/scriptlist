package script_repo

import (
	"context"
	"strconv"

	"github.com/codfrm/cago/database/redis"
	goRedis "github.com/go-redis/redis/v9"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
)

type ScriptWatchRepo interface {
	List(ctx context.Context, script int64) ([]*script_entity.Watch, error)
	Num(ctx context.Context, script int64) (int, error)
	Watch(ctx context.Context, script, user int64, level script_entity.ScriptWatchLevel) error
	Unwatch(ctx context.Context, script, user int64) error
	IsWatch(ctx context.Context, script, user int64) (script_entity.ScriptWatchLevel, error)
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

func (w *scriptWatchRepo) key(issue int64) string {
	return "script:watch:" + strconv.FormatInt(issue, 10)
}

func (w *scriptWatchRepo) List(ctx context.Context, script int64) ([]*script_entity.Watch, error) {
	list, err := redis.Ctx(ctx).HGetAll(w.key(script)).Result()
	if err != nil {
		return nil, err
	}
	ret := make([]*script_entity.Watch, 0)
	for k, v := range list {
		uid, _ := strconv.ParseInt(k, 10, 64)
		level, _ := strconv.Atoi(v)
		ret = append(ret, &script_entity.Watch{UserID: uid, Level: script_entity.ScriptWatchLevel(level)})
	}
	return ret, nil
}

func (w *scriptWatchRepo) Num(ctx context.Context, script int64) (int, error) {
	list, err := redis.Ctx(ctx).HGetAll(w.key(script)).Result()
	if err != nil {
		return 0, err
	}
	return len(list), nil
}

func (w *scriptWatchRepo) Watch(ctx context.Context, script, user int64, level script_entity.ScriptWatchLevel) error {
	return redis.Ctx(ctx).HSet(w.key(script), user, int(level)).Err()
}

func (w *scriptWatchRepo) Unwatch(ctx context.Context, script, user int64) error {
	return redis.Ctx(ctx).HDel(w.key(script), strconv.FormatInt(user, 10)).Err()
}

func (w *scriptWatchRepo) IsWatch(ctx context.Context, script, user int64) (script_entity.ScriptWatchLevel, error) {
	ret, err := redis.Ctx(ctx).HGet(w.key(script), strconv.FormatInt(user, 10)).Int()
	if err != nil {
		if err == goRedis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return script_entity.ScriptWatchLevel(ret), nil
}
