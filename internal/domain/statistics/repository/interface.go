package repository

import (
	"time"

	"github.com/scriptscat/scriptweb/internal/domain/statistics/entity"
)

type Precision struct {
	Time  int64
	Count int64
}

type Statistics interface {
	Download(entity *entity.StatisticsDownload) error
	CheckUpdate(entity *entity.StatisticsUpdate) error
	Query(scriptId, starttime, endtime, precision int64) ([]*Precision, error)
	DayDownload(scriptId int64, day time.Time) (int64, error)
	TotalDownload(scriptId int64) (int64, error)
	DayUpdate(scriptId int64, day time.Time) (int64, error)
	TotalUpdate(scriptId int64) (int64, error)
	DealDay() error
	DealRealtime() error
}
