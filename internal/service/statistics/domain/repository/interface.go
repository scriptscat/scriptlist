package repository

import (
	"time"

	"github.com/scriptscat/scriptlist/internal/service/statistics/entity"
)

type Precision struct {
	Time  int64
	Count int64
}

type Statistics interface {
	Download(entity *entity.StatisticsDownload) (bool, error)
	CheckUpdate(entity *entity.StatisticsUpdate) (bool, error)
	Deal() error
	RealtimeDownload(scriptId int64) ([]int64, error)
	RealtimeUpdate(scriptId int64) ([]int64, error)
	WeeklyUv(scriptId int64, t time.Time) (int64, error)
	WeeklyMember(scriptId int64, t time.Time) (int64, error)
	TotalPv(scriptId int64, op string) (int64, error)
	DayPv(scriptId int64, op string, day time.Time) (int64, error)
	DayUv(scriptId int64, op string, day time.Time) (int64, error)
	DayMember(scriptId int64, op string, day time.Time) (int64, error)
}