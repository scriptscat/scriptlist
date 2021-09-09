package service

import (
	"strconv"
	"time"

	"github.com/scriptscat/scriptweb/internal/domain/safe/repository"
	"github.com/scriptscat/scriptweb/internal/pkg/errs"
)

type Rate interface {
	Rate(userinfo *repository.RateUserInfo, rule *repository.RateRule, f func() error) error
}

type rate struct {
	repo repository.Rate
}

func NewRate(repo repository.Rate) Rate {
	return &rate{repo: repo}
}

func (r *rate) Rate(userinfo *repository.RateUserInfo, rule *repository.RateRule, f func() error) error {
	t, err := r.repo.GetLastOpTime(strconv.FormatInt(userinfo.Uid, 10), rule.Name)
	if err != nil {
		return err
	}
	if t > time.Now().Unix()-rule.Interval {
		return errs.NewOperationTimeToShort(rule)
	}
	c, err := r.repo.GetDayOpCnt(strconv.FormatInt(userinfo.Uid, 10), rule.Name)
	if err != nil {
		return err
	}
	if rule.DayMax > 0 && c > rule.DayMax {
		return errs.NewOperationMax(rule)
	}
	if err := r.repo.SetLastOpTime(strconv.FormatInt(userinfo.Uid, 10), rule.Name, time.Now().Unix()); err != nil {
		return err
	}
	if err := f(); err != nil {
		_ = r.repo.SetLastOpTime(strconv.FormatInt(userinfo.Uid, 10), rule.Name, 0)
		return err
	}
	return r.repo.SetLastOpTime(strconv.FormatInt(userinfo.Uid, 10), rule.Name, time.Now().Unix())
}
