package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/scriptscat/scriptlist/internal/service/safe/domain/repository"
)

type rate struct {
	redis *redis.Client
}

func NewRate(redis *redis.Client) repository.Rate {
	return &rate{
		redis: redis,
	}
}

func (r *rate) GetLastOpTime(user, op string) (int64, error) {
	ret, err := r.redis.Get(context.Background(), r.key(user, op)).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return ret, nil
}

func (r *rate) SetLastOpTime(user, op string, t int64) error {
	err := r.redis.Set(context.Background(), r.key(user, op), t, time.Hour).Err()
	if err != nil {
		return err
	}
	k := r.key(user, op) + ":" + time.Now().Format("20060102")
	if err := r.redis.Incr(context.Background(), k).Err(); err != nil {
		return err
	}
	if err := r.redis.Expire(context.Background(), k, time.Hour*48).Err(); err != nil {
		return err
	}
	return nil
}

func (r *rate) GetDayOpCnt(user, op string) (int64, error) {
	k := r.key(user, op) + ":" + time.Now().Format("20060102")
	ret, err := r.redis.Get(context.Background(), k).Int64()
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
