package migrations

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/codfrm/cago/database/redis"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/scriptlist/internal/model/entity/resource_entity"
	"github.com/scriptscat/scriptlist/internal/repository/resource_repo"
	"gorm.io/gorm"
)

func T1674893758() *gormigrate.Migration {
	// Resource 老版本存在redis中的数据结构
	type Resource struct {
		ID          string `json:"id"`
		Uid         int64  `json:"uid"`
		Comment     string `json:"comment"`
		Name        string `json:"name"`
		Path        string `json:"path"`
		ContentType string `json:"content_type"`
		Createtime  int64  `json:"createtime"`
	}
	return &gormigrate.Migration{
		ID: "1674893758",
		Migrate: func(tx *gorm.DB) error {
			ctx := context.Background()
			if err := tx.AutoMigrate(&resource_entity.Resource{}); err != nil {
				return err
			}
			// 迁移resource数据
			return ScanKeys(ctx, "resource:id:*", func(ctx context.Context, key string) error {
				resourceId := strings.TrimPrefix(key, "resource:id:")
				if resourceId == "" {
					return errors.New("解析resourceId失败")
				}
				ret, err := redis.Ctx(ctx).Get(key).Result()
				if err != nil {
					return err
				}
				res := &Resource{}
				if err := json.Unmarshal([]byte(ret), res); err != nil {
					return err
				}
				m := &resource_entity.Resource{
					ResourceID:  res.ID,
					UserID:      res.Uid,
					LinkID:      -1,
					Comment:     res.Comment,
					Name:        res.Name,
					Path:        res.Path,
					Status:      consts.ACTIVE,
					ContentType: res.ContentType,
					Createtime:  res.Createtime,
				}
				if err := resource_repo.Resource().Create(ctx, m); err != nil {
					return err
				}
				return nil
			})
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropTable(&resource_entity.Resource{}); err != nil {
				return err
			}
			return nil
		},
	}
}
