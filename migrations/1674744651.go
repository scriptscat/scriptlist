package migrations

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/cago-frame/cago/database/redis"
	"github.com/cago-frame/cago/pkg/consts"
	"github.com/cago-frame/cago/pkg/logger"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/scriptlist/internal/model/entity/issue_entity"
	"github.com/scriptscat/scriptlist/internal/repository/issue_repo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func T1674744651() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1674744651",
		Migrate: func(db *gorm.DB) error {
			ctx := context.Background()
			// 迁移issue watch
			if err := db.AutoMigrate(&issue_entity.ScriptIssueWatch{}); err != nil {
				return err
			}
			if err := ScanKeys(ctx, "script:issue:watch:*", "hash", func(ctx context.Context, key string) error {
				logger.Ctx(ctx).Info("迁移反馈关注", zap.String("key", key))
				// 取出id
				issueId, err := strconv.ParseInt(key[strings.LastIndex(key, ":")+1:], 10, 64)
				if err != nil {
					return err
				}
				if issueId == 0 {
					return errors.New("id解析错误")
				}
				list, err := redis.Ctx(ctx).HGetAll(key).Result()
				if err != nil {
					return err
				}
				// v=1 关注 2=关注过,但是取消了
				for k, v := range list {
					uid, _ := strconv.ParseInt(k, 10, 64)
					status, _ := strconv.ParseInt(v, 10, 64)
					// 判断是否重复
					if ok, err := issue_repo.Watch().FindByUser(ctx, issueId, uid); err != nil {
						logger.Ctx(ctx).Error("迁移issue watch失败", zap.Int64("issue_id", issueId),
							zap.Int64("user_id", uid), zap.String("status", v), zap.Error(err))
						continue
					} else if ok != nil {
						continue
					}
					if err := issue_repo.Watch().Create(ctx, &issue_entity.ScriptIssueWatch{
						UserID:     uid,
						IssueID:    issueId,
						Status:     int32(status),
						Createtime: time.Now().Unix(),
					}); err != nil {
						logger.Ctx(ctx).Error("迁移issue watch失败", zap.Int64("issue_id", issueId),
							zap.Int64("user_id", uid), zap.String("status", v), zap.Error(err))
					}
				}
				return nil
			}); err != nil {
				return err
			}
			// 修改issue删除状态为2
			if err := db.Model(&issue_entity.ScriptIssue{}).Where("status = 0").
				Update("status", consts.DELETE).Error; err != nil {
				return err
			}
			// 修改comment删除状态为2
			if err := db.Model(&issue_entity.ScriptIssueComment{}).Where("status = 0").
				Update("status", consts.DELETE).Error; err != nil {
				return err
			}
			return nil
		},
		Rollback: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&issue_entity.ScriptIssueWatch{})
		},
	}
}
