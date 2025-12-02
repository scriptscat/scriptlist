package notification_repo

import (
	"context"
	"time"

	"github.com/cago-frame/cago/database/db"
	"github.com/cago-frame/cago/pkg/consts"
	api "github.com/scriptscat/scriptlist/internal/api/notification"
	"github.com/scriptscat/scriptlist/internal/model/entity/notification_entity"
)

//go:generate mockgen -source=./notification.go -destination=./mock/notification.go
type NotificationRepo interface {
	Find(ctx context.Context, userId, id int64) (*notification_entity.Notification, error)
	FindPage(ctx context.Context, userId int64, req *api.ListRequest) ([]*notification_entity.Notification, int64, error)
	Create(ctx context.Context, notification *notification_entity.Notification) error
	Update(ctx context.Context, notification *notification_entity.Notification) error

	CountUnread(ctx context.Context, userId int64, templateType int32) (int64, error)
	BatchMarkRead(ctx context.Context, userId int64, ids []int64) error
}

var defaultNotification NotificationRepo

func Notification() NotificationRepo {
	return defaultNotification
}

func RegisterNotification(n NotificationRepo) {
	defaultNotification = n
}

type notificationRepo struct {
}

func NewNotificationRepo() NotificationRepo {
	return &notificationRepo{}
}

func (r *notificationRepo) Find(ctx context.Context, userId, id int64) (*notification_entity.Notification, error) {
	ret := &notification_entity.Notification{}
	if err := db.Ctx(ctx).Where("id=? and user_id=? and status=?", id, userId, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (r *notificationRepo) Create(ctx context.Context, notification *notification_entity.Notification) error {
	return db.Ctx(ctx).Create(notification).Error
}

func (r *notificationRepo) Update(ctx context.Context, notification *notification_entity.Notification) error {
	return db.Ctx(ctx).Save(notification).Error
}

func (r *notificationRepo) FindPage(ctx context.Context, userId int64, req *api.ListRequest) ([]*notification_entity.Notification, int64, error) {
	var list []*notification_entity.Notification
	var count int64
	find := db.Ctx(ctx).Model(&notification_entity.Notification{}).
		Where("user_id=? and status=?", userId, consts.ACTIVE)

	// 根据已读状态筛选
	if req.ReadStatus != 0 {
		find = find.Where("read_status=?", req.ReadStatus)
	}

	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := find.Order("createtime desc").Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (r *notificationRepo) CountUnread(ctx context.Context, userId int64, templateType int32) (int64, error) {
	var count int64
	query := db.Ctx(ctx).Model(&notification_entity.Notification{}).
		Where("user_id=? and read_status=? and status=?", userId, notification_entity.StatusUnread, consts.ACTIVE)
	if templateType != 0 {
		query = query.Where("type=?", templateType)
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *notificationRepo) BatchMarkRead(ctx context.Context, userId int64, ids []int64) error {
	db := db.Ctx(ctx).Model(&notification_entity.Notification{}).
		Where("user_id=? and read_status=? and status=?", userId, notification_entity.StatusUnread, consts.ACTIVE)
	if len(ids) > 0 {
		db = db.Where("id IN (?)", ids)
	}
	return db.Updates(map[string]interface{}{
		"read_status": notification_entity.StatusRead,
		"read_time":   time.Now().Unix(),
	}).Error
}
