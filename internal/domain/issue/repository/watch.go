package repository

import (
	"context"
	"strconv"

	goRedis "github.com/go-redis/redis/v8"
	"github.com/scriptscat/scriptlist/internal/pkg/db"
	"github.com/scriptscat/scriptlist/pkg/utils"
)

type watch struct {
}

func NewWatch() IssueWatch {
	return &watch{}
}

func (w *watch) key(issue int64) string {
	return "script:issue:watch:" + strconv.FormatInt(issue, 10)
}

func (w *watch) List(issue int64) ([]*Watch, error) {
	list, err := db.Redis.HGetAll(context.Background(), w.key(issue)).Result()
	if err != nil {
		return nil, err
	}
	ret := make([]*Watch, 0)
	for k := range list {
		ret = append(ret, &Watch{UserId: utils.StringToInt64(k)})
	}
	return ret, nil
}

func (w *watch) Num(issue int64) (int, error) {
	list, err := db.Redis.HGetAll(context.Background(), w.key(issue)).Result()
	if err != nil {
		return 0, err
	}
	return len(list), nil
}

func (w *watch) Watch(issue, user int64) error {
	return db.Redis.HSet(context.Background(), w.key(issue), user, "1").Err()
}

func (w *watch) Unwatch(issue, user int64) error {
	return db.Redis.HDel(context.Background(), w.key(issue), strconv.FormatInt(user, 10)).Err()
}

func (w *watch) IsWatch(issue, user int64) (bool, error) {
	ret, err := db.Redis.HGet(context.Background(), w.key(issue), strconv.FormatInt(user, 10)).Result()
	if err != nil {
		if err == goRedis.Nil {
			return false, nil
		}
		return false, err
	}
	return ret == "1", nil
}
