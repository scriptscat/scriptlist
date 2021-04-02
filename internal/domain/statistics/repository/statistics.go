package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/scriptscat/scriptweb/internal/domain/statistics/entity"
	"github.com/scriptscat/scriptweb/internal/pkg/db"
)

// NOTE: 后续可换TSDB
type statistics struct {
}

func NewStatistics() Statistics {
	return &statistics{}
}

func (s *statistics) Download(entity *entity.StatisticsDownload) error {
	if err := s.cacheDayNum(entity.ScriptId, time.Now(), true); err != nil {
		return err
	}
	if err := s.cacheTotalNum(entity.ScriptId, true); err != nil {
		return err
	}
	if _, err := db.Redis.IncrBy(context.Background(), s.todayKey(entity.ScriptId, time.Now(), false), 1).Result(); err != nil {
		return err
	}
	if _, err := db.Redis.IncrBy(context.Background(), s.totalKey(entity.ScriptId, true), 1).Result(); err != nil {
		return err
	}
	return db.Db.Create(entity).Error
}

func (s *statistics) CheckUpdate(entity *entity.StatisticsUpdate) error {
	if err := s.cacheDayNum(entity.ScriptId, time.Now(), false); err != nil {
		return err
	}
	if err := s.cacheTotalNum(entity.ScriptId, false); err != nil {
		return err
	}
	if _, err := db.Redis.IncrBy(context.Background(), s.todayKey(entity.ScriptId, time.Now(), false), 1).Result(); err != nil {
		return err
	}
	if _, err := db.Redis.IncrBy(context.Background(), s.totalKey(entity.ScriptId, false), 1).Result(); err != nil {
		return err
	}
	return db.Db.Create(entity).Error
}

func (s *statistics) Query(scriptId, starttime, endtime, precision int64) ([]*Precision, error) {
	panic("implement me")
}

func (s *statistics) todayKey(scriptId int64, day time.Time, download bool) string {
	if download {
		return fmt.Sprintf("statis:%d:count:%s:download", scriptId, day.Format("2006:01:02"))
	}
	return fmt.Sprintf("statis:%d:count:%s:update", scriptId, day.Format("2006:01:02"))
}

func (s *statistics) totalKey(scriptId int64, download bool) string {
	if download {
		return fmt.Sprintf("statis:%d:count:download", scriptId)
	}
	return fmt.Sprintf("statis:%d:count:update", scriptId)
}

func (s *statistics) cacheDayNum(scriptId int64, day time.Time, download bool) error {
	todayKey := s.todayKey(scriptId, day, download)
	if ok, err := db.Redis.Exists(context.Background(), todayKey).Result(); err != nil {
		return err
	} else if ok == 0 {
		starttime := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location()).Unix()
		var num int64
		if download {
			num, err = s.DownloadCount(scriptId, starttime,
				starttime+86400)
		} else {
			num, err = s.UpdateCount(scriptId, starttime,
				starttime+86400)
		}
		if err != nil {
			return err
		}
		if err := db.Redis.Set(context.Background(), todayKey, num, time.Hour*24).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (s *statistics) cacheTotalNum(scriptId int64, download bool) error {
	totalKey := s.totalKey(scriptId, download)
	if ok, err := db.Redis.Exists(context.Background(), totalKey).Result(); err != nil {
		return err
	} else if ok == 0 {
		var num int64
		if download {
			num, err = s.DownloadCount(scriptId, 0, time.Now().Unix())
		} else {
			num, err = s.UpdateCount(scriptId, 0, time.Now().Unix())
		}
		if err != nil {
			return err
		}
		if err := db.Redis.Set(context.Background(), totalKey, num, 0).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (s *statistics) DayDownload(scriptId int64, day time.Time) (int64, error) {
	if err := s.cacheDayNum(scriptId, day, true); err != nil {
		return 0, err
	}
	todayKey := s.todayKey(scriptId, day, true)
	return db.Redis.Get(context.Background(), todayKey).Int64()
}

func (s *statistics) DayUpdate(scriptId int64, day time.Time) (int64, error) {
	if err := s.cacheDayNum(scriptId, day, false); err != nil {
		return 0, err
	}
	todayKey := s.todayKey(scriptId, day, false)
	return db.Redis.Get(context.Background(), todayKey).Int64()
}

func (s *statistics) TotalDownload(scriptId int64) (int64, error) {
	if err := s.cacheTotalNum(scriptId, true); err != nil {
		return 0, err
	}
	totalKey := s.totalKey(scriptId, true)
	return db.Redis.Get(context.Background(), totalKey).Int64()
}

func (s *statistics) TotalUpdate(scriptId int64) (int64, error) {
	if err := s.cacheTotalNum(scriptId, false); err != nil {
		return 0, err
	}
	totalKey := s.totalKey(scriptId, false)
	return db.Redis.Get(context.Background(), totalKey).Int64()
}

func (s *statistics) DownloadCount(scriptId, starttime, endtime int64) (int64, error) {
	var cnt int64
	if err := db.Db.Model(&entity.StatisticsDownload{}).Where("script_id=?", scriptId).Where("createtime>? and createtime<?", starttime, endtime).Count(&cnt).Error; err != nil {
		return 0, err
	}
	return cnt, nil
}

func (s *statistics) UpdateCount(scriptId, starttime, endtime int64) (int64, error) {
	var cnt int64
	if err := db.Db.Model(&entity.StatisticsUpdate{}).Where("script_id=?", scriptId).Where("createtime>? and createtime<?", starttime, endtime).Count(&cnt).Error; err != nil {
		return 0, err
	}
	return cnt, nil
}
