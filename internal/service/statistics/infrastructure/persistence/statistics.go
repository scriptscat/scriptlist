package persistence

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/scriptscat/scriptlist/internal/service/statistics/domain/repository"
	"github.com/scriptscat/scriptlist/internal/service/statistics/entity"
	"github.com/scriptscat/scriptlist/pkg/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type statistics struct {
	sync.Mutex
	db    *gorm.DB
	redis *redis.Client
}

// pf   statistics:script:@op:@id:day:uv:@date 30天过期
// hash statistics:script:@op:@id:day:uv @date @num 永不过期
// pf   statistics:script:weekly:@id:uv:@date 15天过期

// pf   statistics:script:@op:@id:day:member:@date 30天过期
// hash statistics:script:@op:@id:day:member @date @num 永不过期
// pf   statistics:script:weekly:@id:member:@date 15天过期

// hash statistics:script:@op:@id:day:pv @date @num 永不过期
// hash statistics:script:@op:@id:total:pv 永不过期

// hash statistics:script:@op:@id:realtime @time @num 永不过期,但定时删除项

// NewStatistics TODO: 遍历key很多,后续换专门的redis库存储
func NewStatistics(db *gorm.DB, redis *redis.Client) repository.Statistics {
	return &statistics{
		db:    db,
		redis: redis,
	}
}

func (s *statistics) Deal() error {
	s.Lock()
	defer s.Unlock()
	var list []string
	var cursor uint64
	var err error
	logrus.Infof("statistics deal start: %d", time.Now().Unix())
	defer logrus.Infof("statistics deal end: %d", time.Now().Unix())
	max := time.Now().Add(-time.Hour).Unix() / 60
	for {
		list, cursor, err = s.redis.ScanType(context.Background(), cursor, "statistics:script:*", 100, "string").Result()
		if err != nil {
			return err
		}
		for _, v := range list {
			split := strings.Split(v, ":")
			if len(split) != 7 {
				continue
			}
			prefix := strings.Join(split[:4], ":")
			if split[4] == "day" {
				date := split[6]
				today := time.Now().Format("2006/01/02")
				if s.redis.HExists(context.Background(), prefix+":day:"+split[5], date).Val() && date != today {
					continue
				}
				number, err := s.redis.PFCount(context.Background(), v).Result()
				if err != nil {
					logrus.Infof("pfcount %s: %v", v, err)
					continue
				}
				s.redis.HSet(context.Background(), prefix+":day:"+split[5], date, number)
				s.redis.Expire(context.Background(), v, time.Hour*24*30)
			} else if strings.HasSuffix(v, "realtime") {
				list, err := s.redis.HGetAll(context.Background(), v).Result()
				if err != nil {
					logrus.Infof("hgetall %s: %v", v, err)
					continue
				}
				for k, num := range list {
					ik := utils.StringToInt64(k)
					if ik < max {
						s.redis.HDel(context.Background(), v, k)
						t := time.Unix(ik*60, 0)
						id := v[strings.Index(v, "download:")+9:]
						id = id[:strings.Index(id, ":")]
						s.redis.ZIncrBy(context.Background(), prefix+":day:rank:"+t.Format("2006/01/02"), float64(utils.StringToInt64(num)), id)
					}
				}
			}
		}
		if cursor == 0 {
			return nil
		}
	}
}

func (s *statistics) Download(entity *entity.StatisticsDownload) (bool, error) {
	return s.save(entity, "download")
}

func (s *statistics) save(entity entity.Statistics, op string) (bool, error) {
	key := "statistics:script:" + op + ":" + fmt.Sprintf("%d", entity.GetScriptId())
	date := time.Now().Format("2006/01/02")
	// 丢定时任务里合并
	s.redis.PFAdd(context.Background(), key+fmt.Sprintf(":day:uv:%s", date), entity.GetStatisticsToken())
	if entity.GetUserId() != 0 {
		s.redis.PFAdd(context.Background(), key+fmt.Sprintf(":day:member:%s", date), entity.GetUserId())
	}
	s.redis.HIncrBy(context.Background(), key+":day:pv", date, 1)
	s.redis.Incr(context.Background(), key+":total:pv")
	// 丢定时任务里清理
	s.redis.HIncrBy(context.Background(), key+":realtime", strconv.FormatInt(time.Now().Unix()/60, 10), 1)
	if err := s.db.Create(entity).Error; err != nil {
		return false, err
	}
	// 判断ip是否操作过了
	return s.redis.SetNX(context.Background(), key+":ip:exist:day:"+date+":"+entity.GetIp(), "1", time.Hour*24).Result()
}

func (s *statistics) CheckUpdate(entity *entity.StatisticsUpdate) (bool, error) {
	return s.save(entity, "update")
}

func (s *statistics) WeeklyUv(scriptId int64, t time.Time) (int64, error) {
	return s.weekly(scriptId, "uv", t)
}

func (s *statistics) WeeklyMember(scriptId int64, t time.Time) (int64, error) {
	return s.weekly(scriptId, "member", t)
}

func (s *statistics) weekly(scriptId int64, op string, t time.Time) (int64, error) {
	weeklyKey := fmt.Sprintf("statistics:script:weekly:%d:%s:%s", scriptId, op, t.Format("2006/01/02"))
	if s.redis.Exists(context.Background(), weeklyKey).Val() != 1 {
		var weeklyDay []string
		t := t
		for i := 1; i <= 7; i++ {
			weeklyDay = append(weeklyDay,
				fmt.Sprintf("statistics:script:download:%d:day:%s:%s", scriptId, op, t.Add(-time.Hour*24*time.Duration(i)).Format("2006/01/02")),
				fmt.Sprintf("statistics:script:update:%d:day:%s:%s", scriptId, op, t.Add(-time.Hour*24*time.Duration(i)).Format("2006/01/02")),
			)
		}
		s.redis.PFMerge(context.Background(), weeklyKey, weeklyDay...)
	}
	ret, err := s.redis.PFCount(context.Background(), weeklyKey).Result()
	if err != nil {
		return 0, err
	}
	s.redis.Expire(context.Background(), weeklyKey, time.Hour*24*15)
	return ret, nil
}

func (s *statistics) RealtimeDownload(scriptId int64) ([]int64, error) {
	return s.realtime(scriptId, "download")
}

func (s *statistics) RealtimeUpdate(scriptId int64) ([]int64, error) {
	return s.realtime(scriptId, "update")
}

func (s *statistics) realtime(scriptId int64, op string) ([]int64, error) {
	var ret []int64
	t := time.Now().Unix() / 60
	for i := int64(0); i < 15; i++ {
		num, _ := s.redis.HGet(context.Background(), fmt.Sprintf("statistics:script:%s:%d:realtime", op, scriptId), strconv.FormatInt(t-i, 10)).Int64()
		ret = append(ret, num)
	}
	return ret, nil
}

func (s *statistics) TotalPv(scriptId int64, op string) (int64, error) {
	key := "statistics:script:" + op + ":" + fmt.Sprintf("%d", scriptId) + ":total:pv"
	return s.redis.Get(context.Background(), key).Int64()
}

func (s *statistics) DayPv(scriptId int64, op string, day time.Time) (int64, error) {
	key := "statistics:script:" + op + ":" + fmt.Sprintf("%d", scriptId) + ":day:pv"
	return s.redis.HGet(context.Background(), key, day.Format("2006/01/02")).Int64()
}

func (s *statistics) DayUv(scriptId int64, op string, day time.Time) (int64, error) {
	if day.Format("2006/01/02") == time.Now().Format("2006/01/02") {
		return s.redis.PFCount(context.Background(), "statistics:script:"+op+":"+fmt.Sprintf("%d", scriptId)+":day:uv").Result()
	}
	key := "statistics:script:" + op + ":" + fmt.Sprintf("%d", scriptId) + ":day:uv"
	return s.redis.HGet(context.Background(), key, day.Format("2006/01/02")).Int64()
}

func (s *statistics) DayMember(scriptId int64, op string, day time.Time) (int64, error) {
	if day.Format("2006/01/02") == time.Now().Format("2006/01/02") {
		return s.redis.PFCount(context.Background(), "statistics:script:"+op+":"+fmt.Sprintf("%d", scriptId)+":day:member").Result()
	}
	key := "statistics:script:" + op + ":" + fmt.Sprintf("%d", scriptId) + ":day:member"
	return s.redis.HGet(context.Background(), key, day.Format("2006/01/02")).Int64()
}
