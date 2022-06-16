package persistence

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/scriptscat/scriptlist/internal/interfaces/api/dto/request"
	"github.com/scriptscat/scriptlist/internal/pkg/cnt"
	"github.com/scriptscat/scriptlist/internal/pkg/errs"
	"github.com/scriptscat/scriptlist/internal/service/script/domain/entity"
	"github.com/scriptscat/scriptlist/internal/service/script/domain/repository"
	"github.com/scriptscat/scriptlist/pkg/gofound"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const SearchHotKeyword = "script:search:hot_keyword"

type script struct {
	statistics repository.Statistics
	db         *gorm.DB
	redis      *redis.Client
	found      *gofound.GOFound
}

func NewScript(db *gorm.DB, redis *redis.Client, found *gofound.GOFound) repository.Script {
	return &script{db: db, redis: redis, found: found,
		statistics: NewStatistics(db),
	}
}

func (s *script) Find(id int64) (*entity.Script, error) {
	ret := &entity.Script{}
	if err := s.db.First(ret, "id=?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.ErrScriptNotFound
		}
		return nil, err
	}
	return ret, nil
}

func (s *script) Save(script *entity.Script) error {
	save := false
	if script.ID > 0 && (script.Public == 0 || script.Unwell == 1 || script.Status != cnt.ACTIVE) {
		if err := s.deleteByGOFound(script.ID); err != nil {
			logrus.WithField("id", script.ID).WithError(err).Errorf("delete by gofound error")
		}
	} else {
		save = true
	}
	if err := s.db.Save(script).Error; err != nil {
		return err
	}
	if save {
		if err := s.PutGoFound(script); err != nil {
			logrus.WithField("id", script.ID).WithError(err).Errorf("put by gofound error")
		}
	}
	return nil
}

type goFoundModel struct {
	*entity.Script
	TodayDownload int64 `json:"today_download"`
	TotalDownload int64 `json:"total_download"`
	Score         int64 `json:"score"`
}

func (s *script) PutGoFound(script *entity.Script) error {
	total, today, err := s.statistics.FindByScriptId(script.ID)
	if err != nil {
		return err
	}
	if total == nil {
		total = &entity.ScriptStatistics{}
	}
	if today == nil {
		today = &entity.ScriptDateStatistics{}
	}
	return s.found.PutIndex("scripts", uint32(script.ID),
		script.Content+script.Name+script.Description,
		&goFoundModel{
			Script:        script,
			TotalDownload: total.Download,
			TodayDownload: today.Download,
			Score:         total.Score,
		})
}

func (s *script) deleteByGOFound(id int64) error {
	return s.found.RemoveIndex("scripts", uint32(id))
}

func (s *script) Search(search *repository.SearchList, page *request.Pages) ([]*entity.Script, int64, error) {
	exp := search.Sort
	if exp != "" {
		exp = "[document." + exp + "]"
	}
	search.Keyword = strings.TrimSpace(search.Keyword)
	ret := make([]*entity.Script, 0)
	info, err := gofound.QueryIndex[*entity.Script](s.found, "scripts", &gofound.QueryIndexRequest{
		Query:    search.Keyword,
		Page:     page.Page(),
		Limit:    page.Size(),
		Order:    gofound.ORDER_DESC,
		ScoreExp: exp,
	})
	if err != nil {
		return nil, 0, err
	}
	for _, v := range info.Documents {
		ret = append(ret, v.Document)
	}
	for _, word := range info.Words {
		s.redis.ZIncrBy(context.Background(), SearchHotKeyword, 1, word)
		//// 高亮
		//for i := range ret {
		//	ret[i].Name = strings.ReplaceAll(ret[i].Name, word,
		//		fmt.Sprintf("%s%s%s", "<span style='color:red'>", word, "</span>"))
		//	ret[i].Description = strings.ReplaceAll(ret[i].Description, word,
		//		fmt.Sprintf("%s%s%s", "<span style='color:red'>", word, "</span>"))
		//}
	}
	s.redis.ZIncrBy(context.Background(), SearchHotKeyword, 1, search.Keyword)
	return ret, int64(info.Total), nil
}

func (s *script) DropGoFound() error {
	return s.found.DropDatabase("scripts")
}

func (s *script) List(search *repository.SearchList, page *request.Pages) ([]*entity.Script, int64, error) {
	list := make([]*entity.Script, 0)
	scriptTbName := (&entity.Script{}).TableName()
	find := s.db.Model(&entity.Script{})
	if !search.Self {
		find.Where("public=? and unwell=?", entity.PUBLIC_SCRIPT, 2)
	}
	if len(search.Category) != 0 {
		tabname := (&entity.ScriptCategory{}).TableName()
		find = find.Joins("left join "+tabname+" on "+tabname+".script_id="+scriptTbName+".id").
			Where(tabname+".category_id in ?", search.Category)
	}
	if search.Domain != "" {
		tabname := (&entity.ScriptDomain{}).TableName()
		find = find.Joins("left join "+tabname+" on "+tabname+".script_id="+scriptTbName+".id").
			Where(tabname+".domain=?", search.Domain)
	}

	switch search.Sort {
	case "today_download":
		tabname := (&entity.ScriptDateStatistics{}).TableName()
		find = find.Joins(fmt.Sprintf("left join %s on %s.script_id=%s.id and %s.date=?", tabname, tabname, scriptTbName, tabname), time.Now().Format("2006-01-02")).
			Order(tabname + ".download desc,createtime desc")
	case "total_download":
		tabname := (&entity.ScriptStatistics{}).TableName()
		find = find.Joins(fmt.Sprintf("left join %s on %s.script_id=%s.id", tabname, tabname, scriptTbName)).
			Order(tabname + ".download desc,createtime desc")
	case "score":
		tabname := (&entity.ScriptStatistics{}).TableName()
		find = find.Joins(fmt.Sprintf("left join %s on %s.script_id=%s.id", tabname, tabname, scriptTbName)).
			Order(tabname + ".score desc,createtime desc")
	case "updatetime":
		find = find.Where("updatetime>0").Order("updatetime desc,createtime desc")
	default:
		find = find.Order("createtime desc")
	}

	if search.Keyword != "" {
		s.redis.ZIncrBy(context.Background(), SearchHotKeyword, 1, search.Keyword)
		find = find.Where("name like ? or description like ?", "%"+search.Keyword+"%", "%"+search.Keyword+"%")
	}
	if search.Status != cnt.UNKNOWN {
		find = find.Where("status=?", search.Status)
	}
	if search.Uid != 0 {
		find = find.Where("user_id=?", search.Uid)
	}
	var num int64
	if err := find.Count(&num).Error; err != nil {
		return nil, 0, err
	}
	find = find.Select(scriptTbName + ".*")
	if page != request.AllPage {
		find = find.Limit(page.Size()).Offset((page.Page() - 1) * page.Size())
	}
	if err := find.Scan(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, num, nil
}

func (s *script) FindSyncPrefix(uid int64, prefix string) ([]*entity.Script, error) {
	ret := make([]*entity.Script, 0)
	if err := s.db.Model(&entity.Script{}).Where("user_id=? and sync_url like ?", uid, prefix+"%").Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *script) FindSyncScript(page *request.Pages) ([]*entity.Script, error) {
	ret := make([]*entity.Script, 0)
	find := s.db.Model(&entity.Script{}).Where(
		"(sync_url!=null or sync_url!='') and status=? and archive=0", cnt.ACTIVE)
	if page != request.AllPage {
		find = find.Limit(page.Size()).Offset(page.Page() - 1*page.Size())
	}
	if err := find.Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *script) HotKeyword() ([]redis.Z, error) {
	return s.redis.ZRevRangeWithScores(context.Background(), SearchHotKeyword, 0, 10).Result()
}
