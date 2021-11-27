package repository

import (
	"github.com/scriptscat/scriptlist/internal/domain/issue/entity"
	"github.com/scriptscat/scriptlist/internal/http/dto/request"
	"github.com/scriptscat/scriptlist/internal/pkg/db"
	"gorm.io/gorm"
)

type issue struct {
	db *gorm.DB
}

func NewIssue() Issue {
	return &issue{db: db.Db}
}

func (i *issue) List(scriptId int64, keyword string, labels []string, status int, page request.Pages) ([]*entity.ScriptIssue, error) {
	list := make([]*entity.ScriptIssue, 0)
	find := i.db.Model(&entity.ScriptIssue{}).Where("script_id=?", scriptId).Order("createtime desc")
	if keyword != "" {
		find = find.Where("title like ?", "%"+keyword+"%")
	}
	if status != 0 {
		find = find.Where("status=?", status)
	} else {
		find = find.Where("status!=0")
	}
	if len(labels) != 0 {
		find = find.Where("labels in ?", labels)
	}
	find = find.Limit(page.Size()).Offset((page.Page() - 1) * page.Size())
	if err := find.Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (i *issue) FindById(issue int64) (*entity.ScriptIssue, error) {
	ret := &entity.ScriptIssue{ID: issue}
	if err := i.db.First(ret).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (i *issue) Save(issue *entity.ScriptIssue) error {
	return i.db.Save(issue).Error
}
