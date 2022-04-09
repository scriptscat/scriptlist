package persistence

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/scriptscat/scriptlist/internal/service/resource/domain/entity"
	"github.com/scriptscat/scriptlist/internal/service/resource/domain/repository"
)

type resource struct {
	redis *redis.Client
}

func NewResource(redis *redis.Client) repository.Resource {
	return &resource{
		redis: redis,
	}
}

func (r *resource) Save(res *entity.Resource) error {
	return r.redis.Set(context.Background(), r.key(res.ID), res, 0).Err()
}

func (r *resource) Find(id string) (*entity.Resource, error) {
	ret, err := r.redis.Get(context.Background(), r.key(id)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	res := &entity.Resource{}
	if err := json.Unmarshal([]byte(ret), res); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *resource) key(id string) string {
	return fmt.Sprintf("resource:id:%s", id)
}
