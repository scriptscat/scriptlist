package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/scriptscat/scriptlist/internal/pkg/db"
)

type rate struct {
}

func NewRate() Rate {
	return &rate{}
}

func (r *rate) GetLastOpTime(user, op string) (int64, error) {
	ret, err := db.Redis.Get(context.Background(), r.key(user, op)).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return ret, nil
}

func (r *rate) SetLastOpTime(user, op string, t int64) error {
	err := db.Redis.Set(context.Background(), r.key(user, op), t, time.Hour).Err()
	if err != nil {
		return err
	}
	k := r.key(user, op) + ":" + time.Now().Format("20060102")
	if err := db.Redis.Incr(context.Background(), k).Err(); err != nil {
		return err
	}
	if err := db.Redis.Expire(context.Background(), k, time.Hour*48).Err(); err != nil {
		return err
	}
	return nil
}

func (r *rate) GetDayOpCnt(user, op string) (int64, error) {
	k := r.key(user, op) + ":" + time.Now().Format("20060102")
	ret, err := db.Redis.Get(context.Background(), k).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return ret, nil
}

func (r *rate) key(user, op string) string {
	return fmt.Sprintf("safe:rate:" + user + ":" + op)
}
