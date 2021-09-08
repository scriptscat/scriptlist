package repository

import (
	"fmt"
	"time"

	"github.com/scriptscat/scriptweb/internal/domain/script/entity"
	"github.com/scriptscat/scriptweb/internal/http/dto/request"
	"github.com/scriptscat/scriptweb/internal/pkg/cnt"
	"github.com/scriptscat/scriptweb/internal/pkg/db"
	"gorm.io/gorm"
)

type script struct {
	db *gorm.DB
}

func NewScript() Script {
	return &script{db: db.Db}
}

func NewTxScript(db *gorm.DB) Script {
	return &script{db: db}
}

func (s *script) Find(id int64) (*entity.Script, error) {
	ret := &entity.Script{}
	if err := s.db.Find(ret, "id=?", id).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *script) Save(script *entity.Script) error {
	return s.db.Save(script).Error
}

func (s *script) List(search *SearchList, page request.Pages) ([]*entity.Script, int64, error) {
	list := make([]*entity.Script, 0)
	scriptTbName := (&entity.Script{}).TableName()
	find := s.db.Model(&entity.Script{})
	if !search.Self {
		find.Where("public=?", entity.PUBLIC_SCRIPT)
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
			Order(tabname + ".score,createtime desc")
	case "updatetime":
		find = find.Order("updatetime desc,createtime desc")
	default:
		find = find.Order("createtime desc")
	}

	if search.Keyword != "" {
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
	if page == request.AllPage {
		if err := find.Select(scriptTbName + ".*").Scan(&list).Error; err != nil {
			return nil, 0, err
		}
	} else if err := find.Select(scriptTbName + ".*").Limit(page.Size()).Offset((page.Page() - 1) * page.Size()).Scan(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, num, nil
}
