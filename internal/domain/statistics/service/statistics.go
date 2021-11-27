package service

import (
	"strconv"
	"time"

	"github.com/scriptscat/scriptlist/internal/domain/statistics/entity"
	"github.com/scriptscat/scriptlist/internal/domain/statistics/repository"
	"github.com/scriptscat/scriptlist/internal/http/dto/respond"
)

type Statistics interface {
	Record(scriptId, scriptCodeId, user int64, ip, ua, statisticsToken string, download bool) error
	TodayDownload(scriptId int64) (int64, error)
	TotalDownload(scriptId int64) (int64, error)
	TodayUpdate(scriptId int64) (int64, error)
	TotalUpdate(scriptId int64) (int64, error)
	Deal() error
	WeeklyUv(scriptId int64) (int64, error)
	WeeklyMember(scriptId int64) (int64, error)
	DownloadUv(scriptId, days int64, date time.Time) (*respond.StatisticsChart, error)
	DownloadPv(scriptId, days int64, date time.Time) (*respond.StatisticsChart, error)
	UpdateUv(scriptId, days int64, date time.Time) (*respond.StatisticsChart, error)
	UpdatePv(scriptId, days int64, date time.Time) (*respond.StatisticsChart, error)
	RealtimeDownload(scriptId int64) (*respond.StatisticsChart, error)
	RealtimeUpdate(scriptId int64) (*respond.StatisticsChart, error)
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

func (s *statistics) Record(scriptId, scriptCodeId, user int64, ip, ua, statisticsToken string, download bool) error {
	item := &entity.StatisticsDownload{
		UserId:          user,
		Ip:              ip,
		ScriptId:        scriptId,
		ScriptCodeId:    scriptCodeId,
		Ua:              ua,
		StatisticsToken: statisticsToken,
		Createtime:      time.Now().Unix(),
	}
	if download {
		return s.repo.Download(item)
	}
	return s.repo.CheckUpdate((*entity.StatisticsUpdate)(item))
}

func (s *statistics) TodayDownload(scriptId int64) (int64, error) {
	return s.repo.DayPv(scriptId, "download", time.Now())
}

func (s *statistics) TotalDownload(scriptId int64) (int64, error) {
	return s.repo.TotalPv(scriptId, "download")
}

func (s *statistics) TodayUpdate(scriptId int64) (int64, error) {
	return s.repo.DayPv(scriptId, "update", time.Now())
}

func (s *statistics) TotalUpdate(scriptId int64) (int64, error) {
	return s.repo.TotalPv(scriptId, "update")
}

func (s *statistics) Deal() error {
	return s.repo.Deal()
}

func (s *statistics) WeeklyUv(scriptId int64) (int64, error) {
	return s.repo.WeeklyUv(scriptId, time.Now())
}

func (s *statistics) WeeklyMember(scriptId int64) (int64, error) {
	return s.repo.WeeklyMember(scriptId, time.Now())
}

func (s *statistics) DownloadUv(scriptId, days int64, date time.Time) (*respond.StatisticsChart, error) {
	return s.daysData(scriptId, days, date, "download", "uv")
}

func (s *statistics) DownloadPv(scriptId, days int64, date time.Time) (*respond.StatisticsChart, error) {
	return s.daysData(scriptId, days, date, "download", "pv")
}

func (s *statistics) UpdateUv(scriptId, days int64, date time.Time) (*respond.StatisticsChart, error) {
	return s.daysData(scriptId, days, date, "update", "uv")
}

func (s *statistics) UpdatePv(scriptId, days int64, date time.Time) (*respond.StatisticsChart, error) {
	return s.daysData(scriptId, days, date, "update", "pv")
}

func (s *statistics) daysData(scriptId, days int64, date time.Time, op string, data string) (*respond.StatisticsChart, error) {
	t := date.Add(-time.Hour * 24 * time.Duration(days))
	var x, y []string
	for i := int64(0); i < days; i++ {
		t = t.Add(time.Hour * 24)
		day := t.Format("2006/01/02")
		var num int64
		switch data {
		case "uv":
			num, _ = s.repo.DayUv(scriptId, op, t)
		case "pv":
			num, _ = s.repo.DayPv(scriptId, op, t)
		case "member":
			num, _ = s.repo.DayMember(scriptId, op, t)
		}
		x = append(x, day)
		y = append(y, strconv.FormatInt(num, 10))
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
	switch op {
	case "download":
		nums, _ = s.repo.RealtimeDownload(scriptId)
	case "update":
		nums, _ = s.repo.RealtimeUpdate(scriptId)
	}
	l := len(nums)
	var x, y = make([]string, l), make([]string, l)
	for n, v := range nums {
		x[l-n-1] = strconv.Itoa(n+1) + "分钟前"
		y[l-n-1] = strconv.FormatInt(v, 10)
	}
	return &respond.StatisticsChart{
		X: x,
		Y: y,
	}, nil
}
