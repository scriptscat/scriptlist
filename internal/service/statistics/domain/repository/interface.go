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
	PageView(entity *entity.StatisticsPageView) (bool, error)
	Deal() error
	Realtime(scriptId int64, op string) ([]int64, error)
	DaysUvNum(scriptId int64, op, member string, days int, t time.Time) (int64, error)
	DaysPvNum(scriptId int64, op string, days int, t time.Time) (int64, error)
	TotalPv(scriptId int64, op string) (int64, error)
}
