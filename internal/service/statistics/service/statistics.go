package service

import (
	"strconv"
	"time"

	"github.com/scriptscat/scriptlist/internal/interfaces/api/dto/respond"
	"github.com/scriptscat/scriptlist/internal/service/statistics/domain/repository"
	"github.com/scriptscat/scriptlist/internal/service/statistics/entity"
)

const (
	VIEW_STATISTICS     = "view"
	DOWNLOAD_STATISTICS = "download"
	UPDATE_STATISTICS   = "update"
)

type Statistics interface {
	Record(scriptId, scriptCodeId, user int64, ip, ua, statisticsToken string, download string) (bool, error)
	TodayDownload(scriptId int64) (int64, error)
	TotalDownload(scriptId int64) (int64, error)
	TodayUpdate(scriptId int64) (int64, error)
	TotalUpdate(scriptId int64) (int64, error)
	Deal() error
	DownloadUv(scriptId, days int64, date time.Time) (*respond.StatisticsChart, error)
	DownloadPv(scriptId, days int64, date time.Time) (*respond.StatisticsChart, error)
	DownloadUvNum(scriptId int64, days int, date time.Time) (int64, error)
	UpdateUv(scriptId, days int64, date time.Time) (*respond.StatisticsChart, error)
	UpdatePv(scriptId, days int64, date time.Time) (*respond.StatisticsChart, error)
	UpdateUvNum(scriptId int64, days int, date time.Time) (int64, error)
	RealtimeDownload(scriptId int64) (*respond.StatisticsChart, error)
	RealtimeUpdate(scriptId int64) (*respond.StatisticsChart, error)
	PageUvNum(scriptId int64, days int, member string, date time.Time) (int64, error)
	PagePvNum(scriptId int64, days int, date time.Time) (int64, error)
}

type statistics struct {
	repo repository.Statistics
}

func NewStatistics(repo repository.Statistics) Statistics {
	ret := &statistics{
		repo: repo,
	}
	return ret
}

func (s *statistics) Record(scriptId, scriptCodeId, user int64, ip, ua, statisticsToken string, download string) (bool, error) {
	item := &entity.StatisticsDownload{
		UserId:          user,
		Ip:              ip,
		ScriptId:        scriptId,
		ScriptCodeId:    scriptCodeId,
		Ua:              ua,
		StatisticsToken: statisticsToken,
		Createtime:      time.Now().Unix(),
	}
	if download == DOWNLOAD_STATISTICS {
		return s.repo.Download(item)
	} else if download == UPDATE_STATISTICS {
		return s.repo.CheckUpdate((*entity.StatisticsUpdate)(item))
	} else if download == VIEW_STATISTICS {
		return s.repo.PageView((*entity.StatisticsPageView)(item))
	}
	return true, nil
}

func (s *statistics) TodayDownload(scriptId int64) (int64, error) {
	return s.repo.DaysPvNum(scriptId, "download", 1, time.Now())
}

func (s *statistics) TotalDownload(scriptId int64) (int64, error) {
	return s.repo.TotalPv(scriptId, "download")
}

func (s *statistics) TodayUpdate(scriptId int64) (int64, error) {
	return s.repo.DaysPvNum(scriptId, "update", 1, time.Now())
}

func (s *statistics) TotalUpdate(scriptId int64) (int64, error) {
	return s.repo.TotalPv(scriptId, "update")
}

func (s *statistics) Deal() error {
	return s.repo.Deal()
}

func (s *statistics) DownloadUv(scriptId, days int64, date time.Time) (*respond.StatisticsChart, error) {
	return s.daysData(scriptId, days, date, "download", "uv")
}

func (s *statistics) DownloadPv(scriptId, days int64, date time.Time) (*respond.StatisticsChart, error) {
	return s.daysData(scriptId, days, date, "download", "pv")
}

func (s *statistics) DownloadUvNum(scriptId int64, days int, date time.Time) (int64, error) {
	return s.repo.DaysUvNum(scriptId, DOWNLOAD_STATISTICS, "uv", days, date)
}

func (s *statistics) UpdateUv(scriptId, days int64, date time.Time) (*respond.StatisticsChart, error) {
	return s.daysData(scriptId, days, date, "update", "uv")
}

func (s *statistics) UpdatePv(scriptId, days int64, date time.Time) (*respond.StatisticsChart, error) {
	return s.daysData(scriptId, days, date, "update", "pv")
}

func (s *statistics) UpdateUvNum(scriptId int64, days int, date time.Time) (int64, error) {
	return s.repo.DaysUvNum(scriptId, UPDATE_STATISTICS, "uv", days, date)
}

// PageUvNum 访客数和平台用户数
func (s *statistics) PageUvNum(scriptId int64, days int, member string, date time.Time) (int64, error) {
	return s.repo.DaysUvNum(scriptId, VIEW_STATISTICS, member, days, date)
}

// PagePvNum 浏览数
func (s *statistics) PagePvNum(scriptId int64, days int, date time.Time) (int64, error) {
	return s.repo.DaysPvNum(scriptId, VIEW_STATISTICS, days, date)
}

func (s *statistics) daysData(scriptId, days int64, date time.Time, op string, data string) (*respond.StatisticsChart, error) {
	t := date.Add(-time.Hour * 24 * time.Duration(days))
	var x []string
	var y []int64
	for i := int64(0); i < days; i++ {
		t = t.Add(time.Hour * 24)
		day := t.Format("2006/01/02")
		var num int64
		switch data {
		case "uv":
			num, _ = s.repo.DaysUvNum(scriptId, op, "uv", 1, t)
		case "pv":
			num, _ = s.repo.DaysPvNum(scriptId, op, 1, t)
		}
		x = append(x, day)
		y = append(y, num)
	}
	return &respond.StatisticsChart{
		X: x,
		Y: y,
	}, nil
}

func (s *statistics) RealtimeDownload(scriptId int64) (*respond.StatisticsChart, error) {
	return s.realtime(scriptId, "download")
}

func (s *statistics) RealtimeUpdate(scriptId int64) (*respond.StatisticsChart, error) {
	return s.realtime(scriptId, "update")
}

func (s *statistics) realtime(scriptId int64, op string) (*respond.StatisticsChart, error) {
	var nums []int64
	nums, _ = s.repo.Realtime(scriptId, op)
	l := len(nums)
	var x = make([]string, l)
	var y = make([]int64, l)
	for n, v := range nums {
		x[l-n-1] = strconv.Itoa(n+1) + "分钟前"
		y[l-n-1] = v
	}
	return &respond.StatisticsChart{
		X: x,
		Y: y,
	}, nil
}
