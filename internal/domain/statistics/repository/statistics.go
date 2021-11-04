package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/scriptscat/scriptweb/internal/domain/statistics/entity"
	"github.com/scriptscat/scriptweb/internal/pkg/db"
	"github.com/scriptscat/scriptweb/pkg/utils"
	"github.com/sirupsen/logrus"
)

type statistics struct {
}

// NewStatistics TODO: 遍历key很多,后续换专门的redis库存储
func NewStatistics() Statistics {
	ret := &statistics{}
	return ret
}

func (s *statistics) DealDay() error {
	if err := s.dealDay("download"); err != nil {
		logrus.Errorf("deal day download: %v", err)
		return err
	}
	if err := s.dealDay("update"); err != nil {
		logrus.Errorf("deal day update: %v", err)
		return err
	}
	return nil
}

func (s *statistics) dealDay(op string) error {
	var list []string
	var cursor uint64
	var err error
	for {
		fmt.Println("statistics:script:" + op + ":*:day:pf:*")
		list, cursor, err = db.Redis.ScanType(context.Background(), cursor, "statistics:script:"+op+":*", 100, "string").Result()
		if err != nil {
			return err
		}
		for _, v := range list {
			op := v[strings.LastIndex(v, ":")+1:]
			if op == "uv" || op == "member" {
				number, err := db.Redis.PFCount(context.Background(), v).Result()
				if err != nil {
					logrus.Infof("pfcount %s: %v", v, err)
					continue
				}
				date := v[strings.Index(v, ":day:pf:")+8:]
				date = date[:strings.Index(v, ":")]
				db.Redis.HIncrBy(context.Background(), v[strings.Index(v, ":day:pf:"):]+":day:"+op, date, number)
				db.Redis.Del(context.Background(), v)
				continue
			}
		}
		if cursor == 0 {
			return nil
		}
	}
}

func (s *statistics) DealRealtime() error {
	if err := s.dealRealtime("download"); err != nil {
		logrus.Errorf("deal realtime download: %v", err)
		return err
	}
	if err := s.dealRealtime("update"); err != nil {
		logrus.Errorf("deal realtime update: %v", err)
		return err
	}
	return nil
}

func (s *statistics) dealRealtime(op string) error {
	var list []string
	var cursor uint64
	var err error
	for {
		list, cursor, err = db.Redis.ScanType(context.Background(), cursor, "statistics:script:"+op+":*", 100, "hash").Result()
		if err != nil {
			return err
		}
		logrus.Infof("deal realtime %s: %d", op, len(list))
		max := time.Now().Add(-time.Hour).Unix() / 60
		for _, v := range list {
			if !strings.HasSuffix(v, "realtime") {
				continue
			}
			list, err := db.Redis.HGetAll(context.Background(), v).Result()
			if err != nil {
				logrus.Infof("hgetall %s: %v", v, err)
				continue
			}
			for k, num := range list {
				ik := utils.StringToInt64(k)
				if ik < max {
					db.Redis.HDel(context.Background(), v, k)
					t := time.Unix(ik*60, 0)
					id := v[strings.Index(v, "download:")+9:]
					id = id[:strings.Index(id, ":")]
					db.Redis.ZIncrBy(context.Background(), "statistics:script:"+op+":day:rank:"+t.Format("2006/01/02"), float64(utils.StringToInt64(num)), id)
				}
			}
		}
		if cursor == 0 {
			return nil
		}
	}
}

func (s *statistics) Download(entity *entity.StatisticsDownload) error {
	return s.save(entity, "download")
}

func (s *statistics) save(entity entity.Statistics, op string) error {
	key := "statistics:script:" + op + ":" + fmt.Sprintf("%d", entity.GetScriptId())
	date := time.Now().Format("2006/01/02")
	// 丢定时任务里合并
	db.Redis.PFAdd(context.Background(), key+fmt.Sprintf(":day:pf:%s:uv", date), entity.GetStatisticsToken())
	if entity.GetUserId() != 0 {
		db.Redis.PFAdd(context.Background(), key+fmt.Sprintf(":day:pf:%s:member", date), entity.GetUserId())
	}
	db.Redis.HIncrBy(context.Background(), key+":day:pv", date, 1)
	db.Redis.Incr(context.Background(), key+":total:pv")
	// 丢定时任务里清理
	db.Redis.HIncrBy(context.Background(), key+":realtime", strconv.FormatInt(time.Now().Unix()/60, 10), 1)
	return db.Db.Create(entity).Error
}

func (s *statistics) CheckUpdate(entity *entity.StatisticsUpdate) error {
	return s.save(entity, "update")
}

func (s *statistics) Query(scriptId, starttime, endtime, precision int64) ([]*Precision, error) {
	panic("implement me")
}

func (s *statistics) DayDownload(scriptId int64, day time.Time) (int64, error) {
	if day.Format("2006/01/02") == time.Now().Format("2006/01/02") {
		return db.Redis.PFCount(context.Background(), "statistics:script:download:"+fmt.Sprintf("%d", scriptId)+":day:uv").Result()
	}
	key := "statistics:script:download:" + fmt.Sprintf("%d", scriptId) + ":day:uv"
	return db.Redis.HGet(context.Background(), key, day.Format("2006/01/02")).Int64()
}

func (s *statistics) DayUpdate(scriptId int64, day time.Time) (int64, error) {
	if day.Format("2006/01/02") == time.Now().Format("2006/01/02") {
		return db.Redis.PFCount(context.Background(), "statistics:script:update:"+fmt.Sprintf("%d", scriptId)+":day:uv").Result()
	}
	key := "statistics:script:update:" + fmt.Sprintf("%d", scriptId) + ":day:uv"
	return db.Redis.HGet(context.Background(), key, day.Format("2006/01/02")).Int64()
}

func (s *statistics) TotalDownload(scriptId int64) (int64, error) {
	key := "statistics:script:download:" + fmt.Sprintf("%d", scriptId) + ":total:pv"
	return db.Redis.Get(context.Background(), key).Int64()
}

func (s *statistics) TotalUpdate(scriptId int64) (int64, error) {
	key := "statistics:script:update:" + fmt.Sprintf("%d", scriptId) + ":total:pv"
	return db.Redis.Get(context.Background(), key).Int64()
}
