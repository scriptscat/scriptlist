package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/scriptscat/scriptweb/internal/pkg/db"
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
	return db.Redis.Set(context.Background(), r.key(user, op), t, time.Hour).Err()
}

func (r *rate) key(user, op string) string {
	return fmt.Sprintf("safe:rate:" + user)
}
