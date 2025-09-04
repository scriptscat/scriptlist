package feedback_repo

import (
	"context"

	"github.com/cago-frame/cago/database/db"
	"github.com/cago-frame/cago/pkg/consts"
	"github.com/cago-frame/cago/pkg/utils/httputils"
	"github.com/scriptscat/scriptlist/internal/model/entity/feedback_entity"
)

type FeedbackRepo interface {
	Find(ctx context.Context, id int64) (*feedback_entity.Feedback, error)
	FindPage(ctx context.Context, page httputils.PageRequest) ([]*feedback_entity.Feedback, int64, error)
	Create(ctx context.Context, feedback *feedback_entity.Feedback) error
	Update(ctx context.Context, feedback *feedback_entity.Feedback) error
	Delete(ctx context.Context, id int64) error
}

var defaultFeedback FeedbackRepo

func Feedback() FeedbackRepo {
	return defaultFeedback
}

func RegisterFeedback(i FeedbackRepo) {
	defaultFeedback = i
}

type feedbackRepo struct {
}

func NewFeedback() FeedbackRepo {
	return &feedbackRepo{}
}

func (u *feedbackRepo) Find(ctx context.Context, id int64) (*feedback_entity.Feedback, error) {
	ret := &feedback_entity.Feedback{}
	if err := db.Ctx(ctx).Where("id=? and status=?", id, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *feedbackRepo) Create(ctx context.Context, feedback *feedback_entity.Feedback) error {
	return db.Ctx(ctx).Create(feedback).Error
}

func (u *feedbackRepo) Update(ctx context.Context, feedback *feedback_entity.Feedback) error {
	return db.Ctx(ctx).Updates(feedback).Error
}

func (u *feedbackRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Model(&feedback_entity.Feedback{}).Where("id=?", id).Update("status", consts.DELETE).Error
}

func (u *feedbackRepo) FindPage(ctx context.Context, page httputils.PageRequest) ([]*feedback_entity.Feedback, int64, error) {
	var list []*feedback_entity.Feedback
	var count int64
	find := db.Ctx(ctx).Model(&feedback_entity.Feedback{}).Where("status=?", consts.ACTIVE)
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := find.Order("createtime desc").Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}
