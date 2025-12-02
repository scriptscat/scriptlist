package sender

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cago-frame/cago/pkg/consts"
	"github.com/cago-frame/cago/pkg/logger"
	"github.com/scriptscat/scriptlist/internal/model/entity/notification_entity"
	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
	"github.com/scriptscat/scriptlist/internal/repository/notification_repo"
	"go.uber.org/zap"
)

type app struct {
}

func NewApp() Sender {
	return &app{}
}

func (s *app) Send(ctx context.Context, user *user_entity.User, content string, options *SendOptions) error {
	now := time.Now().Unix()
	paramsJson, err := json.Marshal(options.Params)
	if err != nil {
		logger.Ctx(ctx).Error("序列化通知参数失败", zap.Error(err))
		return err
	}

	m := &notification_entity.Notification{
		UserID:     user.ID,
		FromUserID: 0,
		Type:       options.Type,
		Title:      options.Title,
		Content:    content,
		ReadStatus: 0,
		ReadTime:   0,
		Link:       options.Link,
		Params:     string(paramsJson),
		Status:     consts.ACTIVE,
		Createtime: now,
		Updatetime: now,
	}

	if options.From != nil {
		m.FromUserID = user.ID
	}

	if err := notification_repo.Notification().Create(ctx, m); err != nil {
		logger.Ctx(ctx).Error("创建应用内通知失败", zap.Error(err))
		return err
	}

	return nil
}
