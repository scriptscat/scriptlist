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
	"github.com/scriptscat/scriptlist/internal/service/statistics/service"
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
				// 实时数据清理
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

func (s *statistics) CheckUpdate(entity *entity.StatisticsUpdate) (bool, error) {
	return s.save(entity, "update")
}

func (s *statistics) PageView(entity *entity.StatisticsPageView) (bool, error) {
	return s.save(entity, "view")
}

func (s *statistics) save(entity entity.Statistics, op string) (bool, error) {
	key := "statistics:script:" + op + ":" + fmt.Sprintf("%d", entity.GetScriptId())
	date := time.Now().Format("2006/01/02")
	// 丢定时任务里合并
	// 储存统计token计算uv
	s.redis.PFAdd(context.Background(), key+fmt.Sprintf(":day:uv:%s", date), entity.GetStatisticsToken())
	if entity.GetUserId() != 0 {
		// 储存用户id计算member
		s.redis.PFAdd(context.Background(), key+fmt.Sprintf(":day:member:%s", date), entity.GetUserId())
	}
	// 日pv
	s.redis.HIncrBy(context.Background(), key+":day:pv", date, 1)
	// 总pv
	s.redis.Incr(context.Background(), key+":total:pv")
	// 丢定时任务里清理
	s.redis.HIncrBy(context.Background(), key+":realtime", strconv.FormatInt(time.Now().Unix()/60, 10), 1)
	if op != service.VIEW_STATISTICS {
		if err := s.db.Create(entity).Error; err != nil {
			return false, err
		}
	}
	// 判断ip是否操作过了
	return s.redis.SetNX(context.Background(), key+":ip:exist:day:"+date+":"+entity.GetIp(), "1", time.Hour*24).Result()
}

// DaysUvNum 获取某一段时间某一个操作的uv或者member数量
// op: download, update, view
// member: member, uv
func (s *statistics) DaysUvNum(scriptId int64, op, member string, days int, t time.Time) (int64, error) {
	if days == 1 {
		ret, err := s.redis.PFCount(context.Background(),
			fmt.Sprintf("statistics:script:%s:%d:day:%s:%s", op, scriptId, member, t.Format("2006/01/02")),
		).Result()
		if err != nil {
			return 0, err
		}
		return ret, nil
	}
	key := fmt.Sprintf("statistics:script:cache:uv:%d:%s:%s", scriptId, op, t.Format("2006/01/02"))
	if s.redis.Exists(context.Background(), key).Val() != 1 {
		var dayKey []string
		t := t
		for i := 1; i <= days; i++ {
			dayKey = append(dayKey,
				fmt.Sprintf("statistics:script:%s:%d:day:%s:%s", op, scriptId, member, t.Add(-time.Hour*24*time.Duration(i)).Format("2006/01/02")),
			)
		}
		s.redis.PFMerge(context.Background(), key, dayKey...)
	}
	ret, err := s.redis.PFCount(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}
	s.redis.Expire(context.Background(), key, time.Hour*24*15)
	return ret, nil
}

// DaysPvNum 获取某一段时间某一个操作的pv数量
// op: download, update, view
func (s *statistics) DaysPvNum(scriptId int64, op string, days int, t time.Time) (int64, error) {
	var num int64
	key := fmt.Sprintf("statistics:script:%s:%d:day:pv", op, scriptId)
	for i := 0; i < days; i++ {
		val, err := s.redis.HGet(context.Background(), key, t.Add(-time.Hour*24*time.Duration(i)).Format("2006/01/02")).Int64()
		if err != nil {
			continue
		}
		num += val
	}
	return num, nil
}

// Realtime 获取某操作的实时记录
func (s *statistics) Realtime(scriptId int64, op string) ([]int64, error) {
	var ret []int64
	t := time.Now().Unix() / 60
	for i := int64(0); i < 15; i++ {
		num, _ := s.redis.HGet(context.Background(), fmt.Sprintf("statistics:script:%s:%d:realtime", op, scriptId), strconv.FormatInt(t-i, 10)).Int64()
		ret = append(ret, num)
	}
	return ret, nil
}

// TotalPv 获取某操作的总pv数量
func (s *statistics) TotalPv(scriptId int64, op string) (int64, error) {
	key := "statistics:script:" + op + ":" + fmt.Sprintf("%d", scriptId) + ":total:pv"
	return s.redis.Get(context.Background(), key).Int64()
}
