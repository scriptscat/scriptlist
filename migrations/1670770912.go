package migrations

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/codfrm/cago/database/redis"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"github.com/scriptscat/scriptlist/internal/repository/user_repo"
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
			if err := ScanKeys(ctx, "script:watch:*", "hash", func(ctx context.Context, key string) error {
				logger.Ctx(ctx).Info("迁移关注", zap.String("key", key))
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
					// 判断是否重复
					if ok, err := script_repo.ScriptWatch().FindByUser(ctx, scriptId, uid); err != nil {
						logger.Ctx(ctx).Error("迁移关注失败",
							zap.Int64("level", int64(level)),
							zap.Int64("script_id", scriptId), zap.Int64("user_id", uid), zap.Error(err))
						continue
					} else if ok != nil {
						// 存在
						continue
					}
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
			// 激活状态修改为2
			if err := db.Model(&script_entity.Script{}).Where("archive = 0 or archive is null").
				Update("archive", script_entity.IsActive).Error; err != nil {
				return err
			}
			// 删除状态修改为2
			if err := db.Model(&script_entity.Script{}).Where("status = 0").
				Update("status", consts.DELETE).Error; err != nil {
				return err
			}
			// 迁移webhook
			if err := db.AutoMigrate(&user_entity.UserConfig{}); err != nil {
				return err
			}
			if err := ScanKeys(ctx, "user:token:user:*", "string", func(ctx context.Context, key string) error {
				logger.Ctx(ctx).Info("迁移webhook", zap.String("key", key))
				suid := key[strings.LastIndex(key, ":")+1:]
				token, err := redis.Ctx(ctx).Get(key).Result()
				if err != nil {
					return err
				}
				uid, err := strconv.ParseInt(suid, 10, 64)
				if err != nil {
					return err
				}
				cfg, err := user_repo.UserConfig().FindByUserID(ctx, uid)
				if err != nil {
					return err
				}
				if cfg == nil {
					cfg = &user_entity.UserConfig{
						Uid:        uid,
						Token:      token,
						Notify:     &user_entity.Notify{},
						Createtime: time.Now().Unix(),
					}
					return user_repo.UserConfig().Create(ctx, cfg)
				} else {
					//cfg.Token = token
					//cfg.Updatetime = time.Now().Unix()
					//return user_repo.UserConfig().Update(ctx, cfg)
					logger.Ctx(ctx).Info("用户已存在token", zap.Int64("uid", uid), zap.String("old_token", token), zap.String("token", cfg.Token))
					return nil
				}
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
