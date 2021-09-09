package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/scriptscat/scriptweb/internal/domain/resource/entity"
	"github.com/scriptscat/scriptweb/internal/pkg/db"
)

type resource struct {
}

func NewResource() Resource {
	return &resource{}
}

func (r *resource) Save(res *entity.Resource) error {
	return db.Redis.Set(context.Background(), r.key(res.ID), res, 0).Err()
}

func (r *resource) Find(id string) (*entity.Resource, error) {
	ret, err := db.Redis.Get(context.Background(), r.key(id)).Result()
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
