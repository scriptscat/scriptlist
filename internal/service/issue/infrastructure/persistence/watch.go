package persistence

import (
	"context"
	"strconv"

	goRedis "github.com/go-redis/redis/v8"
	"github.com/scriptscat/scriptlist/internal/service/issue/domain/repository"
	"github.com/scriptscat/scriptlist/pkg/utils"
)

type watch struct {
	redis *goRedis.Client
}

func NewWatch(redis *goRedis.Client) repository.IssueWatch {
	return &watch{redis: redis}
}

func (w *watch) key(issue int64) string {
	return "script:issue:watch:" + strconv.FormatInt(issue, 10)
}

func (w *watch) List(issue int64) ([]*repository.Watch, error) {
	list, err := w.redis.HGetAll(context.Background(), w.key(issue)).Result()
	if err != nil {
		return nil, err
	}
	ret := make([]*repository.Watch, 0)
	for k, v := range list {
		if v != "1" {
			continue
		}
		ret = append(ret, &repository.Watch{UserId: utils.StringToInt64(k)})
	}
	return ret, nil
}

func (w *watch) Num(issue int64) (int, error) {
	list, err := w.redis.HGetAll(context.Background(), w.key(issue)).Result()
	if err != nil {
		return 0, err
	}
	return len(list), nil
}

func (w *watch) Watch(issue, user int64) error {
	return w.redis.HSet(context.Background(), w.key(issue), user, "1").Err()
}

func (w *watch) Unwatch(issue, user int64) error {
	return w.redis.HSet(context.Background(), w.key(issue), user, "2").Err()
}

// IsWatch 0 从未关注过 1 关注 2 关注但是取消了
func (w *watch) IsWatch(issue, user int64) (int, error) {
	ret, err := w.redis.HGet(context.Background(), w.key(issue), strconv.FormatInt(user, 10)).Result()
	if err != nil {
		if err == goRedis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return utils.StringToInt(ret), nil
}
