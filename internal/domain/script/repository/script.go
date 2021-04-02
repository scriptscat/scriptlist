package repository

import (
	"github.com/scriptscat/scriptweb/internal/domain/script/entity"
	"github.com/scriptscat/scriptweb/internal/interfaces/dto/request"
	"github.com/scriptscat/scriptweb/internal/pkg/cnt"
	"github.com/scriptscat/scriptweb/internal/pkg/db"
)

type script struct {
}

func NewScript() Script {
	return &script{}
}

func (s *script) Find(id int64) (*entity.Script, error) {
	ret := &entity.Script{}
	if err := db.Db.Find(ret, "id=?", id).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *script) Save(script *entity.Script) error {
	return db.Db.Save(script).Error
}

func (s *script) List(search *SearchList, page request.Pages) ([]*entity.Script, int64, error) {
	list := make([]*entity.Script, 0)
	find := db.Db.Model(&entity.Script{}).Order("createtime desc")
	if search.Category != 0 {
		tabname := (&entity.ScriptCategory{}).TableName()
		find = find.Joins("left join "+tabname+" on "+tabname+".script_id="+(&entity.Script{}).TableName()+".id").
			Where(tabname+".category_id=?", search.Category)
	}
	if search.Keyword != "" {
		find = find.Where("name like ? or description like ?", "%"+search.Keyword+"%", "%"+search.Keyword+"%")
	}
	if search.Status == cnt.UNKNOWN {
		find = find.Where("status=?", search.Status)
	}
	if search.Uid != 0 {
		find = find.Where("user_id=?", search.Uid)
	}
	var num int64
	if err := find.Count(&num).Error; err != nil {
		return nil, 0, err
	}
	if err := find.Limit(page.Size()).Offset((page.Page() - 1) * page.Size()).Scan(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, num, nil
}
