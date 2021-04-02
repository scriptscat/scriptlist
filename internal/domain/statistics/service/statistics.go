package service

import (
	"time"

	"github.com/scriptscat/scriptweb/internal/domain/statistics/entity"
	"github.com/scriptscat/scriptweb/internal/domain/statistics/repository"
)

type Statistics interface {
	Record(scriptId, scriptCodeId, user int64, ip, ua string, download bool) error
	TodayDownload(scriptId int64) (int64, error)
	TotalDownload(scriptId int64) (int64, error)
	TodayUpdate(scriptId int64) (int64, error)
	TotalUpdate(scriptId int64) (int64, error)
}

type statistics struct {
	repo repository.Statistics
}

func NewStatistics(repo repository.Statistics) Statistics {
	return &statistics{
		repo: repo,
	}
}

func (s *statistics) Record(scriptId, scriptCodeId, user int64, ip, ua string, download bool) error {
	item := &entity.StatisticsDownload{
		UserId:       user,
		Ip:           ip,
		ScriptId:     scriptId,
		ScriptCodeId: scriptCodeId,
		Ua:           ua,
		Createtime:   time.Now().Unix(),
	}
	if download {
		return s.repo.Download(item)
	}
	return s.repo.CheckUpdate((*entity.StatisticsUpdate)(item))
}

func (s *statistics) TodayDownload(scriptId int64) (int64, error) {
	return s.repo.DayDownload(scriptId, time.Now())
}

func (s *statistics) TotalDownload(scriptId int64) (int64, error) {
	return s.repo.TotalDownload(scriptId)
}

func (s *statistics) TodayUpdate(scriptId int64) (int64, error) {
	return s.repo.DayUpdate(scriptId, time.Now())
}

func (s *statistics) TotalUpdate(scriptId int64) (int64, error) {
	return s.repo.TotalUpdate(scriptId)
}
