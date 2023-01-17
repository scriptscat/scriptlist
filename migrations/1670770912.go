package migrations

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/codfrm/cago/database/redis"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// T1670770912 迁移关注、资源
func T1670770912() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1670770912",
		Migrate: func(db *gorm.DB) error {
			ctx := context.Background()
			if err := db.AutoMigrate(&script_entity.ScriptWatch{}); err != nil {
				return err
			}
			if err := ScanKeys(ctx, "script:watch:*", func(ctx context.Context, key string) error {
				// 取出id
				scriptId, err := strconv.ParseInt(key[strings.LastIndex(key, ":")+1:], 10, 64)
				if err != nil {
					return err
				}
				if scriptId == 0 {
					return errors.New("id解析错误")
				}
				list, err := redis.Ctx(ctx).HGetAll(key).Result()
				if err != nil {
					return err
				}
				for k, v := range list {
					uid, _ := strconv.ParseInt(k, 10, 64)
					level, _ := strconv.Atoi(v)
					if err := script_repo.ScriptWatch().Create(ctx, &script_entity.ScriptWatch{
						UserID:     uid,
						ScriptID:   scriptId,
						Level:      script_entity.ScriptWatchLevel(level),
						Createtime: time.Now().Unix(),
					}); err != nil {
						logger.Ctx(ctx).Error("迁移关注失败", zap.Int64("script_id", scriptId),
							zap.Int64("user_id", uid), zap.Int("level", level), zap.Error(err))
					}
				}
				return nil
			}); err != nil {
				return err
			}

			return nil
		},
		Rollback: func(db *gorm.DB) error {
			return nil
		},
	}
}
