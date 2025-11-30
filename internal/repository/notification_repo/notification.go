package notification_repo

import (
	"context"

	"github.com/cago-frame/cago/database/db"
	"github.com/cago-frame/cago/pkg/consts"
	api "github.com/scriptscat/scriptlist/internal/api/notification"
	"github.com/scriptscat/scriptlist/internal/model/entity/notification_entity"
)

//go:generate mockgen -source=./notification.go -destination=./mock/notification.go
type NotificationRepo interface {
	Find(ctx context.Context, id int64) (*notification_entity.Notification, error)
	FindByUserID(ctx context.Context, userID int64, id int64) (*notification_entity.Notification, error)
	FindPage(ctx context.Context, userID int64, req *api.ListRequest) ([]*notification_entity.Notification, int64, error)
	Create(ctx context.Context, notification *notification_entity.Notification) error
	BatchCreate(ctx context.Context, notifications []*notification_entity.Notification) error
	Update(ctx context.Context, notification *notification_entity.Notification) error
	Delete(ctx context.Context, userID int64, id int64) error
	BatchDelete(ctx context.Context, userID int64, ids []int64) (int64, error)
	MarkRead(ctx context.Context, userID int64, id int64, readTime int64) error
	BatchMarkRead(ctx context.Context, userID int64, ids []int64, readTime int64) (int64, error)
	MarkAllRead(ctx context.Context, userID int64, notifyType int32, readTime int64) (int64, error)
	CountUnread(ctx context.Context, userID int64) (int64, error)
	CountUnreadByType(ctx context.Context, userID int64) (map[int32]int64, error)
	ClearRead(ctx context.Context, userID int64) (int64, error)
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

func (r *notificationRepo) Find(ctx context.Context, id int64) (*notification_entity.Notification, error) {
	ret := &notification_entity.Notification{}
	if err := db.Ctx(ctx).Where("id=? and status=?", id, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (r *notificationRepo) FindByUserID(ctx context.Context, userID int64, id int64) (*notification_entity.Notification, error) {
	ret := &notification_entity.Notification{}
	if err := db.Ctx(ctx).Where("id=? and user_id=? and status=?", id, userID, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (r *notificationRepo) FindPage(ctx context.Context, userID int64, req *api.ListRequest) ([]*notification_entity.Notification, int64, error) {
	var list []*notification_entity.Notification
	var count int64
	query := db.Ctx(ctx).Model(&notification_entity.Notification{}).
		Where("user_id=? and status=?", userID, consts.ACTIVE)

	// 按已读状态筛选
	if req.ReadStatus != 0 {
		query = query.Where("read_status=?", req.ReadStatus)
	}

	// 按通知类型筛选
	if req.Type != 0 {
		query = query.Where("type=?", req.Type)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("createtime desc").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, count, nil
}

func (r *notificationRepo) Create(ctx context.Context, notification *notification_entity.Notification) error {
	return db.Ctx(ctx).Create(notification).Error
}

func (r *notificationRepo) BatchCreate(ctx context.Context, notifications []*notification_entity.Notification) error {
	if len(notifications) == 0 {
		return nil
	}
	return db.Ctx(ctx).Create(&notifications).Error
}

func (r *notificationRepo) Update(ctx context.Context, notification *notification_entity.Notification) error {
	return db.Ctx(ctx).Model(&notification).
		Select("*").
		Updates(notification).Error
}

func (r *notificationRepo) Delete(ctx context.Context, userID int64, id int64) error {
	return db.Ctx(ctx).Model(&notification_entity.Notification{}).
		Where("id=? and user_id=?", id, userID).
		Update("status", consts.DELETE).Error
}

func (r *notificationRepo) BatchDelete(ctx context.Context, userID int64, ids []int64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	result := db.Ctx(ctx).Model(&notification_entity.Notification{}).
		Where("user_id=? and id in (?)", userID, ids).
		Update("status", consts.DELETE)
	return result.RowsAffected, result.Error
}

func (r *notificationRepo) MarkRead(ctx context.Context, userID int64, id int64, readTime int64) error {
	return db.Ctx(ctx).Model(&notification_entity.Notification{}).
		Where("id=? and user_id=? and read_status=?", id, userID, notification_entity.StatusUnread).
		Updates(map[string]interface{}{
			"read_status": notification_entity.StatusRead,
			"read_time":   readTime,
		}).Error
}

func (r *notificationRepo) BatchMarkRead(ctx context.Context, userID int64, ids []int64, readTime int64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	result := db.Ctx(ctx).Model(&notification_entity.Notification{}).
		Where("user_id=? and id in (?) and read_status=?", userID, ids, notification_entity.StatusUnread).
		Updates(map[string]interface{}{
			"read_status": notification_entity.StatusRead,
			"read_time":   readTime,
		})
	return result.RowsAffected, result.Error
}

func (r *notificationRepo) MarkAllRead(ctx context.Context, userID int64, notifyType int32, readTime int64) (int64, error) {
	query := db.Ctx(ctx).Model(&notification_entity.Notification{}).
		Where("user_id=? and read_status=? and status=?", userID, notification_entity.StatusUnread, consts.ACTIVE)

	if notifyType != 0 {
		query = query.Where("type=?", notifyType)
	}

	result := query.Updates(map[string]interface{}{
		"read_status": notification_entity.StatusRead,
		"read_time":   readTime,
	})
	return result.RowsAffected, result.Error
}

func (r *notificationRepo) CountUnread(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := db.Ctx(ctx).Model(&notification_entity.Notification{}).
		Where("user_id=? and read_status=? and status=?", userID, notification_entity.StatusUnread, consts.ACTIVE).
		Count(&count).Error
	return count, err
}

func (r *notificationRepo) CountUnreadByType(ctx context.Context, userID int64) (map[int32]int64, error) {
	type Result struct {
		Type  int32 `json:"type"`
		Count int64 `json:"count"`
	}
	var results []Result
	err := db.Ctx(ctx).Model(&notification_entity.Notification{}).
		Select("type, count(*) as count").
		Where("user_id=? and read_status=? and status=?", userID, notification_entity.StatusUnread, consts.ACTIVE).
		Group("type").
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	countMap := make(map[int32]int64)
	for _, r := range results {
		countMap[r.Type] = r.Count
	}
	return countMap, nil
}

func (r *notificationRepo) ClearRead(ctx context.Context, userID int64) (int64, error) {
	result := db.Ctx(ctx).Model(&notification_entity.Notification{}).
		Where("user_id=? and read_status=? and status=?", userID, notification_entity.StatusRead, consts.ACTIVE).
		Update("status", consts.DELETE)
	return result.RowsAffected, result.Error
}
