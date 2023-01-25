package migrations

import (
	"context"

	"github.com/codfrm/cago/database/redis"
	"github.com/codfrm/cago/pkg/logger"
	"go.uber.org/zap"
)

func ScanKeys(ctx context.Context, match string, callback func(ctx context.Context, key string) error) error {
	var list []string
	var cursor uint64
	var err error
	for {
		list, cursor, err = redis.Ctx(ctx).ScanType(cursor, match, 100, "hash").Result()
		if err != nil {
			logger.Ctx(ctx).Error("扫描key错误", zap.String("match", match), zap.Uint64("cursor", cursor), zap.Error(err))
			return err
		}
		for _, v := range list {
			if err := callback(ctx, v); err != nil {
				logger.Ctx(ctx).Error("扫描key处理错误", zap.String("match", match), zap.Uint64("cursor", cursor),
					zap.String("key", v), zap.Error(err))
			}
		}
		if cursor == 0 {
			return nil
		}
	}
}
